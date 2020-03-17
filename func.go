package dailytoggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	conf   *config
	toggl  TogglClient
	pixela PixelaClient
	loc    *time.Location
)

type config struct {
	AuthToken        string `split_words:"true" required:"true"`
	PixelaToken      string `split_words:"true" required:"true"`
	PixelaUsername   string `split_words:"true" required:"true"`
	PixelaGraphID    string `split_words:"true" required:"true"`
	TogglTimeZone    string `split_words:"true" required:"true"`
	TogglAPIToken    string `split_words:"true" required:"true"`
	TogglProjectID   string `split_words:"true" required:"true"`
	TogglWorkspaceID string `split_words:"true" required:"true"`
}

type requestBody struct {
	AuthToken string `json:"auth_token"`
	Date      string `json:"date"`
}

func init() {
	conf = &config{}
	if err := envconfig.Process("", conf); err != nil {
		log.Fatal(err.Error())
	}

	toggl = newToggl(conf)
	pixela = newPixelaClient(conf)

	var err error
	if loc, err = time.LoadLocation(conf.TogglTimeZone); err != nil {
		log.Fatal(err.Error())
	}
}

// DailyToggl aggregate toggl times.
func DailyToggl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("request method is: %s", r.Method)
		http.Error(w, "unsupported http method", http.StatusBadRequest)
		return
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "failed to read request", http.StatusBadRequest)
		return
	}

	body := requestBody{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "failed to parse request", http.StatusBadRequest)
		return
	}

	if body.AuthToken != conf.AuthToken {
		http.Error(w, "authentication error", http.StatusUnauthorized)
		return
	}

	date, err := getTargetDate(body.Date)
	if err != nil {
		log.Printf("getTargetDate: %v", err)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	total, err := toggl.getDayTotal(date)
	if err != nil {
		log.Printf("toggl.getDayTotal: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := pixela.update(date, total); err != nil {
		log.Printf("pixela.update: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("%d", total))
}

// getTargetDate returns target date as time.Time.
// If date is given, getTargetDate uses it.
// If not, getTargetDate returns yesterday.
func getTargetDate(str string) (time.Time, error) {
	if str != "" {
		return time.Parse("2006-01-02", str)
	}

	yesterday := time.Now().Add(-24 * time.Hour)
	return yesterday.In(loc), nil
}

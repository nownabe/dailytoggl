package dailytoggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	conf    *config
	toggl   TogglClient
	pixela  PixelaClient
	twitter TwitterClient
	loc     *time.Location
)

type config struct {
	AuthToken                string `split_words:"true" required:"true"`
	PixelaToken              string `split_words:"true" required:"true"`
	PixelaUsername           string `split_words:"true" required:"true"`
	PixelaGraphID            string `split_words:"true" required:"true"`
	TogglTimeZone            string `split_words:"true" required:"true"`
	TogglAPIToken            string `split_words:"true" required:"true"`
	TogglProjectID           string `split_words:"true" required:"true"`
	TogglWorkspaceID         string `split_words:"true" required:"true"`
	TwitterAccessToken       string `split_words:"true"`
	TwitterAccessTokenSecret string `split_words:"true"`
	TwitterConsumerKey       string `split_words:"true"`
	TwitterConsumerSecret    string `split_words:"true"`
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
	if conf.TwitterAccessToken != "" {
		twitter = newTwitterClient(conf)
	}

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

	if twitter != nil {
		if err := tweet(total, date); err != nil {
			log.Printf("tweet: %v", err)
		}
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

func tweet(total int64, date time.Time) error {
	svg, err := pixela.getGraph()
	if err != nil {
		return err
	}

	svgFile, err := ioutil.TempFile("", "dailytoggl.*.svg")
	if err != nil {
		return err
	}
	defer os.Remove(svgFile.Name())
	log.Printf("svg tempfile: %s", svgFile.Name())

	pngFile, err := ioutil.TempFile("", "dailytoggl.*.png")
	if err != nil {
		return err
	}
	defer os.Remove(pngFile.Name())
	log.Printf("png tempfile: %s", pngFile.Name())

	if _, err := svgFile.Write(svg); err != nil {
		svgFile.Close()
		return err
	}

	if err := svgFile.Close(); err != nil {
		return err
	}

	if err := exec.Command("convert",
		"-density", "100",
		"-background", "none",
		svgFile.Name(), pngFile.Name()).Run(); err != nil {
		return err
	}

	png, err := ioutil.ReadAll(pngFile)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("On %s, I studied English for %d minutes.\n%s.html",
		date.Format("2006-01-02"), total, pixela.graphURL())

	if err := twitter.post(msg, png); err != nil {
		return err
	}

	return nil
}

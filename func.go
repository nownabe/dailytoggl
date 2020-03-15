package dailytoggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

var conf *config

type config struct {
	AuthToken string `split_words:"true" required:"true"`
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
}

// DailyToggl aggregate toggl times.
func DailyToggl(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, "ok")
}

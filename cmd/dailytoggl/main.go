package main

import (
	"log"
	"net/http"

	"github.com/nownabe/dailytoggl"
)

func main() {
	http.HandleFunc("/", dailytoggl.DailyToggl)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

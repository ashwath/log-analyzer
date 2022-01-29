package main

import (
	"logAnalyzer/internal/analyzer"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {

	// All public routes are populated here
	s, err := analyzer.New()
	if err != nil {
		log.WithError(err).Error("error creating API")
		return
	}

	err = http.ListenAndServe(":8080", s)
	if err != nil {
		log.WithError(err).Error("error starting  API ")
		return
	}
}

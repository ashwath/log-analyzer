package analyzer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	logFilesPath = "/var/log"
)

func GetFileName(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	if len(fileName) == 0 {
		log.WithError(errors.New("filename not provided"))
		return fileName, errors.New("filename not provided")
	}
	return fileName, nil
}

func GetLogLinesRequested(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	logLinesStr := vars["logLines"]
	logLines, err := strconv.Atoi(logLinesStr)

	if err != nil || logLines == 0 {
		// log and swallow it
		log.WithError(err).Error("invalid input for 'logLines'")
		return 0, errors.New("invalid input for 'logLines'")
	}
	return logLines, nil
}

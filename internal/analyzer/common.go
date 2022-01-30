package analyzer

import (
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	logFilesPath = "/var/log"
)

func GetFileName(r *http.Request) (string, error) {
	fileName := r.FormValue("filename")
	log.Infof("GetFileName():fileName: %+v\n", fileName)
	if len(fileName) == 0 {
		log.WithError(errors.New("filename not provided"))
		return fileName, errors.New("filename not provided")
	}
	return fileName, nil
}

func GetLogLinesRequested(r *http.Request) (int, error) {
	logLinesStr := r.FormValue("logLines")
	logLines, err := strconv.Atoi(logLinesStr)
	if err != nil || logLines == 0 {
		// log and swallow it
		log.WithError(err).Error("invalid input for 'logLines'")
		return 0, errors.New("invalid input for 'logLines'")
	}
	return logLines, nil
}

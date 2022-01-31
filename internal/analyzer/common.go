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

var (
	fileNameErr       = errors.New("filename not provided")
	numberLogLinesErr = errors.New("invalid input for 'logs lines'")
	searchKeyWordErr  = errors.New("search keyword not provided")
)

type SearchResponse struct {
	FileName string   `json:"file_path"`
	Logs     []string `json:"logs"`
}

func GetFileName(r *http.Request) (string, error) {
	fileName := r.FormValue("file_name")
	log.Infof("GetFileName():fileName: %+v\n", fileName)
	//if len(fileName) == 0 {
	//	log.WithError(fileNameErr)
	//	return fileName, fileNameErr
	//}
	return fileName, nil
}

func GetLogLinesRequested(r *http.Request) (int, error) {
	logLinesStr := r.FormValue("log_lines")
	logLines, err := strconv.Atoi(logLinesStr)
	if err != nil || logLines == 0 {
		log.WithError(err).Error(numberLogLinesErr)
		return 0, numberLogLinesErr
	}
	return logLines, nil
}

func GetSearchKeyWord(r *http.Request) (string, error) {
	keyword := r.FormValue("keyword")
	if len(keyword) == 0 {
		log.WithError(errors.New("search keyword not provided"))
		return keyword, errors.New("search keyword not provided")
	}
	return keyword, nil
}

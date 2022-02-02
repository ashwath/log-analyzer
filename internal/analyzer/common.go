package analyzer

import (
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	logFilesPath         = "/var/log"
	defaultLogEntryLimit = 100
)

var (
	fileNameErr       = errors.New("filename not provided")
	numberLogLinesErr = errors.New("invalid input for 'logs lines'")
	searchKeyWordErr  = errors.New("search keyword not provided")
)

type Response struct {
	Results  []SearchResults  `json:"results"`
	MetaData ResponseMetadata `json:"response_metadata"`
}

type SearchResults struct {
	FileName   string   `json:"file_path"`
	LogEntries []string `json:"log_entries"`
}

type ResponseMetadata struct {
	CurrentFile string `json:"current_file"`
	NextCursor  string `json:"next_cursor"`
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

func GetLogEntryLimit(r *http.Request) (int, error) {
	limitStr := r.FormValue("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.WithError(err).Error(numberLogLinesErr)
		return 0, numberLogLinesErr
	}
	if limit == 0 {
		limit = defaultLogEntryLimit
	}
	return limit, nil
}

func GetSearchKeyword(r *http.Request) (string, error) {
	keyword := r.FormValue("keyword")
	if len(keyword) == 0 {
		log.WithError(errors.New("search keyword not provided"))
		return keyword, errors.New("search keyword not provided")
	}
	return keyword, nil
}

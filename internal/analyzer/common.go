package analyzer

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	logFilesPath         = "/var/log"
	defaultLogEntryLimit = 20
	defaultLastN         = 10
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
	NextFile   string `json:"next_file"`
	NextCursor int64  `json:"next_cursor"`
}

func GetFileName(r *http.Request) string {
	fileName := r.FormValue("file_name")
	return fileName
}

func GetLogLimit(r *http.Request) (int, error) {
	limit := defaultLogEntryLimit
	limitStr := r.FormValue("limit")
	if len(limitStr) > 0 {
		limitReq, err := strconv.Atoi(limitStr)
		if err != nil {
			log.WithError(err).Error(LimitErr)
			return 0, LimitErr
		}
		if limitReq > 0 {
			limit = limitReq
		}
	}
	return limit, nil
}

func GetLastN(r *http.Request) (int, error) {
	lastN := defaultLastN
	lastNStr := r.FormValue("lastN")
	if len(lastNStr) > 0 {
		lastNReq, err := strconv.Atoi(lastNStr)
		if err != nil {
			log.WithError(err).Error(LastNErr)
			return 0, LastNErr
		}
		if lastNReq > 0 {
			lastN = lastNReq
		}
	}
	return lastN, nil
}

func GetNextCursor(r *http.Request) (int64, error) {
	cursor := int64(0)
	cursorStr := r.FormValue("next_cursor")
	if len(cursorStr) > 0 {
		cursorProvided, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || cursorProvided < 0 {
			log.WithError(err).Error(NextCursorErr)
			return cursor, NextCursorErr
		}
		cursor = cursorProvided
	}
	return cursor, nil
}

func GetNextFile(r *http.Request) string {
	nextFile := r.FormValue("next_file")
	return nextFile
}

func GetSearchKeyword(r *http.Request) string {
	keyword := r.FormValue("keyword")
	if len(keyword) == 0 {
		log.Infof("Search keyword not provided, will tail logs")
	}
	return keyword
}

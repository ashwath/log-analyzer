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
	nextCursorErr     = errors.New("invalid input for 'next cursor'")
	searchKeyWordErr  = errors.New("search keyword not provided")
	pagingMetadataErr = errors.New("invalid metadata provided to page the results")
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
	log.Infof("GetFileName(): %+v\n", fileName)
	//if len(fileName) == 0 {
	//	log.WithError(fileNameErr)
	//	return fileName, fileNameErr
	//}
	return fileName
}

func GetLogEntryLimit(r *http.Request) (int, error) {
	limit := defaultLogEntryLimit
	limitStr := r.FormValue("limit")
	if len(limitStr) > 0 {
		limitReq, err := strconv.Atoi(limitStr)
		if err != nil {
			log.WithError(err).Error(numberLogLinesErr)
			return 0, numberLogLinesErr
		}
		if limitReq > 0 {
			limit = limitReq
		}
	}
	return limit, nil
}

func GetNextCursor(r *http.Request) (int64, error) {
	cursor := int64(0)
	cursorStr := r.FormValue("next_cursor")
	if len(cursorStr) > 0 {
		cursorProvided, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || cursorProvided < 0 {
			log.WithError(err).Error(nextCursorErr)
			return cursor, nextCursorErr
		}
		cursor = cursorProvided
	}
	return cursor, nil
}

func GetNextFile(r *http.Request) string {
	nextFile := r.FormValue("next_file")
	log.Infof("GetNextFile(): %+v\n", nextFile)
	return nextFile
}

func GetSearchKeyword(r *http.Request) (string, error) {
	keyword := r.FormValue("keyword")
	if len(keyword) == 0 {
		log.WithError(errors.New("search keyword not provided"))
		return keyword, errors.New("search keyword not provided")
	}
	return keyword, nil
}

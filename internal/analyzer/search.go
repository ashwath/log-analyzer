package analyzer

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) search(w http.ResponseWriter, r *http.Request) error {
	log.Infof("http.Request: %+v\n", *r)

	// fetch number of log entries requested
	limit, err := GetLogLimit(r)
	if err != nil {
		return err
	}
	log.Infof("logs: %+v\n", limit)

	// fetch number of log lines requested
	keyword, err := GetSearchKeyword(r)
	if err != nil {
		return err
	}
	log.Infof("keyword: %+v\n", keyword)

	// fetch filename
	var allLogFiles bool
	fileName := GetFileName(r)
	if len(fileName) == 0 {
		allLogFiles = true
		log.Infof("Requested to search across all log files, filename: %+v\n", fileName)
	} else {
		log.Infof("filename: %+v\n", fileName)
	}

	cursor, err := GetNextCursor(r)
	if err != nil {
		return err
	}
	log.Infof("cursor: %+v\n", cursor)

	nextFile := GetNextFile(r)
	if len(nextFile) == 0 && cursor != 0 {
		return pagingMetadataErr
	}

	// Search in the logs
	response := Response{}
	startScaning := false
	if allLogFiles {
		/*
			Filepath.Walk()
			The files are walked in lexical order, which makes the output deterministic.
			We need the output to be deterministic for paging logic to work
		*/
		err := filepath.Walk(logFilesPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				log.Debugf("path:%v, info.Size():%v\n", path, info.Size())
				if len(nextFile) > 0 {
					if path == nextFile {
						startScaning = true
					}
					if startScaning && !info.IsDir() {
						err = ScanLogFile(path, keyword, limit, cursor, &response)
						if err != nil {
							return err
						}
					}
				} else {
					if !info.IsDir() {
						err = ScanLogFile(path, keyword, limit, cursor, &response)
						if err != nil {
							return err
						}
					}
				}
				return nil
			})
		if err != nil {
			log.Error(err)
		}
	} else {
		// fetch corresponding log lines
		path := logFilesPath + "/" + fileName
		err := ScanLogFile(path, keyword, limit, cursor, &response)
		if err != nil {
			return err
		}
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(response)
	return nil
}

func ScanLogFile(filepath, keyword string, limit int, cursor int64, response *Response) error {
	log.Infof("Scaning file: %v\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Seek(cursor, 0)
	reader := bufio.NewReader(file)
	logs := []string{}
	for {
		bytes, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("Error while reading the file %v\n", err)
			return err
		}
		line := string(bytes)
		if len(keyword) > 0 {
			if limit > 0 && strings.Contains(line, keyword) {
				logs = append(logs, line)
				limit--
			}
		} else {
			logs = append(logs, line)
			limit--
		}

		// check if limit is reached
		if limit == 0 {
			log.Infof("limit == 0")
			results := response.Results
			results = append(results, SearchResults{FileName: filepath, LogEntries: logs})
			response.Results = results
			response.MetaData.NextFile = filepath
			offset, err := file.Seek(0, io.SeekCurrent) // offset is the current position
			if err != nil {
				log.Errorf("Error while finindg file:%s offset %v\n", filepath, err)
				return err
			}
			response.MetaData.NextCursor = offset
		}
	}

	if len(logs) > 0 {
		results := response.Results
		results = append(results, SearchResults{FileName: filepath, LogEntries: logs})
		response.Results = results
	}
	return nil
}

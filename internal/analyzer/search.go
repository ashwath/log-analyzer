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
	limit, err := GetLogEntryLimit(r)
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
	fileName, err := GetFileName(r)
	if len(fileName) == 0 {
		allLogFiles = true
		log.Infof("Requested to search across all log files, filename: %+v\n", fileName)
	} else {
		log.Infof("filename: %+v\n", fileName)
	}

	// Search in the logs
	response := Response{}
	if allLogFiles {
		/*
			Filepath.Walk()
			The files are walked in lexical order, which makes the output deterministic
		*/
		err := filepath.Walk(logFilesPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				log.Infof("path:%v, info.Size():%v\n", path, info.Size())
				if !info.IsDir() {
					err = ScanLogFile(path, keyword, limit, &response)
					if err != nil {
						return err
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
		err := ScanLogFile(path, keyword, limit, &response)
		if err != nil {
			return err
		}
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(response)
	return nil
}

func ScanLogFile(filepath, keyword string, limit int, response *Response) error {
	log.Infof("Scaning file: %v\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Seek(0, 0)
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
		if limit > 0 && strings.Contains(line, keyword) {
			logs = append(logs, line)
			limit--
		}
		if limit == 0 {
			results := response.Results
			results = append(results, SearchResults{FileName: filepath, LogEntries: logs})
			response.Results = results
			response.MetaData.CurrentFile = filepath
		}
	}
	if len(logs) > 0 {
		results := response.Results
		results = append(results, SearchResults{FileName: filepath, LogEntries: logs})
		response.Results = results
		response.MetaData.CurrentFile = filepath
		response.MetaData.NextCursor = ""
	}
	return nil
}

// Reads each line at a time
func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		/*
			// ReadLine is a low-level line-reading primitive
			// ReadLine tries to return a single line, not including the end-of-line bytes.
			// If the line was too long for the buffer then isPrefix is set and the
			// beginning of the line is returned. The rest of the line will be returned
			// from future calls. isPrefix will be false when returning the last fragment
			// of the line
		*/
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

package analyzer

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) search(w http.ResponseWriter, r *http.Request) error {
	log.Infof("http.Request: %+v\n", *r)

	// fetch number of log lines requested
	keyword, err := GetSearchKeyWord(r)
	if err != nil {
		return err
	}
	log.Infof("keyword: %+v\n", keyword)

	// fetch Filename
	var allLogFiles bool
	fileName, err := GetFileName(r)
	if len(fileName) == 0 {
		allLogFiles = true
		log.Infof("Requested to search across all log files, filename: %+v\n", fileName)
	} else {
		log.Infof("filename: %+v\n", fileName)
	}

	results := []SearchResponse{}
	if allLogFiles {
		err := filepath.Walk(logFilesPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				log.Infof("path:%v, info.Size():%v\n", path, info.Size())
				if !info.IsDir() {
					err = h.searchInFiles(path, keyword, &results)
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
		err := h.searchInFiles(logFilesPath+"/"+fileName, keyword, &results)
		if err != nil {
			return err
		}
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(results)
	return nil
}

func (h *Handler) searchInFiles(filepath, keyword string, results *[]SearchResponse) error {
	log.Infof("Reading file: %v\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// assumed max length of log line is 64000 characters
	const maxCapacity = 64000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	logs := []string{}
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			logs = append(logs, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(logs) > 0 {
		*results = append(*results, SearchResponse{FileName: filepath, Logs: logs})
	}

	return nil
}

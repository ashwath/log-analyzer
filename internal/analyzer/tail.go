package analyzer

import (
	"encoding/json"
	"logAnalyzer/internal/file_utils"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) tailLogs(w http.ResponseWriter, r *http.Request) error {
	// fetch Filename
	fileName := GetFileName(r)
	if len(fileName) == 0 {
		log.WithError(fileNameErr)
		return fileNameErr
	}

	// fetch number of log lines requested
	lastN, err := GetLastN(r)
	if err != nil {
		return err
	}

	// fetch corresponding log lines
	results, err := h.fetchLastN(fileName, lastN)
	if err != nil {
		return err
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(results)
}

func (h *Handler) fetchLastN(fileName string, n int) (Response, error) {
	f, err := os.Open(logFilesPath + "/" + fileName)
	if err != nil {
		log.Errorf("File not found, %v", err)
		return Response{}, err
	}
	fi, err := f.Stat() // returns file info
	if err != nil {
		log.Error(err)
		return Response{}, err
	}
	defer f.Close()
	log.Printf("fileName:%v lastN:%v\n", fileName, n)
	scanner := file_utils.NewScanner(f, int(fi.Size())) // set pos to end of file
	logs := []string{}
	for {
		log.Println("here")
		line, _, err := scanner.Line()
		if err != nil {
			log.Println("here")
			log.Error(err)
			return Response{}, err
			break
		}
		log.Printf("line: %s\n", line)
		logs = append(logs, line)
		if len(logs) == n {
			break
		}
	}
	response := Response{Results: []SearchResults{{FileName: fileName, LogEntries: logs}}}
	return response, nil
}

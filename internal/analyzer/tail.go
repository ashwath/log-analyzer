package analyzer

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) tailLogs(w http.ResponseWriter, r *http.Request) error {
	log.Infof("http.Request: %+v\n", *r)

	// fetch Filename
	fileName, err := GetFileName(r)
	if len(fileName) == 0 {
		log.WithError(fileNameErr)
		return fileNameErr
	}
	log.Infof("filename: %+v\n", fileName)

	// fetch number of log lines requested
	logLines, err := GetLogEntryLimit(r)
	if err != nil {
		return err
	}
	log.Infof("logs: %+v\n", logLines)

	// fetch corresponding log lines
	results, err := h.fetchLastN(fileName, logLines)
	if err != nil {
		return err
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(results)
}

func (h *Handler) fetchLastN(fileName string, n int) (SearchResults, error) {
	file, err := os.Open(logFilesPath + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// assumed max length of log line is 64000 characters
	const maxCapacity = 64000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	result := []string{}
	for scanner.Scan() {
		if len(result) == n {
			result = result[1:]
		}
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return SearchResults{FileName: fileName, LogEntries: result}, nil
}

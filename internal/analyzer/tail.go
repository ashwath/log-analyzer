package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (h *Handler) tailLogs(w http.ResponseWriter, r *http.Request) error {
	// fetch Filename
	fileName, err := GetFileName(r)
	if err != nil {
		return err
	}

	// fetch number of log lines requested
	logLines, err := GetLogLinesRequested(r)
	if err != nil {
		return err
	}

	// fetch corresponding log lines
	lastNLogLines, err := h.fetchLastN(fileName, logLines)
	if err != nil {
		return err
	}

	for _, v := range lastNLogLines {
		fmt.Println(v)
	}

	// encode the results to http.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(lastNLogLines)
}

func (h *Handler) fetchLastN(fileName string, n int) ([]string, error) {
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

	return result, nil
}

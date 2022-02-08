//go:build integration
// +build integration

package automated_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logAnalyzer/internal/analyzer"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTailAPISuccessCases(t *testing.T) {
	fileName := "sample.log"
	lastN := 5

	resourceURL := fmt.Sprintf("%s?file_name=%s&lastN=%d", TailURL, fileName, lastN)
	log.Infof("GET:%s", resourceURL)
	resp, err := http.Get(resourceURL)
	assert.NoError(t, err, "expected no errors")
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no errors")

	var response analyzer.Response
	err = json.Unmarshal(body, &response)
	log.Infof("Response:%+v", response)
	assert.NoError(t, err, "Can't unmarshal into Response")
}

func TestTailAPIFailureCases(t *testing.T) {
	fileName := "sample.log"
	lastN := 5

	testCases := []struct {
		name          string
		requestUrl    string
		expectedError analyzer.ErrorResponse
	}{
		{
			name:       "Invalid LastN value",
			requestUrl: fmt.Sprintf("%s?file_name=%s&lastN=%s", TailURL, fileName, "xxx"),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         analyzer.LastNErr.Error(),
				InternalMessage: analyzer.LastNErr.Error(),
			},
		},
		{
			name:       "File Not found",
			requestUrl: fmt.Sprintf("%s?file_name=%s&lastN=%d", TailURL, "xxx", lastN),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusInternalServerError,
				Name:            http.StatusText(http.StatusInternalServerError),
				Message:         fmt.Sprintf("open %s: no such file or directory", "/var/log/xxx"),
				InternalMessage: fmt.Sprintf("open %s: no such file or directory", "/var/log/xxx"),
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			log.Infof("GET on: %+v", tc.requestUrl)
			resp, err := http.Get(tc.requestUrl)
			defer resp.Body.Close()
			var output analyzer.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&output)
			log.Infof("Response: %+v", output)
			assert.NoError(t, err, "unexpected error during unmarshalling payload")
			assert.Equal(t, tc.expectedError.Status, output.Status)
			assert.Equal(t, tc.expectedError.Message, output.Message)
			assert.Equal(t, tc.expectedError.InternalMessage, output.InternalMessage)
		})
	}
}

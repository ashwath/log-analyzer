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

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var (
	BaseURL   = "http://localhost:4200"
	SearchURL = fmt.Sprintf("%s/v1/logs/search", BaseURL)
	TailURL   = fmt.Sprintf("%s/v1/logs/tail", BaseURL)
)

func TestSearchAPISuccessCases(t *testing.T) {
	fileName := "sample.log"
	keyword := "1198913921"

	resourceURL := fmt.Sprintf("%s?file_name=%s&keyword=%s", SearchURL, fileName, keyword)
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

func TestSearchAPIFailureCases(t *testing.T) {
	fileName := "sample.log"
	keyword := "1198913921"
	limit := 5
	nextCursor := 98978

	testCases := []struct {
		name          string
		requestUrl    string
		expectedError analyzer.ErrorResponse
	}{
		{
			name:       "Invalid Limit",
			requestUrl: fmt.Sprintf("%s?file_name=%s&keyword=%s&limit=%s", SearchURL, fileName, keyword, "xxx"),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         analyzer.LimitErr.Error(),
				InternalMessage: analyzer.LimitErr.Error(),
			},
		},
		{
			name: "Invalid Next Cursor",
			requestUrl: fmt.Sprintf("%s?file_name=%s&keyword=%s&limit=%d&next_cursor=%s", SearchURL, fileName,
				keyword, limit, "xxx"),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         analyzer.NextCursorErr.Error(),
				InternalMessage: analyzer.NextCursorErr.Error(),
			},
		},
		{
			name: "Invalid Metadata",
			requestUrl: fmt.Sprintf("%s?file_name=%s&keyword=%s&limit=%d&next_cursor=%d", SearchURL, fileName,
				keyword, limit, nextCursor),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         analyzer.PagingMetadataErr.Error(),
				InternalMessage: analyzer.PagingMetadataErr.Error(),
			},
		},
		{
			name: "Invalid Next File",
			requestUrl: fmt.Sprintf("%s?file_name=%s&keyword=%s&limit=%d&next_cursor=%d&next_file=%s", SearchURL, fileName,
				keyword, limit, nextCursor, "xxx"),
			expectedError: analyzer.ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         analyzer.InvalidNextFileErr.Error(),
				InternalMessage: analyzer.InvalidNextFileErr.Error(),
			},
		},
		{
			name:       "File Not found",
			requestUrl: fmt.Sprintf("%s?file_name=%s&keyword=%s&limit=%d", SearchURL, "xxx", keyword, limit),
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

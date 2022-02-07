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

var (
	baseURL   = "http://localhost:4200"
	searchURL = fmt.Sprintf("%s/v1/logs/search", baseURL)
)

func TestSearchAPI(t *testing.T) {
	fileName := "sample.log"
	keyword := "1198913921"

	resourceURL := fmt.Sprintf("%s?file_name=%s&keyword=%s", searchURL, fileName, keyword)
	log.Infof("Executing GET:%s", resourceURL)
	resp, err := http.Get(resourceURL)
	log.Debugf("Response:%+v", resp)
	assert.NoError(t, err, "expected no errors")
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	defer resp.Body.Close()
	log.Infof("Response Code:%+v", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no errors")

	var response analyzer.Response
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err, "Can't unmarshal into Response")
}

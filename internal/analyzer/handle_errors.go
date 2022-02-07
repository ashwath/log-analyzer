package analyzer

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	fileNameErr        = errors.New("{file_name} not provided")
	lastNErr           = errors.New("invalid input for {lastN}")
	limitErr           = errors.New("invalid input for {limit}")
	nextCursorErr      = errors.New("invalid input for {next_cursor}")
	searchKeywordErr   = errors.New("search {keyword} not provided")
	pagingMetadataErr  = errors.New("invalid metadata provided, please provide {next_file} with {cursor}")
	invalidNextFileErr = errors.New("invalid input combination for {next_file} & {file_name}")
)

// apiError error response for API requests
type apiError struct {
	Status          int    `json:"status"`
	Name            string `json:"name"`
	Message         string `json:"message"`
	InternalMessage string `json:"internal_message"`
}

type handler func(w http.ResponseWriter, req *http.Request) error

// HandleError error handler
func HandleError(f handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		log.WithError(err).Error(errors.WithStack(err))

		var e apiError
		switch err {
		case fileNameErr, limitErr, lastNErr, nextCursorErr, searchKeywordErr, pagingMetadataErr:
			e = apiError{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         err.Error(),
				InternalMessage: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)

		default:
			e = apiError{
				Status:          http.StatusInternalServerError,
				Name:            http.StatusText(http.StatusInternalServerError),
				Message:         err.Error(),
				InternalMessage: err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
		}

		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(e)
	}
}

package analyzer

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	FileNameErr        = errors.New("{file_name} not provided")
	LastNErr           = errors.New("invalid input for {lastN}")
	LimitErr           = errors.New("invalid input for {limit}")
	NextCursorErr      = errors.New("invalid input for {next_cursor}")
	SearchKeywordErr   = errors.New("search {keyword} not provided")
	PagingMetadataErr  = errors.New("invalid metadata provided, please provide {next_file} with {cursor}")
	InvalidNextFileErr = errors.New("invalid input combination for {next_file} & {file_name}")
)

// ErrorResponse error response for API requests
type ErrorResponse struct {
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

		var e ErrorResponse
		switch err {
		case FileNameErr, LimitErr, LastNErr, NextCursorErr, SearchKeywordErr, PagingMetadataErr, InvalidNextFileErr:
			e = ErrorResponse{
				Status:          http.StatusBadRequest,
				Name:            http.StatusText(http.StatusBadRequest),
				Message:         err.Error(),
				InternalMessage: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)

		default:
			e = ErrorResponse{
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

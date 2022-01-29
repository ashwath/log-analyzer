package analyzer

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type handler func(w http.ResponseWriter, req *http.Request) error

type apiError struct {
	Code            int    `json:"code"`
	Status          int    `json:"status"`
	Name            string `json:"name"`
	Message         string `json:"message"`
	InternalMessage string `json:"internal_message"`
	MoreInfo        string `json:"more_info"`
}

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
		default:

			e = apiError{
				Code:            http.StatusInternalServerError,
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

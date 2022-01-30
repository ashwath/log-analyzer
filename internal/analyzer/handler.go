package analyzer

import (
	"logAnalyzer/config"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	getSearchPath = "/v1/logs/search"
	getTailPath   = "/v1/logs/tail"
)

// Handler implements http.Handler
type Handler struct {
	router http.Handler
	cfg    config.Config
}

// ServeHTTP serves HTTP
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// New returns a server allowing for dependency injection
func New(opts ...func(h *Handler)) (*Handler, error) {

	h := &Handler{
		cfg: config.Get(),
	}

	for _, opt := range opts {
		opt(h)
	}

	if h.router == nil {
		r := mux.NewRouter()
		r.Handle(getSearchPath, HandleError(h.search)).Methods(http.MethodGet)
		r.Handle(getTailPath, HandleError(h.tailLogs)).Methods(http.MethodGet)
		h.router = r
	}

	log.Debug("server created")

	return h, nil
}

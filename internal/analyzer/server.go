package analyzer

import (
	"logAnalyzer/config"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	getSearchPath = "/v1/search"
	getTailPath   = "/v1/tail"
)

// Server implements http.Handler
type Server struct {
	router http.Handler
	cfg    config.Config
}

// ServeHTTP serves HTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// New returns a server allowing for dependency injection
func New(opts ...func(s *Server)) (*Server, error) {

	s := &Server{
		cfg: config.Get(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.router == nil {
		r := mux.NewRouter()
		r.Handle(getSearchPath, HandleError(s.search)).Methods(http.MethodGet)
		r.Handle(getTailPath, HandleError(s.tailLogs)).Methods(http.MethodGet)
		s.router = r
	}

	log.Debug("server created")

	return s, nil
}

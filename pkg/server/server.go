package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	json "github.com/pquerna/ffjson/ffjson"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/logger"
	"github.com/saromanov/golang-developer-test-task/pkg/storage"
	"gopkg.in/tylerb/graceful.v1"
)

// Server provides implementation of the main server
type Server struct {
	store  storage.Storage
	logger *logger.Logger
}

// search provides searching by the data
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func(){
		finishTime := time.Since(startTime)
		responseTime.Observe(float64(finishTime.Milliseconds()))

	}()
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	totalRequests.Inc()
	if len(r.URL.Query()) == 0 {
		http.Error(w, "search parameters is not defined. Should be provides id, global_id or mode_id", http.StatusBadRequest)
		return
	}
	data, err := s.store.Find(s.prepareSearchRequest(r))
	if err != nil {
		writeStatusCode(http.StatusInternalServerError, "GET")
		failedRequests.Inc()
		http.Error(w, fmt.Sprintf("unable to find data: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		failedRequests.Inc()
		writeStatusCode(http.StatusInternalServerError, "GET")
		http.Error(w, fmt.Sprintf("unable to marshal data: %v", err), http.StatusInternalServerError)
		return
	}
	writeStatusCode(http.StatusOK, "GET")
	w.WriteHeader(http.StatusOK)
	w.Write(resp) //nolint
}

// mustParseInt always returns number from request
// but if its contains errors, its logging
func (s *Server) mustParseInt(d string) int64 {
	i, err := strconv.ParseInt(d, 10, 32)
	if err != nil {
		s.logger.Errorf("unable to parse input data: %v", err)
		return 0
	}
	return i

}

// makeServer provides create of the new server
func makeServer(st storage.Storage, c *config.Config) (*http.ServeMux, *http.Server) {
	server := &Server{
		store:  st,
		logger: c.Logger,
	}

	s := http.NewServeMux()
	s.HandleFunc("/v1/search", server.search)
	s.Handle("/metrics", promhttp.Handler())
	s.Handle("/favicon.ico", http.NotFoundHandler())
	c.Logger.Infof("starting of the server at %s...", c.Address)

	srv := &http.Server{
		Handler:      s,
		Addr:         c.Address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s, srv
}

// Make provides starting of the server
func Make(st storage.Storage, c *config.Config) error {
	s, srv := makeServer(st, c)
	initPrometheus()
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	graceful.Run(c.Address, 10*time.Second, s)
	return nil
}

package server

import (
	"fmt"
	"net/http"
	"time"

	json "github.com/pquerna/ffjson/ffjson"
	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/storage"
	"gopkg.in/tylerb/graceful.v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server provides implementation of the main server
type Server struct {
	store storage.Storage
}

// search provides searching by the data
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.Find(&storage.FindConfig{})
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to find data: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal data: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// prometheus provides exporting of prometheus data
func (s *Server) prometheus(w http.ResponseWriter, r *http.Request) {

}

// Make provides starting of the server
func Make(st storage.Storage, c *config.Config) error {
	server := &Server{
		store: st,
	}
	s := http.NewServeMux()
	s.HandleFunc("/v1/search", server.search)
	s.HandleFunc("/v1/prom", server.prometheus)
	s.Handle("/metrics", promhttp.Handler())
	c.Logger.Infof("starting of the server...")
	initPrometheus()
	http.ListenAndServe(c.Address, nil)
	graceful.Run(c.Address, 10*time.Second, s)
	return nil
}

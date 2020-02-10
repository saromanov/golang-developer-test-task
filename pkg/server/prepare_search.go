package server

import (
	"net/http"

	"github.com/saromanov/golang-developer-test-task/pkg/storage"
)

// prepareSearchRequest creates request for find on storage
func (s *Server) prepareSearchRequest(r *http.Request) *storage.FindConfig {
	response := &storage.FindConfig{}
	globalID, ok := r.URL.Query()["global_id"]
	if ok && len(globalID[0]) > 0 {
		response.GlobalID = s.mustParseInt(globalID[0])
	}
	id, ok := r.URL.Query()["id"]
	if ok && len(id[0]) > 0 {
		response.ID = s.mustParseInt(id[0])
	}
	mode, ok := r.URL.Query()["mode"]
	if ok && len(mode[0]) > 0 {
		response.ModeID = mode[0]
	}
	return response
}

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/logger"
	"github.com/saromanov/golang-developer-test-task/pkg/models"
	"github.com/saromanov/golang-developer-test-task/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var errNoElement = errors.New("unable to find element")

type storageTmp struct {
	mu   sync.RWMutex
	data map[string]models.Parking
}

func (s *storageTmp) Insert(data []models.Parking) (int, error) {
	for _, d := range data {
		s.mu.Lock()
		s.data[fmt.Sprintf("id_%d", d.ID)] = d
		s.data[fmt.Sprintf("gid_%d", d.GlobalID)] = d
		s.data[fmt.Sprintf("mode_%s", d.Mode)] = d
		s.mu.Unlock()
	}
	return len(data), nil
}

func (s *storageTmp) Find(req *storage.FindConfig) ([]models.Parking, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fmt.Println("REQ: ", req.GlobalID)
	if req.ID != 0 {
		res, ok := s.data[fmt.Sprintf("id_%d", req.ID)]
		if !ok {
			return nil, errNoElement
		}
		return []models.Parking{res}, nil
	}
	if req.GlobalID != 0 {
		res, ok := s.data[fmt.Sprintf("gid_%d", req.GlobalID)]
		if !ok {
			return nil, errNoElement
		}
		return []models.Parking{res}, nil
	}
	if req.ModeID != "" {
		res, ok := s.data[fmt.Sprintf("mode_%s", req.ModeID)]
		if !ok {
			return nil, errNoElement
		}
		return []models.Parking{res}, nil
	}
	return nil, nil
}

func insertToStorage(t *testing.T, st storage.Storage) {
	_, err := st.Insert([]models.Parking{
		models.Parking{
			ID:       1,
			GlobalID: 1,
			Mode:     "test",
		},
		models.Parking{
			ID:       2,
			GlobalID: 2,
			Mode:     "test2",
		},
		models.Parking{
			ID:       3,
			GlobalID: 3,
			Mode:     "test3",
		},
	})
	assert.NoError(t, err)
}

func TestSearchEndpoint(t *testing.T) {
	st := &storageTmp{
		mu:   sync.RWMutex{},
		data: map[string]models.Parking{},
	}
	insertToStorage(t, st)
	addr := "http://127.0.0.1:8083"
	go func() {
		err := Make(st, &config.Config{
			Address: ":8083",
			Logger:  logger.New(),
		})
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)
	response, err := http.Post(fmt.Sprintf("%s/v1/search", addr), "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusMethodNotAllowed, "status is not equal")
	response, err = http.Get(fmt.Sprintf("%s/v1/search", addr))
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusBadRequest)

	response, err = http.Get(fmt.Sprintf("%s/v1/search?id=1", addr))
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK, "status is not equal")
	dec := []models.Parking{}
	err = json.NewDecoder(response.Body).Decode(&dec)
	assert.NoError(t, err)
	assert.Equal(t, len(dec), 1)

	response, err = http.Get(fmt.Sprintf("%s/v1/search?global_id=1", addr))
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK, "status is not equal")
	dec = []models.Parking{}
	err = json.NewDecoder(response.Body).Decode(&dec)
	assert.NoError(t, err)
	assert.Equal(t, len(dec), 1)

	response, err = http.Get(fmt.Sprintf("%s/v1/search?mode=test", addr))
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK, "status is not equal")
	dec = []models.Parking{}
	err = json.NewDecoder(response.Body).Decode(&dec)
	assert.NoError(t, err)
	assert.Equal(t, len(dec), 1)

	response, err = http.Get(fmt.Sprintf("%s/v1/search?id=5", addr))
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusInternalServerError)

}

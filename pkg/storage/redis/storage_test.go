package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/models"
	"github.com/saromanov/golang-developer-test-task/pkg/storage"
	"github.com/stretchr/testify/assert"
)

type testFind struct {
	Input   []models.Parking
	Request *storage.FindConfig
}

func newTestRedis() storage.Storage {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	client, err := New(&config.Config{
		StorageAddress: mr.Addr(),
	}, redisClient)
	if err != nil {
		panic(err)
	}
	return client
}

func checkParkingStructFind(t *testing.T, r storage.Storage, req *storage.FindConfig) {
	data, err := r.Find(req)
	assert.NoError(t, err)
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].ID, int64(1))
	assert.Equal(t, data[0].Mode, "test")
	assert.Equal(t, data[0].GlobalID, int64(1))
}

func TestInsertParkingData(t *testing.T) {
	r := newTestRedis()
	count, err := r.Insert([]models.Parking{models.Parking{
		ID:       1,
		Mode:     "test",
		GlobalID: 1,
	}})
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	_, err = r.Insert([]models.Parking{models.Parking{
		ID: 0,
	}})
	assert.EqualError(t, err, errNoID.Error())

	count, err = r.Insert([]models.Parking{models.Parking{
		ID:       1,
		Mode:     "test",
		GlobalID: 1,
	},
		models.Parking{
			ID:       2,
			Mode:     "test2",
			GlobalID: 2,
		}})
	assert.NoError(t, err)
	assert.Equal(t, count, 2)
}

func TestFindParkingData(t *testing.T) {
	r := newTestRedis()
	inp := []testFind{testFind{
		Input: []models.Parking{models.Parking{
			ID:       1,
			Mode:     "test",
			GlobalID: 1,
		}},
		Request: &storage.FindConfig{
			ID: 1,
		},
	},

		testFind{
			Input: []models.Parking{models.Parking{
				ID:       1,
				Mode:     "test",
				GlobalID: 1,
			}},
			Request: &storage.FindConfig{
				GlobalID: 1,
			},
		},

		testFind{
			Input: []models.Parking{models.Parking{
				ID:       1,
				Mode:     "test",
				GlobalID: 1,
			}},
			Request: &storage.FindConfig{
				ModeID: "test",
			},
		},
	}

	for _, d := range inp {
		count, err := r.Insert(d.Input)
		assert.NoError(t, err)
		assert.Equal(t, count, len(d.Input))
		checkParkingStructFind(t, r, d.Request)
	}
}

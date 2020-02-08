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
	count, err := r.Insert([]models.Parking{models.Parking{
		ID:       1,
		Mode:     "test",
		GlobalID: 1,
	}})
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	data, err := r.Find(&storage.FindConfig{
		ID: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, len(data), 1)
	assert.Equal(t, data[0].ID, int64(1))
	assert.Equal(t, data[0].Mode, "test")
	assert.Equal(t, data[0].GlobalID, int64(1))
}

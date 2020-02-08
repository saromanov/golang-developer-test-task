package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	json "github.com/pquerna/ffjson/ffjson"
	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/models"
	"github.com/saromanov/golang-developer-test-task/pkg/storage"
)

type Redis struct {
	client *redis.Client
}

// New provides initialization of the redis client
func New(c *config.Config) (storage.Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.StorageAddress,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping redis")
	}

	return &Redis{
		client: client,
	}, nil
}

// Insert provides inserting of list of the data
func (r *Redis) Insert(data []models.Parking) error {
	if len(data) == 0 {
		return nil
	}

	for _, d := range data {
		if d.ID == 0 {
			return errNoID
		}
		d.Name = ""
		d.Mode = ""
		key := fmt.Sprintf("id_%d", d.ID)
		result, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("unable to marshal data: %v", err)
		}
		err = r.client.Do("SET", key, string(result)).Err()
		if err != nil {
			return fmt.Errorf("unable to set data: %v", err)
		}
	}
	return nil
}

// Find provides searhing of the data
func (r *Redis) Find(req *storage.FindConfig) ([]models.Parking, error) {
	return nil, nil
}

// createIndex provides creating of the index for searching data
func createIndex(client *redis.Client, index, parentID string) error {
	if err := client.HSet(index, "field", parentID).Err(); err != nil {
		return fmt.Errorf("unable to create index: %v", err)
	}

	return nil
}
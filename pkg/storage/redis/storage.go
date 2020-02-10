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

// Redis provides implementation of the storage for redis
type Redis struct {
	client *redis.Client
}

// New provides initialization of the redis client
func New(c *config.Config, client *redis.Client) (storage.Storage, error) {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     c.StorageAddress,
			Password: c.StoragePassword,
			DB:       0,
		})
	}

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping redis")
	}

	return &Redis{
		client: client,
	}, nil
}

// Insert provides inserting of slisw with parking data
func (r *Redis) Insert(data []models.Parking) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	modes := make(map[string][]interface{})
	for _, d := range data {
		if d.ID == 0 {
			return 0, errNoID
		}
		key := fmt.Sprintf("id_%d", d.ID)
		result, err := json.Marshal(d)
		if err != nil {
			return 0, fmt.Errorf("unable to marshal data with id %v: %v", d.ID, err)
		}
		err = r.client.Do("SET", key, string(result)).Err()
		if err != nil {
			return 0, fmt.Errorf("unable to set data: %v", err)
		}
		if err := createIndex(r.client, fmt.Sprintf("global_id_%d", d.GlobalID), fmt.Sprintf("%d", d.ID)); err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("unable to create index: %v", err))
		}
		modes[getMD5Hash(d.Mode)] = append(modes[getMD5Hash(d.Mode)], d.ID)
	}

	return len(data), createOneToMany(r.client, "mode", modes)
}

// Find provides searhing of the data
func (r *Redis) Find(req *storage.FindConfig) ([]models.Parking, error) {
	if req == nil {
		return nil, nil
	}
	var (
		key string
		err error
	)
	if req.ModeID != "" {
		return findOneToMany(r.client, fmt.Sprintf("mode_%s", getMD5Hash(req.ModeID)))
	}

	if req.GlobalID != 0 {
		key, err = findByKeys(r.client, fmt.Sprintf("global_id_%d", req.GlobalID))
		if err != nil {
			return nil, errors.Wrap(err, "unable to find by the key")
		}
	}

	if req.ID != 0 {
		key = fmt.Sprintf("%d", req.ID)
	}
	obj, err := getObject(r.client, key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get object")
	}
	return []models.Parking{obj}, nil
}

// findByKeys provides searching by additional indexes
func findByKeys(conn *redis.Client, key string) (string, error) {
	data, err := conn.HGet(key, "field").Result()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to get key %s", key))
	}
	return data, nil
}

// findOneToMany provides searching of data by one to many relationships
// "one-to-many" implemented by Sets
func findOneToMany(conn *redis.Client, key string) ([]models.Parking, error) {
	members, err := conn.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	response := []models.Parking{}
	for _, id := range members {
		obj, err := getObject(conn, id)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get object")
		}
		response = append(response, obj)
	}
	return response, nil
}

// createOneToMany provides creating one to many relationships
// Its implemented by Sets
func createOneToMany(conn *redis.Client, key string, data map[string][]interface{}) error {
	for k, v := range data {
		if err := conn.SAdd(fmt.Sprintf("%s_%s", key, k), v...).Err(); err != nil {
			return errors.Wrap(err, "unable to create data")
		}
	}
	return nil
}

// createIndex provides creating of the index for searching data
func createIndex(client *redis.Client, index, parentID string) error {
	if err := client.HSet(index, "field", parentID).Err(); err != nil {
		return fmt.Errorf("unable to create index: %v", err)
	}

	return nil
}

func getObject(conn *redis.Client, key string) (models.Parking, error) {
	parking := models.Parking{}
	objStr, err := conn.Do("GET", fmt.Sprintf("id_%s", key)).String()
	if err != nil {
		return parking, errors.Wrap(err, "unable to find by the key")
	}
	b := []byte(objStr)
	err = json.Unmarshal(b, &parking)
	if err != nil {
		return parking, errors.Wrap(err, "unable to find by the key")
	}

	return parking, nil
}

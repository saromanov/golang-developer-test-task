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

	modes := make(map[string][]interface{})

	for _, d := range data {
		if d.ID == 0 {
			return errNoID
		}
		key := fmt.Sprintf("id_%d", d.ID)
		result, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("unable to marshal data: %v", err)
		}
		err = r.client.Do("SET", key, string(result)).Err()
		if err != nil {
			return fmt.Errorf("unable to set data: %v", err)
		}
		if err := createIndex(r.client, fmt.Sprintf("global_id_%d", d.GlobalID), key); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to create index: %v", err))
		}
		modes[d.Mode] = append(modes[d.Mode], d.ID)
	}

	return createOneToMany(r.client, "mode", modes)
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
	if req.GlobalID != 0 {
		key, err = findByKeys(fmt.Sprintf("global_id_%d", req.GlobalID))
		if err != nil {
			return nil, errors.Wrap(err, "unable to find by the key")
		}
	}
	if req.ModeID != "" {
		fmt.Println("SS: ", fmt.Sprintf("mode_%x", req.ModeID))
		return findOneToMany(r.client, fmt.Sprintf("mode_%x", req.ModeID))
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
func findByKeys(key string) (string, error) {
	return "", nil
}

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
		fmt.Println("M: ", fmt.Sprintf("%s_%s", key, k))
		if err := conn.SAdd(fmt.Sprintf("%s_%s", key, "d0bad180d183d0b3d0bbd0bed181d183d182d0bed187d0bdd0be"), v...).Err(); err != nil {
			return errors.Wrap(err, "unable to create data")
		}
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

// createIndex provides creating of the index for searching data
func createIndex(client *redis.Client, index, parentID string) error {
	if err := client.HSet(index, "field", parentID).Err(); err != nil {
		return fmt.Errorf("unable to create index: %v", err)
	}

	return nil
}

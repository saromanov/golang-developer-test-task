package redis


import (
	"testing"
	"github.com/go-redis/redis"
	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock"
	"errors"
	"github.com/stretchr/testify/assert"
)

func newTestRedis() *redismock.ClientMock {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return redismock.NewNiceMock(client)
}

func RedisIsAvailable(client redis.Cmdable) bool {
	return client.Ping().Err() == nil
}

func TestRedisCannotBePinged(t *testing.T) {
	r := newTestRedis()
	r.On("Ping").
		Return(redis.NewStatusResult("", errors.New("server not available")))

	assert.False(t, RedisIsAvailable(r))
}

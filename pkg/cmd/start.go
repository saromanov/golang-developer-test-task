package cmd

import (
	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/data"
	"github.com/saromanov/golang-developer-test-task/pkg/logger"
	"github.com/saromanov/golang-developer-test-task/pkg/server"
	"github.com/saromanov/golang-developer-test-task/pkg/storage/redis"
)

// Start provides starting of the app
func Start() {
	log := logger.New()
	c := &config.Config{
		Address: "localhost:3000",
		Logger:  log,
	}
	storage, err := redis.New(c)
	if err != nil {
		log.Fatalf("unable to init storage redis: %v", err)
	}
	d, err := data.Load("data.json")
	if err != nil {
		log.Fatalf("unable to load data: %v", err)
	}
	if err := storage.Insert(d); err != nil {
		log.Fatalf("unable to save data to the storage: %v", err)
	}
	if err := server.Make(storage, c); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

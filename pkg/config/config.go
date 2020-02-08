package config

import "github.com/saromanov/golang-developer-test-task/pkg/logger"

// Config provides definition of the config
type Config struct {
	Address         string
	StorageAddress  string
	StoragePassword string
	Logger          *logger.Logger
}

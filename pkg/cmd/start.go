package cmd

import (
	"os"

	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/loader/local"
	"github.com/saromanov/golang-developer-test-task/pkg/logger"
	"github.com/saromanov/golang-developer-test-task/pkg/server"
	"github.com/saromanov/golang-developer-test-task/pkg/storage/redis"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "golang-developer-test-task",
	Usage: "create puppet for the project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "storage-address",
			Value: "redis:6379",
			Usage: "address to storage connection",
		},
		&cli.StringFlag{
			Name:  "storage-password",
			Value: "",
			Usage: "storage password",
		},
		&cli.StringFlag{
			Name:  "address",
			Value: ":3000",
			Usage: "address to web",
		},
		&cli.StringFlag{
			Name:  "path-to-data",
			Value: "./assets/data.json",
			Usage: "path to loading of data",
		},
	},
	Action: initialize,
}

// Start provides starting of the app
func Start(args []string) {
	err := app.Run(os.Args)
	if err != nil {
		return
	}
}

// initialize provides initialization of the app
func initialize(ctx *cli.Context) error {
	log := logger.New()
	c := makeConfig(ctx)
	c.Logger = log
	storage, err := redis.New(c, nil)
	if err != nil {
		log.Fatalf("unable to init storage redis: %v", err)
	}
	redis.InitPrometheus()
	dLocal, err := local.New(ctx.String("path-to-data"))
	if err != nil {
		log.Fatalf("unable to load data: %v", err)
	}
	d, err := dLocal.Load()
	if err != nil {
		log.Fatalf("unable to load data: %v", err)
	}
	if _, err := storage.Insert(d); err != nil {
		log.Fatalf("unable to save data to the storage: %v", err)
	}
	if err := server.Make(storage, c); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	return nil
}

// makeConfig provides construction of the config data
// from input. Its supports flags or environment variables
func makeConfig(ctx *cli.Context) *config.Config {
	c := &config.Config{}
	address := os.Getenv("ADDRESS")
	if address != "" {
		c.Address = address
	} else {
		c.Address = ctx.String("address")
	}
	storageAddress := os.Getenv("STORAGE_ADDRESS")
	if storageAddress != "" {
		c.StorageAddress = storageAddress
	} else {
		c.StorageAddress = ctx.String("storage-address")
	}
	storagePassword := os.Getenv("STORAGE_PASSWORD")
	if storagePassword != "" {
		c.StoragePassword = storagePassword
	} else {
		c.StoragePassword = ctx.String("address-password")
	}

	return c
}

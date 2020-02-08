package cmd

import (
	"os"

	"github.com/saromanov/golang-developer-test-task/pkg/config"
	"github.com/saromanov/golang-developer-test-task/pkg/data"
	"github.com/saromanov/golang-developer-test-task/pkg/logger"
	"github.com/saromanov/golang-developer-test-task/pkg/server"
	"github.com/saromanov/golang-developer-test-task/pkg/storage/redis"
	"github.com/urfave/cli/v2"
)

// Start provides starting of the app
func Start(args []string) {
	app := &cli.App{
		Name:  "golang-developer-test-task",
		Usage: "create puppet for the project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "storage-address",
				Value: "localhost:6379",
				Usage: "address to storage connection",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "localhost:3000",
				Usage: "address to web",
			},
		},
		Action: initialize,
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}

// initialize provides initialization of the app
func initialize(ctx *cli.Context) error {
	log := logger.New()
	c := &config.Config{
		Address:        ctx.String("address"),
		Logger:         log,
		StorageAddress: ctx.String("storage-address"),
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
	return nil
}

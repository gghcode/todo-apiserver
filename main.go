package main

import (
	"apas-todo-apiserver/app"
	"apas-todo-apiserver/config"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func main() {
	logger := logrus.New()

	builder := config.NewViperBuilder()
	builder.SetBasePath(".")
	builder.AddJsonFile("config")
	builder.AddEnvironmentVariables()

	configuration, err := builder.Build()
	if err != nil {
		logger.Fatalln(errors.Wrap(err, "Configuration build failed."))
	}

	server := app.NewServer(configuration)
	server.Initialize()

	if err := server.Run(); err != nil {
		logger.Fatalln(err)
	}
}

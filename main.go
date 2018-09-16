package main

import (
	"apas-todo-apiserver/app"
	"apas-todo-apiserver/config"
	"apas-todo-apiserver/todo"
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

	controllers := []app.ApiController{
		todo.NewController(configuration),
	}

	server := app.NewServer(configuration, controllers)
	server.Initialize()

	if err := server.Run(); err != nil {
		logger.Fatalln(err)
	}
}

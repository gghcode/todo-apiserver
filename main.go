package main

import (
	"apas-todo-apiserver/app"
	logger "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func main() {
	builder := app.NewViperConfigurationBuilder()
	builder.SetBasePath(".")
	builder.AddJsonFile("config")
	builder.AddEnvironmentVariables()

	configuration, err := builder.Build()
	if err != nil {
		logger.Fatalln(errors.Wrap(err, "Configuration build failed."))
	}

	if err := app.NewServerAndRun(*configuration); err != nil {
		logger.Fatalln(err)
	}
}

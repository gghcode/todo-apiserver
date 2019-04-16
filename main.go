package main

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

const (
	envPrefix = "REST"
)

func main() {
	conf, err := config.NewBuilder().
		AddConfigFile("config.yaml").
		BindEnvs(envPrefix).
		Build()

	if err != nil {
		panic(err)
	}

	server := app.New(conf)
	if err := server.Run(); err != nil {
		panic(err)
	}
}

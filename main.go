package main

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

const (
	envPrefix = "REST"
)

func main() {
	_, err := config.NewBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()

	if err != nil {
		panic(err)
	}
}

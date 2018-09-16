package app

import (
	"apas-todo-apiserver/config"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type TodoApiServer struct {
	engine        ServerEngine
	configuration config.Configuration
	controllers   []ApiController
	logger        *logrus.Entry
}

func NewServer(configuration config.Configuration, controllers []ApiController) *TodoApiServer {
	logger := logrus.New().WithField("host", "server")

	return &TodoApiServer{
		configuration: configuration,
		controllers:   controllers,
		logger:        logger,
	}
}

func (apiServer *TodoApiServer) Initialize() {
	engine := NewGinEngine()

	RegisterControllers(engine, apiServer.controllers)

	apiServer.engine = engine
}

func (apiServer *TodoApiServer) Run() error {
	listenPort := apiServer.configuration.ListenPort
	listenAddr := getAddrString(listenPort)

	return apiServer.engine.Run(listenAddr)
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

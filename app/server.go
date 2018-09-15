package app

import (
	"apas-todo-apiserver/controllers"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type TodoApiServer struct {
	engine        ServerEngine
	configuration Configuration
	logger        *logrus.Entry
}

func NewServer(configuration Configuration) *TodoApiServer {
	logger := logrus.New().WithField("host", "server")

	return &TodoApiServer{
		configuration: configuration,
		logger:        logger,
	}
}

func (apiServer *TodoApiServer) Initialize() {
	engine := NewGinEngine()

	setRoutes(engine, getControllerList())

	apiServer.engine = engine
}

func getControllerList() []controllers.Controller {
	return []controllers.Controller{
		&controllers.TodoController{},
	}
}

func setRoutes(router ServerEngine, controllerList []controllers.Controller) {
	for _, controller := range controllerList {
		addRoutes(router, controller)
	}
}

func addRoutes(router ServerEngine, controller controllers.Controller) {
	handlers := controller.GetHandlerInfos()

	for _, handler := range handlers {
		router.Handle(handler)
	}
}

func (apiServer *TodoApiServer) Run() error {
	return apiServer.engine.Run(getAddrString(apiServer.configuration.ListenPort))
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

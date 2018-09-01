package app

import (
	"apas-todo-apiserver/controllers"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type TodoApiServer struct {
	engine        *gin.Engine
	logger        *logrus.Entry
	configuration Configuration
}

func NewServer(configuration Configuration) *TodoApiServer {
	logger := logrus.New().WithField("host", "server")

	return &TodoApiServer{
		configuration: configuration,
		logger:        logger,
	}
}

func (apiServer *TodoApiServer) Initialize() {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	setRoutes(engine, getControllerList())

	apiServer.engine = engine
}

func getControllerList() []controllers.Controller {
	return []controllers.Controller{
		&controllers.TodoController{},
	}
}

func setRoutes(router *gin.Engine, controllerList []controllers.Controller) {
	for _, controller := range controllerList {
		addRoutes(router, controller)
	}
}

func addRoutes(router *gin.Engine, controller controllers.Controller) {
	handlers := controller.GetHandlers()

	for _, handler := range handlers {
		router.Handle(handler.Method, handler.Path, handler.Handle)
	}
}

func (apiServer *TodoApiServer) Run() error {
	return apiServer.engine.Run(getAddrString(apiServer.configuration.ListenPort))
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

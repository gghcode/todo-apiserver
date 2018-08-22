package app

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type TodoApiServer struct {
	configuration Configuration
	logger        logrus.Entry
}

func NewServerAndRun(configuration Configuration) error {
	return RunServer(NewTodoApiServer(configuration))
}

func NewTodoApiServer(configuration Configuration) *TodoApiServer {
	logger := logrus.New().WithField("host", "server")

	return &TodoApiServer{
		configuration: configuration,
		logger:        *logger,
	}
}

func RunServer(apiServer *TodoApiServer) error {
	router := gin.Default()

	return router.Run(getAddrString(apiServer.configuration.ListenPort))
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

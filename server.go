package main

import (
	"apas-todo-service/config"
	"apas-todo-service/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(configuration.MongoConnStr)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	controllerList := createAllControllers(session)

	router := gin.Default()
	router.Group("/api")
	{
		for _, controller := range controllerList {
			controller.InitializeRoutes(router)
		}
	}

	router.Run(fmt.Sprintf(":%d", configuration.Port))
}

func createAllControllers(session *mgo.Session) []controllers.Controller {
	var result []controllers.Controller

	result = append(result, controllers.CreateTodoController(session))

	return result
}

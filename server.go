package main

import (
	"apas-todo-service/app"
	"apas-todo-service/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func main() {
	environment, err := app.LoadEnvironment()
	if err != nil {
		panic(err)
	}

	db, err := mgo.Dial(environment.Config.MongoConnStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	controllerList := createAllControllers(db)
	router := gin.Default()

	for _, controller := range controllerList {
		controller.InitializeRoutes(router)
	}

	router.Run(getAddrString(environment.Config.Port))
}

func createAllControllers(db *mgo.Session) []controllers.Controller {
	var result []controllers.Controller

	result = append(result, controllers.CreateTodoController(db))

	return result
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

package main

import (
	"apas-todo-service/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func main() {
	configuration := config.Load()

	session, err := mgo.Dial(configuration.MongoConnStr)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	router := gin.Default()
	router.Run(fmt.Sprintf(":%d", configuration.Port))
}

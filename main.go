package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/gghcode/apas-todo-apiserver/docs"
)

// @title APAS TODO API
// @version 1.0
// @description This is a apas todo api server.
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	srv, cleanupFunc, err := InitializeRouter()
	defer cleanupFunc()

	if err != nil {
		panic(err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

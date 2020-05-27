package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gghcode/apas-todo-apiserver/config"
	_ "github.com/gghcode/apas-todo-apiserver/docs"
)

const (
	envPrefix = "REST"
)

// @title APAS TODO API
// @version 1.0
// @description This is a apas todo api server.
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewViperBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()

	if err != nil {
		panic(err)
	}

	ginRouter, cleanupFunc, err := InitializeRouter(cfg)
	defer cleanupFunc()

	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: ginRouter,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown: ", err)
	}

	fmt.Println("Shutdown was successful")
}

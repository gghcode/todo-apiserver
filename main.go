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
	"github.com/gghcode/apas-todo-apiserver/db"
	_ "github.com/gghcode/apas-todo-apiserver/docs"
	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	infraSecurity "github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webApp "github.com/gghcode/apas-todo-apiserver/web/api/app"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"
	webTodo "github.com/gghcode/apas-todo-apiserver/web/api/todo"
	webUser "github.com/gghcode/apas-todo-apiserver/web/api/user"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/defval/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
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

	var ginRouter *gin.Engine

	c := setupContainer(cfg)
	if err := c.Extract(&ginRouter); err != nil {
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

	disposeDBConnections(c)

	fmt.Println("Shutdown was successful")
}

func setupContainer(cfg config.Configuration) *inject.Container {
	provideConfig := func() config.Configuration {
		return cfg
	}

	container, _ := inject.New(
		inject.Provide(provideConfig),

		inject.Provide(db.NewPostgresConn),
		inject.Provide(db.NewRedisConn),
		inject.Provide(infraSecurity.NewBcryptPassport),

		inject.Provide(afero.NewOsFs),
		inject.Provide(file.NewAferoFileReader),
		inject.Provide(app.NewService),
		inject.Provide(webApp.NewController, injectAsGinController()),

		inject.Provide(repository.NewRedisTokenRepository),
		inject.Provide(repository.NewUserRepository),
		inject.Provide(auth.NewService),
		inject.Provide(jwt.NewJwtAccessTokenGeneratorFunc),
		inject.Provide(jwt.NewJwtRefreshTokenGeneratorFunc),
		inject.Provide(webAuth.NewController, injectAsGinController()),

		inject.Provide(user.NewService),
		inject.Provide(webUser.NewController, injectAsGinController()),

		inject.Provide(repository.NewGormTodoRepository),
		inject.Provide(todo.NewTodoService),
		inject.Provide(webTodo.NewController, injectAsGinController()),

		inject.Provide(newGinRouter),
	)

	return container
}

func injectAsGinController() inject.ProvideOption {
	return inject.As(api.GinControllerToken)
}

func newGinRouter(cfg config.Configuration, controllers []api.GinController) *gin.Engine {
	router := gin.New()
	registerMiddlewares(cfg, router)

	for _, c := range controllers {
		c.RegisterRoutes(router.Group(""))
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func registerMiddlewares(cfg config.Configuration, router gin.IRouter) {
	router.Use(middleware.AddAccessTokenHandler(
		jwt.NewJwtAccessTokenVerifyHandlerFactory(cfg),
	))
}

func disposeDBConnections(c *inject.Container) error {
	var gormConn db.GormConnection
	var redisConn db.RedisConnection

	if err := c.Extract(&gormConn); err != nil {
		return err 
	}

	if err := c.Extract(&redisConn); err != nil {
		return err
	}
	
	gormConn.Close()
	redisConn.Close()

	return nil
}

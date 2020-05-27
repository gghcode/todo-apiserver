//+build wireinject

package main

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webApp "github.com/gghcode/apas-todo-apiserver/web/api/app"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"
	webTodo "github.com/gghcode/apas-todo-apiserver/web/api/todo"
	webUser "github.com/gghcode/apas-todo-apiserver/web/api/user"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/afero"
)

var dbSet = wire.NewSet(
	db.NewPostgresConn,
)

var redisSet = wire.NewSet(
	db.NewRedisConn,
)

var todoSet = wire.NewSet(
	repository.NewGormTodoRepository,
	todo.NewTodoService,
	webTodo.NewController,
)

var securitySet = wire.NewSet(
	security.NewBcryptPassport,
)

var authSet = wire.NewSet(
	repository.NewRedisTokenRepository,
	jwt.NewJwtAccessTokenGeneratorFunc,
	jwt.NewJwtRefreshTokenGeneratorFunc,
	auth.NewService,
	webAuth.NewController,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	user.NewService,
	webUser.NewController,
)

var appSet = wire.NewSet(
	afero.NewOsFs,
	file.NewAferoFileReader,
	app.NewService,
	webApp.NewController,
)

var routerSet = wire.NewSet(
	provideControllers,
	newGinRouter,
)

func InitializeRouter(cfg config.Configuration) (*gin.Engine, func(), error) {
	wire.Build(
		dbSet,
		redisSet,
		securitySet,
		todoSet,
		authSet,
		userSet,
		appSet,
		routerSet,
	)
	return nil, nil, nil
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

func provideControllers(
	appController *webApp.Controller,
	todoController *webTodo.Controller,
	authController *webAuth.Controller,
	userController *webUser.Controller,
) []api.GinController {
	return []api.GinController{
		appController,
		todoController,
		authController,
		userController,
	}
}

func registerMiddlewares(cfg config.Configuration, router gin.IRouter) {
	router.Use(middleware.AddAccessTokenHandler(
		jwt.NewJwtAccessTokenVerifyHandlerFactory(cfg),
	))
}

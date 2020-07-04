//+build wireinject

package main

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/app"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/gghcode/apas-todo-apiserver/web"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webApp "github.com/gghcode/apas-todo-apiserver/web/api/app"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"
	webTodo "github.com/gghcode/apas-todo-apiserver/web/api/todo"
	webUser "github.com/gghcode/apas-todo-apiserver/web/api/user"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/afero"
)

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

func provideMiddlewares(
	accessTokenHandlerMiddleware middleware.AccessTokenHandlerMiddleware,
	corsMiddleware middleware.CorsMiddleware,
) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.HandlerFunc(accessTokenHandlerMiddleware),
		gin.HandlerFunc(corsMiddleware),
	}
}

var configSet = wire.NewSet(
	config.FromEnvs,
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
	jwt.NewJwtAccessTokenVerifyHandlerFactory,
	middleware.NewAccessTokenHandler,
	middleware.NewCors,
	provideMiddlewares,
	provideControllers,
	web.NewGinRouter,
)

func InitializeRouter() (*http.Server, func(), error) {
	wire.Build(
		configSet,
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

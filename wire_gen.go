// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/app"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/infra/bcrypt"
	"github.com/gghcode/apas-todo-apiserver/infra/file"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/repository"
	"github.com/gghcode/apas-todo-apiserver/infra/jwt"
	"github.com/gghcode/apas-todo-apiserver/web"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	app2 "github.com/gghcode/apas-todo-apiserver/web/api/app"
	auth2 "github.com/gghcode/apas-todo-apiserver/web/api/auth"
	todo2 "github.com/gghcode/apas-todo-apiserver/web/api/todo"
	user2 "github.com/gghcode/apas-todo-apiserver/web/api/user"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/afero"
	"net/http"
)

import (
	_ "github.com/gghcode/apas-todo-apiserver/docs"
)

// Injectors from wire.go:

func InitializeRouter() (*http.Server, func(), error) {
	configuration, err := config.FromEnvs()
	if err != nil {
		return nil, nil, err
	}
	accessTokenHandlerFactory := jwt.NewJwtAccessTokenVerifyHandlerFactory(configuration)
	accessTokenHandlerMiddleware := middleware.NewAccessTokenHandler(accessTokenHandlerFactory)
	corsMiddleware := middleware.NewCors(configuration)
	v := provideMiddlewares(accessTokenHandlerMiddleware, corsMiddleware)
	fs := afero.NewOsFs()
	fileReader := file.NewAferoFileReader(fs)
	useCase := app.NewService(fileReader)
	gormConnection, cleanup, err := db.NewPostgresConn(configuration)
	if err != nil {
		return nil, nil, err
	}
	redisConnection, cleanup2 := db.NewRedisConn(configuration)
	controller := app2.NewController(useCase, gormConnection, redisConnection)
	todoRepository := repository.NewGormTodoRepository(gormConnection)
	todoUseCase := todo.NewTodoService(todoRepository)
	todoController := todo2.NewController(todoUseCase)
	passwordAuthenticator := bcrypt.NewPasswordAuthenticator()
	tokenRepository := repository.NewRedisTokenRepository(redisConnection)
	userRepository := repository.NewUserRepository(gormConnection)
	userDataSource := provideDataSource(userRepository)
	accessTokenGeneratorFunc := jwt.NewJwtAccessTokenGeneratorFunc(configuration)
	refreshTokenGeneratorFunc := jwt.NewJwtRefreshTokenGeneratorFunc(configuration)
	authUseCase := auth.NewService(configuration, passwordAuthenticator, tokenRepository, userDataSource, accessTokenGeneratorFunc, refreshTokenGeneratorFunc)
	authController := auth2.NewController(authUseCase)
	passwordEncryptor := bcrypt.NewPasswordEncryptor(configuration)
	userUseCase := user.NewService(userRepository, passwordEncryptor)
	userController := user2.NewController(userUseCase)
	v2 := provideControllers(controller, todoController, authController, userController)
	server, cleanup3 := web.NewGinRouter(configuration, v, v2)
	return server, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func provideControllers(
	appController *app2.Controller,
	todoController *todo2.Controller,
	authController *auth2.Controller,
	userController *user2.Controller,
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
	return []gin.HandlerFunc{gin.HandlerFunc(accessTokenHandlerMiddleware), gin.HandlerFunc(corsMiddleware)}
}

func provideDataSource(userRepo user.Repository) auth.UserDataSource {
	return userRepo
}

var configSet = wire.NewSet(config.FromEnvs)

var dbSet = wire.NewSet(db.NewPostgresConn)

var redisSet = wire.NewSet(db.NewRedisConn)

var todoSet = wire.NewSet(repository.NewGormTodoRepository, todo.NewTodoService, todo2.NewController)

var bcryptSet = wire.NewSet(bcrypt.NewPasswordAuthenticator, bcrypt.NewPasswordEncryptor)

var authSet = wire.NewSet(repository.NewRedisTokenRepository, jwt.NewJwtAccessTokenGeneratorFunc, jwt.NewJwtRefreshTokenGeneratorFunc, auth.NewService, auth2.NewController)

var userSet = wire.NewSet(repository.NewUserRepository, user.NewService, provideDataSource, user2.NewController)

var appSet = wire.NewSet(afero.NewOsFs, file.NewAferoFileReader, app.NewService, app2.NewController)

var routerSet = wire.NewSet(jwt.NewJwtAccessTokenVerifyHandlerFactory, middleware.NewAccessTokenHandler, middleware.NewCors, provideMiddlewares,
	provideControllers, web.NewGinRouter,
)

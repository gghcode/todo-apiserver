//+build wireinject

package main

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/app"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/infra/bcrypt"
	"github.com/gghcode/apas-todo-apiserver/infra/file"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/repository"
	"github.com/gghcode/apas-todo-apiserver/infra/jwt"
	"github.com/gghcode/apas-todo-apiserver/infra/redis"
	"github.com/gghcode/apas-todo-apiserver/infra/redis/repo"
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

func provideDataSource(userRepo user.Repository) auth.UserDataSource {
	return userRepo
}

var configSet = wire.NewSet(
	config.FromEnvs,
)

var dbSet = wire.NewSet(
	gorm.NewPostgresConn,
)

var redisSet = wire.NewSet(
	redis.NewConnection,
)

var todoSet = wire.NewSet(
	repository.NewGormTodoRepository,
	todo.NewTodoService,
	webTodo.NewController,
)

var bcryptSet = wire.NewSet(
	bcrypt.NewPasswordAuthenticator,
	bcrypt.NewPasswordEncryptor,
)

var authSet = wire.NewSet(
	repo.NewRedisTokenRepository,
	jwt.NewJwtAccessTokenGeneratorFunc,
	jwt.NewJwtRefreshTokenGeneratorFunc,
	auth.NewService,
	webAuth.NewController,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	user.NewService,
	provideDataSource,
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
		bcryptSet,
		todoSet,
		authSet,
		userSet,
		appSet,
		routerSet,
	)
	return nil, nil, nil
}

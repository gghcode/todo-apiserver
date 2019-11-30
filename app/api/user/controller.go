package user

import (
	"net/http"
	"strconv"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/infra"
	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gin-gonic/gin"
)

// Controller godoc
type Controller struct {
	userRepository Repository
	passport       infra.Passport
}

// NewController godoc
func NewController(userRepository Repository, passport infra.Passport) *Controller {
	return &Controller{
		userRepository: userRepository,
		passport:       passport,
	}
}

// RegisterRoutes godocs
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.POST("api/users", controller.CreateUser)
	router.GET("api/users", controller.UserByID)
	router.GET("api/users/:username", controller.UserByName)

	authorized := router.Use(middleware.JwtAuthRequired())
	{
		authorized.GET("api/user", controller.AuthenticatedUser)
	}
}

// CreateUser godoc
// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.CreateUserRequest true "user payload"
// @Success 201 {object} user.UserResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid user payload"
// @Failure 409 {object} api.ErrorResponse "Already exists user"
// @Tags User API
// @Router /api/users [post]
func (controller *Controller) CreateUser(ctx *gin.Context) {
	createUserRequestValidator := CreateUserRequestValidator{}
	if err := createUserRequestValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	reqPayload := createUserRequestValidator.Model

	passwordHash, _ := controller.passport.HashPassword(reqPayload.Password)
	createdUser, err := controller.userRepository.CreateUser(User{
		UserName:     reqPayload.UserName,
		PasswordHash: passwordHash,
	})

	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, createdUser.Response())
}

// UserByID godoc
// @Description Fetch user by user id
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} api.ErrorResponse "User Not Found"
// @Tags User API
// @Router /api/users [get]
func (controller *Controller) UserByID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if err != nil {
		api.WriteErrorResponse(ctx, ErrInvalidUserID)
		return
	}

	user, err := controller.userRepository.UserByID(userID)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

// UserByName godoc
// @Description Fetch user by username
// @Accept json
// @Produce json
// @Param username path string true "User Name"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} api.ErrorResponse "User Not Found"
// @Tags User API
// @Router /api/users/{username} [get]
func (controller *Controller) UserByName(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := controller.userRepository.UserByUserName(username)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

// AuthenticatedUser godoc
// @Description Fetch user itself by access_token
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} api.ErrorResponse "User Not Found"
// @Tags User API
// @Router /api/user [get]
func (controller *Controller) AuthenticatedUser(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	user, err := controller.userRepository.UserByID(userID)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

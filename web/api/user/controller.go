package user

import (
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller is user controller
type Controller struct {
	userService user.UsecaseInteractor
}

// NewController return user controller
func NewController(userService user.UsecaseInteractor) *Controller {
	return &Controller{
		userService: userService,
	}
}

// RegisterRoutes register handler routes.
func (c *Controller) RegisterRoutes(router gin.IRouter) {
	router.POST("api/users", c.CreateUser)

	authorized := router.Use(middleware.RequiredAccessToken())
	{
		authorized.GET("api/user", c.User)
	}
}

// CreateUser is api that create user
// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.createUserRequestDTO true "user payload"
// @Success 201 {object} user.userResponseDTO "ok"
// @Failure 400 {object} api.ErrorResponseDTO "Invalid user payload"
// @Failure 409 {object} api.ErrorResponseDTO "Already exists user"
// @Tags User API
// @Router /api/users [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var req user.CreateUserRequest
	if err := validateCreateUserRequestDTO(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponseDTO(err))
		return
	}

	res, err := c.userService.CreateUser(req)
	if err == user.ErrAlreadyExistUser {
		ctx.JSON(http.StatusConflict, api.MakeErrorResponseDTO(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	serializer := newUserResponseSerializer(res)

	ctx.JSON(http.StatusCreated, serializer.Response())
}

// User godoc
// @Description Fetch user itself by access_token
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} user.userResponseDTO "ok"
// @Failure 404 {object} api.ErrorResponseDTO "User Not Found"
// @Tags User API
// @Router /api/user [get]
func (c *Controller) User(ctx *gin.Context) {
	res, err := c.userService.GetUserByUserID(middleware.AuthUserID(ctx))
	if err == user.ErrUserNotFound {
		ctx.JSON(http.StatusNotFound, api.MakeErrorResponseDTO(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	serializer := newUserResponseSerializer(res)

	ctx.JSON(http.StatusOK, serializer.Response())
}

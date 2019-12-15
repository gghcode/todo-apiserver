package user

import (
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/web/api"
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
}

// CreateUser is api that create user
// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.createUserRequestDto true "user payload"
// @Success 201 {object} user.UserResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid user payload"
// @Failure 409 {object} api.ErrorResponse "Already exists user"
// @Tags User API
// @Router /api/users [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var req user.CreateUserRequest
	if err := validateCreateUserRequestDto(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponse(err))
		return
	}

	res, err := c.userService.CreateUser(req)
	if err == user.ErrAlreadyExistUser {
		ctx.JSON(http.StatusConflict, api.MakeErrorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponse(err))
		return
	}

	serializer := newUserResponseSerializer(res)

	ctx.JSON(http.StatusCreated, serializer.Response())
}

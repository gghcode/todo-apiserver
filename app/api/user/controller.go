package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
)

// Controller godoc
type Controller struct {
}

// NewController godoc
func NewController() *Controller {
	return &Controller{}
}

// RegisterRoutes godocs
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.POST("/users", controller.CreateUser)
	// router.GET("/healthy", controller.Healthy)
}

// CreateUser godoc
// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.CreateUserRequest true "user payload"
// @Success 201 {object} user.UserResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid user payload"
// @Tags User API
// @Router /users [post]
func (controller *Controller) CreateUser(ctx *gin.Context) {
	var reqPayload CreateUserRequest
	if err := ctx.ShouldBindJSON(&reqPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{})
}

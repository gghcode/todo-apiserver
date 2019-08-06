package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
)

// Controller godoc
type Controller struct {
	userRepository Repository
}

// NewController godoc
func NewController(userRepository Repository) *Controller {
	return &Controller{
		userRepository: userRepository,
	}
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

	createdUser, err := controller.userRepository.CreateUser(User{
		UserName:     reqPayload.UserName,
		PasswordHash: []byte(reqPayload.Password),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{
		ID:        createdUser.ID,
		UserName:  createdUser.UserName,
		CreatedAt: time.Unix(createdUser.CreatedAt, 0),
	})
}

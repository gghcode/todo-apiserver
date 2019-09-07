package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/infra"
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
	router.POST("/users", controller.CreateUser)
	// router.GET("/users/:id", controller.UserByID)
	router.GET("/users/:username", controller.UserByName)
}

// CreateUser godoc
// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.CreateUserRequest true "user payload"
// @Success 201 {object} user.UserResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid user payload"
// @Failure 409 {object} common.ErrorResponse "Already exists user"
// @Tags User API
// @Router /users [post]
func (controller *Controller) CreateUser(ctx *gin.Context) {
	var reqPayload CreateUserRequest
	if err := ctx.ShouldBindJSON(&reqPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	passwordHash, _ := controller.passport.HashPassword(reqPayload.Password)
	createdUser, err := controller.userRepository.CreateUser(User{
		UserName:     reqPayload.UserName,
		PasswordHash: passwordHash,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, createdUser.Response())
}

// UserByID godoc
// @Description Fetch user by user id
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} common.ErrorResponse "User Not Found"
// @Tags User API
// @Router /users/{id} [get]
func (controller *Controller) UserByID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		// ctx.JSON(http.StatusBadRequest, common.NewErrResp(common.))
		return
	}

	user, err := controller.userRepository.UserByID(userID)
	if err == ErrUserNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
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
// @Failure 404 {object} common.ErrorResponse "User Not Found"
// @Tags User API
// @Router /users/{username} [get]
func (controller *Controller) UserByName(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := controller.userRepository.UserByUserName(username)
	if err == ErrUserNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

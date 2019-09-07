package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/middleware"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/val"
)

// Controller godoc
type Controller struct {
	todoRepo Repository
}

// NewController godoc
func NewController(todoRepo Repository) *Controller {
	return &Controller{
		todoRepo: todoRepo,
	}
}

// RegisterRoutes godoc
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	routes := router.Group("/todos")
	{
		routes.GET("", controller.AllTodosByUserID)
	}

	authorized := routes.Use(middleware.JwtAuthRequired())
	{
		authorized.POST("", controller.AddTodo)
	}
}

// AddTodo godoc
// @Description Add new todo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body todo.AddTodoRequest true "payload"
// @Success 201 {object} todo.TodoResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid payload"
// @Tags Todo API
// @Router /todos [post]
func (controller *Controller) AddTodo(ctx *gin.Context) {
	todoValidator := NewAddTodoValidator()
	if err := todoValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	todoEntity := Todo{
		Title:      todoValidator.Model.Title,
		Contents:   todoValidator.Model.Contents,
		AssignorID: ctx.GetInt64(val.UserID),
	}

	todo, err := controller.todoRepo.AddTodo(todoEntity)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	todoSerializer := TodoSerializer{Model: todo}

	ctx.JSON(http.StatusCreated, todoSerializer.Response())
}

// AllTodosByUserID godoc
// @Description Fetch todos by user id
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} todo.TodoResponse "ok"
// @Failure 404 {object} api.ErrorResponse "User Not Found"
// @Tags Todo API
// @Router /todos [get]
func (controller *Controller) AllTodosByUserID(ctx *gin.Context) {
	userIDStr, exists := ctx.GetQuery("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if !exists || err != nil {
		api.WriteErrorResponse(ctx, user.ErrInvalidUserID)
		return
	}

	todos, err := controller.todoRepo.AllTodosByUserID(userID)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	todosSerializer := TodosSerializer{Model: todos}

	ctx.JSON(http.StatusOK, todosSerializer.Response())
}

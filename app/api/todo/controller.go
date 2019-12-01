package todo

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gin-gonic/gin"
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
	router.GET("api/todos", controller.AllTodosByUserID)

	authorized := router.Use(middleware.RequiredAccessToken())
	{
		authorized.POST("api/todos", controller.AddTodo)
		authorized.PATCH("api/todos/:todo_id", controller.UpdateTodoByTodoID)
		authorized.DELETE("api/todos/:todo_id", controller.RemoveTodoByTodoID)
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
// @Router /api/todos [post]
func (controller *Controller) AddTodo(ctx *gin.Context) {
	todoValidator := NewAddTodoValidator()
	if err := todoValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	todoEntity := Todo{
		Title:      todoValidator.Model.Title,
		Contents:   todoValidator.Model.Contents,
		AssignorID: ctx.GetInt64("user_id"),
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
// @Router /api/todos [get]
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

// UpdateTodoByTodoID godoc
// @Description Update todo by todo id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param todo_id path string true "Todo ID"
// @Param payload body todo.UpdateTodoRequest true "payload"
// @Success 200 {object} todo.TodoResponse "ok"
// @Failure 404 {object} api.ErrorResponse "Todo Not Found"
// @Tags Todo API
// @Router /api/todos/{todo_id} [patch]
func (controller *Controller) UpdateTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("todo_id")
	todoID = strings.Trim(todoID, " ")
	if len(todoID) <= 0 {
		api.WriteErrorResponse(ctx, ErrEmptyTodoID)
		return
	}

	todoValidator := NewUpdateTodoValidator()
	if err := todoValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	patchData := todoValidator.Model.Map()

	todo, err := controller.todoRepo.UpdateTodo(todoID, patchData)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	todoSerializer := TodoSerializer{Model: todo}

	ctx.JSON(http.StatusOK, todoSerializer.Response())
}

// RemoveTodoByTodoID godoc
// @Description Remove todo by todo id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param todo_id path string true "Todo ID"
// @Success 204 {object} todo.TodoResponse "ok"
// @Failure 404 {object} api.ErrorResponse "Todo Not Found"
// @Tags Todo API
// @Router /api/todos/{todo_id} [delete]
func (controller *Controller) RemoveTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("todo_id")
	todoID = strings.Trim(todoID, " ")
	if len(todoID) <= 0 {
		api.WriteErrorResponse(ctx, ErrEmptyTodoID)
		return
	}

	if err := controller.todoRepo.RemoveTodo(todoID); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

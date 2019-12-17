package todo

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"github.com/gin-gonic/gin"
)

// Controller is todo controller
type Controller struct {
	todoService todo.UsecaseInteractor
}

// NewController return todo controller
func NewController(todoService todo.UsecaseInteractor) *Controller {
	return &Controller{
		todoService: todoService,
	}
}

// RegisterRoutes register handler routes.
func (c *Controller) RegisterRoutes(router gin.IRouter) {
	authorized := router.Use(middleware.RequiredAccessToken())
	{
		authorized.POST("api/todos", c.AddTodo)
		authorized.GET("api/todos", c.Todos)
		authorized.DELETE("api/todos/:todo_id", c.RemoveTodoByTodoID)
	}
}

// AddTodo is api that add todo
// @Description Add new todo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body todo.addTodoRequestDTO true "payload"
// @Success 201 {object} todo.todoResponseDTO "ok"
// @Failure 400 {object} api.ErrorResponseDTO "Invalid payload"
// @Tags Todo API
// @Router /api/todos [post]
func (c *Controller) AddTodo(ctx *gin.Context) {
	var req todo.AddTodoRequest
	if err := validateAddTodoRequestDTO(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponseDTO(err))
		return
	}

	res, err := c.todoService.AddTodo(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	s := newTodoResponseSerializer(res)

	ctx.JSON(http.StatusCreated, s.Response())
}

// Todos fetch todos of authenticated user
// @Description Fetch todos of authenticated user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} todo.todoResponseDTO "ok"
// @Failure 404 {object} api.ErrorResponseDTO "User Not Found"
// @Tags Todo API
// @Router /api/todos [get]
func (c *Controller) Todos(ctx *gin.Context) {
	res, err := c.todoService.GetTodosByUserID(middleware.AuthUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	resArr := make([]todoResponseDTO, len(res))
	for i, t := range res {
		resArr[i] = newTodoResponseSerializer(t).Response()
	}

	ctx.JSON(http.StatusOK, resArr)
}

// RemoveTodoByTodoID godoc
// @Description Remove todo by todo id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param todo_id path string true "Todo ID"
// @Success 204 {object} todo.todoResponseDTO "ok"
// @Failure 404 {object} api.ErrorResponseDTO "Todo Not Found"
// @Tags Todo API
// @Router /api/todos/{todo_id} [delete]
func (c *Controller) RemoveTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("todo_id")

	if err := c.todoService.RemoveTodo(todoID); err == todo.ErrNotFoundTodo {
		ctx.JSON(http.StatusNotFound, api.MakeErrorResponseDTO(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

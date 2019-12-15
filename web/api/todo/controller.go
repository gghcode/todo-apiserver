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
	}
}

// AddTodo is api that add todo
// @Description Add new todo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body todo.AddTodoRequest true "payload"
// @Success 201 {object} todo.TodoResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid payload"
// @Tags Todo API
// @Router /api/todos [post]
func (c *Controller) AddTodo(ctx *gin.Context) {
	var req todo.AddTodoRequest
	if err := validateAddTodoRequestDto(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponse(err))
		return
	}

	res, err := c.todoService.AddTodo(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponse(err))
		return
	}

	s := newTodoResponseSerializer(res)

	ctx.JSON(http.StatusCreated, s.Response())
}

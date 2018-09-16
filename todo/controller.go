package todo

import (
	"apas-todo-apiserver/app"
	"apas-todo-apiserver/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	repository *Repository
}

func NewController(configuration config.Configuration) *Controller {
	controller := Controller{
		repository: NewRepository(configuration),
	}

	return &controller
}

func (controller *Controller) RouteInfos() app.RouteInfos {
	return app.RouteInfos{
		app.Route("GET", "api/v1/todos/:id", controller.todo),
		app.Route("POST", "api/v1/todos", controller.addTodo),
		app.Route("PUT", "api/v1/todos/:id", controller.updateTodo),
		app.Route("DELETE", "api/v1/todos/:id", controller.removeTodo),
	}
}

func (controller *Controller) todo(ctx *gin.Context) {
	todoId := ctx.Param("id")

	todo, err := controller.repository.Todo(todoId)
	if err != nil {
		ctx.AbortWithError(400, err)
	}

	ctx.JSON(200, todo)
}

func (controller *Controller) addTodo(ctx *gin.Context) {
	var todo Todo

	if err := ctx.Bind(&todo); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	todo, err := controller.repository.AddTodo(todo)
	if err != nil {
		ctx.AbortWithError(400, err)
	}

	ctx.JSON(http.StatusCreated, todo)
}

func (controller *Controller) updateTodo(ctx *gin.Context) {

}

func (controller *Controller) removeTodo(ctx *gin.Context) {
	todoId := ctx.Param("id")

	err := controller.repository.RemoveTodo(todoId)
	if err != nil {
		ctx.AbortWithError(400, err)
	}

	ctx.Status(http.StatusNoContent)
}

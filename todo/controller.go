package todo

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"net/http"
)

type Controller struct {
	repository *Repository
	logger *logrus.Entry
}

func NewController(configuration config.Configuration) *Controller {
	controller := Controller{
		repository: NewRepository(configuration),
		logger: logrus.New().WithField("caller", "todo.Controller"),
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

		controller.logger.Warn(err)
		ctx.JSON(500, err)
		//ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, todo)
}

func (controller *Controller) addTodo(ctx *gin.Context) {
	var todo Todo

	if err := ctx.Bind(&todo); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	todo, err := controller.repository.AddTodo(todo)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
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
		return
	}

	ctx.Status(http.StatusNoContent)
}

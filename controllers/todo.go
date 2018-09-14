package controllers

import (
	"apas-todo-apiserver/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"net/http"
)

type TodoController struct {
	//logger
}

func NewTodoController(session *mgo.Session) *TodoController {
	result := TodoController{}

	return &result
}

func (controller *TodoController) GetHandlerInfos() []HandlerInfo {
	return []HandlerInfo{
		HandlerInfo{
			Method: "GET",
			Path: "api/v1/todos",
			Handle: controller.getTodos,
		},
		HandlerInfo{
			Method: "GET",
			Path: "api/v1/todos/:id",
			Handle: controller.getTodo,
		},
		HandlerInfo{
			Method: "POST",
			Path: "api/v1/todos",
			Handle: controller.addTodo,
		},
		HandlerInfo{
			Method: "PUT",
			Path: "api/v1/todos/:id",
			Handle: controller.updateTodo,
		},
		HandlerInfo{
			Method: "DELETE",
			Path: "api/v1/todos/:id",
			Handle: controller.removeTodo,
		},
	}
}

func (controller *TodoController) getTodos(ctx *gin.Context) {

}

func (controller *TodoController) getTodo(ctx *gin.Context) {
	todo := models.Todo{}

	ctx.JSON(200, todo)
}

func (controller *TodoController) addTodo(ctx *gin.Context) {
	var todo models.Todo

	if err := ctx.Bind(&todo); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}


	ctx.JSON(http.StatusCreated, todo)
}

func (controller *TodoController) updateTodo(ctx *gin.Context) {

}

func (controller *TodoController) removeTodo(ctx *gin.Context) {

}
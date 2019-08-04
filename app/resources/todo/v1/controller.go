package v1

import "github.com/gin-gonic/gin"

// Controller provides request handler.
type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (controller Controller) getAllHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (controller Controller) getByIdHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (controller Controller) createHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (controller Controller) removeByIdHandler(ctx *gin.Context) {
	panic("Not implement...")
}

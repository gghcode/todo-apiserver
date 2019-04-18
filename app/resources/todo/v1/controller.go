package v1

import "github.com/gin-gonic/gin"

// Controller provides request handler.
type controller struct {
}

func NewController *controller {
	return &controller{
		
	}
}

func (controller controller) getHandler(ctx *gin.Context) {
	ctx.String(200, "string", "Hello")
}

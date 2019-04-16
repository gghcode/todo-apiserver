package v1

import "github.com/gin-gonic/gin"

// Controller provides request handler.
type Controller struct {
}

func (controller Controller) getHandler(ctx *gin.Context) {
	ctx.String(200, "string", "Hello")
}

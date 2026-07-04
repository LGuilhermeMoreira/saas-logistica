package rest

import (
	"github.com/gin-gonic/gin"
)

func NewController() *Controller {
	return &Controller{}
}

type Controller struct{}

func (c Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}

	// todo

}

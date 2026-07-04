package main

import (
	"gateway/internal/transport/rest"

	"github.com/gin-gonic/gin"
)

func main() {

	// http
	//
	mux := gin.Default()
	ctrl := rest.NewController()
	mux.POST("/login", ctrl.Login)

	go func() {

	}()
}

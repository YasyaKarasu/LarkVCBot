package app

import (
	"LarkVCBot/app/controller"
	"LarkVCBot/app/dispatcher"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// register your controllers here
	// example
	r.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
}

func Init(r *gin.Engine) {
	controller.InitEvent()
	Register(r)
}

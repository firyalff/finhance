package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/version", getVersionHandler)
}

func getVersionHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"version": appVersion,
	})
}
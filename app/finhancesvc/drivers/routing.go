package drivers

import (
	"github.com/gin-gonic/gin"
)

type HTTPMethod int32

const (
	HTTPGET HTTPMethod = iota
	HTTPPOST
	HTTPPUT
	HTTPDELETE
)

type AppRoute struct {
	Path    string
	Method  HTTPMethod
	Handler func(ctx *gin.Context)
}

func InitRouting() (router *gin.Engine) {
	router = gin.Default()
	return router
}

func StartRouteServer(router *gin.Engine) (err error) {
	return router.Run()
}

func RegisterRoutes(router *gin.Engine, appRoutes []AppRoute) {
	for _, appRoute := range appRoutes {
		switch appRoute.Method {
		case HTTPGET:
			router.GET(appRoute.Path, appRoute.Handler)
		}
	}
}

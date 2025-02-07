package api

import (
	"template/internal/common"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method  string
	Path    string
	Handler []gin.HandlerFunc
}

func Setup(routes ...Route) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Log())
	router.Use(Errors())

	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler...)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.Error(common.NotFoundError{Message: "nothing to do here"})
	})

	return router
}

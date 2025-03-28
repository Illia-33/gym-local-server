package main

import (
	"gymlocalserver/server/handlers"

	"github.com/gin-gonic/gin"
)

func wrap(f handlers.HandlerFunc, ctx *handlers.Context) gin.HandlerFunc {
	return func(r *gin.Context) {
		f(r, ctx)
	}
}

func registerApi(r *gin.Engine, ctx *handlers.Context) {
	r.GET("/cameras", wrap(handlers.GetCamerasInfo, ctx))

	r.POST("/cameras/:id/ptz", wrap(handlers.StartPtz, ctx))
	r.DELETE("/cameras/:id/ptz", wrap(handlers.EndPtz, ctx))
}

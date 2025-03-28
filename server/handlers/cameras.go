package handlers

import (
	"gymlocalserver/server/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func doGetCamerasInfo(ctx *Context, _ api.CamerasRequest) api.CamerasResponse {
	return api.CamerasResponse{
		Count: len(ctx.Cameras),
	}
}

func GetCamerasInfo(r *gin.Context, ctx *Context) {
	r.IndentedJSON(http.StatusOK, doGetCamerasInfo(ctx, api.CamerasRequest{}))
}

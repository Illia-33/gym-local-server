package handlers

import (
	"gymlocalserver/server/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func doGetCamerasInfo(ctx *Context, _ api.CamerasRequest) api.CamerasResponse {
	cameras := make([]api.CameraDescription, 0, len(ctx.Cameras))
	for _, camera := range ctx.Cameras {
		cameras = append(cameras, api.CameraDescription{
			Label:       camera.Config.Label,
			Description: camera.Config.Description,
		})
	}

	return api.CamerasResponse{
		Cameras: cameras,
	}
}

func GetCamerasInfo(r *gin.Context, ctx *Context) {
	r.IndentedJSON(http.StatusOK, doGetCamerasInfo(ctx, api.CamerasRequest{}))
}

package handlers

import (
	"net/http"

	api "github.com/Illia-33/gym-localserver/api/localserver"

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

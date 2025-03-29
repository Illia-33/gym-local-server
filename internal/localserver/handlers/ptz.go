package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Illia-33/gym-localserver/camera"
	"github.com/Illia-33/gym-localserver/server/api"

	"github.com/gin-gonic/gin"
)

func checkPlusMinusOne(x float64) bool {
	return -1.0 <= x && x <= 1.0
}

func getCameraId(r *gin.Context, ctx *Context) (int, error) {
	id, err := strconv.Atoi(r.Param("id"))
	if err != nil {
		return 0, err
	}

	if id >= len(ctx.Cameras) {
		return 0, errors.New("not found")
	}

	return id, nil
}

func doStartPtz(ctx *Context, cameraId int, r api.StartPtzRequest) api.StartPtzResponse {
	deadline, err := time.ParseDuration(r.Deadline)
	if err != nil {
		return api.StartPtzResponse{}
	}
	ctx.Cameras[cameraId].PtzController.StartPtz(camera.PtzVelocity{
		Pan:  r.Velocity.Pan,
		Tilt: r.Velocity.Tilt,
		Zoom: r.Velocity.Zoom,
	}, deadline)
	return api.StartPtzResponse{}
}

func doStopPtz(ctx *Context, cameraId int, _ api.EndPtzRequest) api.EndPtzResponse {
	ctx.Cameras[cameraId].PtzController.StopPtz()
	return api.EndPtzResponse{}
}

func StartPtz(r *gin.Context, ctx *Context) {
	id, err := getCameraId(r, ctx)
	if err != nil {
		r.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var request api.StartPtzRequest
	if err := r.BindJSON(&request); err != nil {
		r.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !checkPlusMinusOne(request.Velocity.Pan) ||
		!checkPlusMinusOne(request.Velocity.Tilt) ||
		!checkPlusMinusOne(request.Velocity.Zoom) {
		r.JSON(http.StatusBadRequest, "velocities must be in range [-1;1]")
		return
	}

	r.IndentedJSON(http.StatusOK, doStartPtz(ctx, id, request))
}

func EndPtz(r *gin.Context, ctx *Context) {
	id, err := getCameraId(r, ctx)
	if err != nil {
		r.JSON(http.StatusBadRequest, err)
	}

	r.IndentedJSON(http.StatusOK, doStopPtz(ctx, id, api.EndPtzRequest{}))
}

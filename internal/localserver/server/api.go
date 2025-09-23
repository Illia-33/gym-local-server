package server

import (
	"context"
	"errors"
	"net/http"

	api "github.com/Illia-33/gym-localserver/api/localserver"
	"github.com/Illia-33/gym-localserver/internal/localserver/service"
	"github.com/gin-gonic/gin"
)

type empty struct{}

func registerAPI(r *gin.Engine, service *service.GymCameraService) {
	restAPI := r.Group("/api/v1")

	restAPI.GET("/cameras", handle(
		func(ctx context.Context, rp requestParams[empty]) (api.GetCamerasResponse, error) {
			return service.GetCamerasInfo(ctx)
		},
	))

	restAPI.POST("/camera/:camera_id/ptz",
		withCameraId(),
		withBody[api.StartPtzRequest](),
		handle(
			func(ctx context.Context, rp requestParams[api.StartPtzRequest]) (empty, error) {
				return empty{}, service.StartPtz(ctx, rp.cameraId, &rp.body)
			},
		),
	)
	restAPI.DELETE("/camera/:camera_id/ptz",
		withCameraId(),
		handle(
			func(ctx context.Context, rp requestParams[empty]) (empty, error) {
				return empty{}, service.StopPtz(ctx, rp.cameraId)
			},
		),
	)

	restAPI.POST("/camera/:camera_id/webrtc",
		withCameraId(),
		withBody[api.SetupWebRTCRequest](),
		handle(
			func(ctx context.Context, rp requestParams[api.SetupWebRTCRequest]) (api.SetupWebRTCResponse, error) {
				return service.SetupWebRTC(ctx, rp.cameraId, &rp.body)
			},
		),
	)
}

type requestParams[BodyType any] struct {
	cameraId int
	body     BodyType
}

type handler[BodyType any, ResponseType any] func(context.Context, requestParams[BodyType]) (ResponseType, error)

func handle[BodyType any, ResponseType any](handler handler[BodyType, ResponseType]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body BodyType
		if _, isEmpty := any(body).(empty); !isEmpty {
			value, exists := ctx.Get(key_request_body)
			if !exists {
				ctx.AbortWithError(http.StatusInternalServerError, errors.New("request body is not set"))
				return
			}

			body = value.(BodyType)
		}

		cameraId := func() int {
			if value, exists := ctx.Get(key_camera_id); exists {
				return value.(int)
			}

			return -1
		}()

		response, err := handler(ctx, requestParams[BodyType]{
			cameraId: cameraId,
			body:     body,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}

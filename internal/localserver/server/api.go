package server

import (
	"context"
	"net/http"
	"strconv"

	api "github.com/Illia-33/gym-localserver/api/localserver"
	"github.com/Illia-33/gym-localserver/internal/localserver/service"
	"github.com/gin-gonic/gin"
)

type empty struct{}

func registerAPI(r *gin.Engine, service *service.Service) {
	restAPI := r.Group("/api/v1")

	restAPI.GET("/cameras", handle(func(ctx context.Context, rp requestParams[empty]) (api.GetCamerasResponse, error) {
		return service.GetCamerasInfo(ctx)
	})...)

	restAPI.POST("/cameras/:id/ptz", handle(
		func(ctx context.Context, rp requestParams[api.StartPtzRequest]) (empty, error) {
			return empty{}, service.StartPtz(ctx, rp.cameraId, &rp.body)
		},
	)...)
	restAPI.DELETE("/cameras/:id/ptz", handle(
		func(ctx context.Context, rp requestParams[empty]) (empty, error) {
			return empty{}, service.StopPtz(ctx, rp.cameraId)
		},
	)...)
}

type requestParams[BodyType any] struct {
	cameraId int
	body     BodyType
}

type handler[BodyType any, ResponseType any] func(context.Context, requestParams[BodyType]) (ResponseType, error)

func handle[BodyType any, ResponseType any](handler handler[BodyType, ResponseType]) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			var body BodyType
			if _, isEmpty := any(body).(empty); !isEmpty {
				err := ctx.BindJSON(&body)
				if err != nil {
					ctx.AbortWithError(http.StatusInternalServerError, err)
					return
				}
			}

			var cameraId int
			{
				s := ctx.Param("camera_id")
				if s != "" {
					id, err := strconv.Atoi(s)
					if err != nil {
						ctx.AbortWithError(http.StatusBadRequest, err)
						return
					}

					cameraId = id
				}
			}

			response, err := handler(ctx, requestParams[BodyType]{
				cameraId: cameraId,
				body:     body,
			})
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusOK, response)
		},
	}
}

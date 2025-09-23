package server

import (
	"context"

	"github.com/Illia-33/gym-localserver/internal/localserver/service"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	service *service.GymCameraService
	engine  *gin.Engine
}

func Create(config *cfg.Config) (Server, error) {
	service := service.GymCameraService{}
	err := service.InitWithConfig(config)
	if err != nil {
		return Server{}, err
	}

	r := gin.Default()
	registerAPI(r, &service)

	return Server{
		service: &service,
		engine:  r,
	}, nil
}

func (s *Server) Run(addr string) error {
	s.service.Start(context.Background())
	return s.engine.Run(addr)
}

package server

import (
	"github.com/Illia-33/gym-localserver/internal/localserver/service"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	service *service.Service
	engine  *gin.Engine
}

func Create(config *cfg.Config) (Server, error) {
	service := service.Service{}
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
	return s.engine.Run(addr)
}

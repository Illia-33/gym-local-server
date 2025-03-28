package main

import (
	cfg "gymlocalserver/config"
	"gymlocalserver/server/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ctx    handlers.Context
	engine *gin.Engine
	addr   string
}

func Create(bind string, config *cfg.Config) Server {
	r := gin.Default()
	ctx := handlers.CreateContext(config.Cameras)
	registerApi(r, &ctx)

	return Server{
		engine: r,
		addr:   bind,
		ctx:    ctx,
	}
}

func (s *Server) Run() {
	s.engine.Run(s.addr)
}

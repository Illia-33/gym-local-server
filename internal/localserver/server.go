package localserver

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Illia-33/gym-localserver/internal/localserver/handlers"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ctx    handlers.Context
	engine *gin.Engine
	addr   string
}

func Create(bind string, config *cfg.Config) Server {
	r := gin.Default()
	ctx := handlers.CreateContext(config)
	registerApi(r, &ctx)

	return Server{
		engine: r,
		addr:   bind,
		ctx:    ctx,
	}
}

func (s *Server) Run() {
	go func() {
		authKey := s.ctx.Settings.AuthKey
		for {
			time.Sleep(30 * time.Second)
			jsonBody := []byte(fmt.Sprintf(`{"auth_key":"%s"}`, authKey))
			reader := bytes.NewReader(jsonBody)
			req, err := http.NewRequest(http.MethodPost, "http://89.169.174.232:8080/api/gym/local/assign", reader)
			if err != nil {
				log.Printf("error while gym local assign: %+v", err)
			} else {
				log.Printf("success gym local assign: %+v", req)
			}
		}
	}()
	s.engine.Run(s.addr)
}

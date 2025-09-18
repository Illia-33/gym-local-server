package main

import (
	"flag"
	"log"
	"os"

	"github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
	"github.com/Illia-33/gym-localserver/pkg/onvif"

	"github.com/Illia-33/gym-localserver/internal/localserver/server"

	"gopkg.in/yaml.v3"
)

var (
	bind       = flag.String("bind", "0.0.0.0:8080", "address to bind server on")
	configFile = flag.String("config", "./config.yml", "path to config file")
)

func main() {
	onvifFactory := &onvif.OnvifCameraFactory{}
	camera.RegisterFactory(string(cfg.TypeOnvif), onvifFactory)

	rawConfig, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	var config cfg.Config
	err = yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatalf("yaml unmarshal failed: %v", err)
	}

	server, err := server.Create(&config)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	server.Run(*bind)
}

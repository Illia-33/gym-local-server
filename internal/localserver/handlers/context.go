package handlers

import (
	"log"

	cam "github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
)

type Camera struct {
	Config        cfg.Camera
	PtzController cam.PtzCamera
}

type Context struct {
	Settings *cfg.Settings
	Cameras  []Camera
}

func CreateContext(config *cfg.Config) Context {
	cameras := make([]Camera, 0, len(config.Cameras))
	for _, c := range config.Cameras {
		camera, err := cam.Create(&c)
		if err != nil {
			log.Printf("failed to create %s camera: %v\n", c.Label, err)
			continue
		}
		cameras = append(cameras, Camera{
			Config:        c,
			PtzController: camera,
		})
	}

	log.Printf("Found %d cameras", len(cameras))

	return Context{
		Settings: &config.Settings,
		Cameras:  cameras,
	}
}

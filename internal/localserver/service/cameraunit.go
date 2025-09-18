package service

import (
	cam "github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
)

type cameraUnit struct {
	config cfg.Camera
	camera cam.Camera
}

func newCameraUnit(cfg cfg.Camera) (cameraUnit, error) {
	camera, err := cam.Create(string(cfg.Type), cam.Config{
		Ip:       cfg.Ip,
		Port:     int(cfg.Port),
		Login:    cfg.Login,
		Password: cfg.Password,
	})
	if err != nil {
		return cameraUnit{}, err
	}

	return cameraUnit{
		config: cfg,
		camera: camera,
	}, nil
}

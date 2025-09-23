package service

import (
	cam "github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
)

type CameraUnit struct {
	Config cfg.Camera
	Camera cam.Camera
}

func newCameraUnit(cfg cfg.Camera) (CameraUnit, error) {
	camera, err := cam.Create(string(cfg.Type), cam.Config{
		Ip:       cfg.Ip,
		Port:     int(cfg.Port),
		Login:    cfg.Login,
		Password: cfg.Password,
	})
	if err != nil {
		return CameraUnit{}, err
	}

	return CameraUnit{
		Config: cfg,
		Camera: camera,
	}, nil
}

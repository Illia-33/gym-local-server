package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	api "github.com/Illia-33/gym-localserver/api/localserver"
	"github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
)

type Service struct {
	Cameras  []cameraUnit
	Settings cfg.Settings
}

func (s *Service) InitWithConfig(cfg *cfg.Config) error {
	cameras := make([]cameraUnit, 0, len(cfg.Cameras))
	for _, c := range cfg.Cameras {
		unit, err := newCameraUnit(c)
		if err != nil {
			log.Printf("failed to create %s camera: %v", c.Label, err)
			continue
		}

		cameras = append(cameras, unit)
	}

	if len(cameras) == 0 {
		return ErrorNoWorkingCameras{}
	}

	log.Printf("Found %d cameras", len(cameras))
	s.Settings = cfg.Settings
	s.Cameras = cameras

	return nil
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		authKey := s.Settings.AuthKey
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
}

func (s *Service) GetCamerasInfo(ctx context.Context) (api.GetCamerasResponse, error) {
	cameras := make([]api.CameraDescription, 0, len(s.Cameras))
	for _, camera := range s.Cameras {
		cameras = append(cameras, api.CameraDescription{
			Label:       camera.config.Label,
			Description: camera.config.Description,
		})
	}

	return api.GetCamerasResponse{
		Cameras: cameras,
	}, nil
}

func (s *Service) StartPtz(ctx context.Context, cameraId int, r *api.StartPtzRequest) error {
	if !(0 <= cameraId && cameraId < len(s.Cameras)) {
		return ErrorCameraNotFound{idx: cameraId}
	}

	dur, err := time.ParseDuration(r.Deadline)
	if err != nil {
		return err
	}

	vel := camera.PtzVelocity{
		Pan:  r.Velocity.Pan,
		Tilt: r.Velocity.Tilt,
		Zoom: r.Velocity.Zoom,
	}

	return s.Cameras[cameraId].camera.Ptz.Start(ctx, vel, dur)
}

func (s *Service) StopPtz(ctx context.Context, cameraId int) error {
	if !(0 <= cameraId && cameraId < len(s.Cameras)) {
		return ErrorCameraNotFound{idx: cameraId}
	}

	return s.Cameras[cameraId].camera.Ptz.Stop(ctx)
}

package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	api "github.com/Illia-33/gym-localserver/api/localserver"
	"github.com/Illia-33/gym-localserver/internal/webrtc"
	"github.com/Illia-33/gym-localserver/pkg/camera"
	cfg "github.com/Illia-33/gym-localserver/pkg/config"
)

type GymCameraService struct {
	Units    []CameraUnit
	Settings cfg.Settings
}

func (s *GymCameraService) InitWithConfig(cfg *cfg.Config) error {
	cameras := make([]CameraUnit, 0, len(cfg.Cameras))
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
	s.Units = cameras

	return nil
}

func (s *GymCameraService) Start(ctx context.Context) {
	go func() {
		authKey := s.Settings.AuthKey
		for {
			jsonBody := []byte(fmt.Sprintf(`{"auth_key":"%s"}`, authKey))
			reader := bytes.NewReader(jsonBody)
			req, err := http.NewRequest(http.MethodPost, "http://89.169.174.232:8080/api/gym/local/assign", reader)
			if err != nil {
				log.Printf("error while creating gym assign request: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			client := http.Client{
				Timeout: 10 * time.Second,
			}
			response, err := client.Do(req)
			if err != nil {
				log.Printf("error while doing gym assign request: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if response.StatusCode != http.StatusOK {
				log.Printf("gym assign request status code is %d", response.StatusCode)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Println("successful gym assign")
			time.Sleep(30 * time.Second)
		}
	}()
}

func (s *GymCameraService) GetCamerasInfo(ctx context.Context) (api.GetCamerasResponse, error) {
	cameras := make([]api.CameraDescription, 0, len(s.Units))
	for _, camera := range s.Units {
		cameras = append(cameras, api.CameraDescription{
			Label:       camera.Config.Label,
			Description: camera.Config.Description,
		})
	}

	return api.GetCamerasResponse{
		Cameras: cameras,
	}, nil
}

func (s *GymCameraService) StartPtz(ctx context.Context, cameraId int, r *api.StartPtzRequest) error {
	if err := s.checkCameraIdExists(cameraId); err != nil {
		return err
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

	return s.Units[cameraId].Camera.Ptz.Start(ctx, vel, dur)
}

func (s *GymCameraService) StopPtz(ctx context.Context, cameraId int) error {
	if err := s.checkCameraIdExists(cameraId); err != nil {
		return err
	}

	return s.Units[cameraId].Camera.Ptz.Stop(ctx)
}

func (s *GymCameraService) SetupWebRTC(ctx context.Context, cameraId int, r *api.SetupWebRTCRequest) (api.SetupWebRTCResponse, error) {
	if err := s.checkCameraIdExists(cameraId); err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	stream, err := s.Units[cameraId].Camera.Stream.Setup(ctx)
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	streamDesc, err := stream.Describe(ctx)
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	sdpCodec := streamDesc.FindVideoCodec()
	if sdpCodec == "" {
		return api.SetupWebRTCResponse{}, ErrorNoVideoTrack{}
	}

	peer, err := webrtc.NewRtpPeer(webrtc.RtpPeerConfig{
		Codec: SdpCodecToWebRTCCodec(sdpCodec),
	})
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	offer, err := decodeFromBase64Json[webrtc.SessionDescription](r.OfferBase64)
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	desc, err := peer.Start(webrtc.StartConfig{
		Offer: offer,
	})
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	go func() {
		ctx := context.Background()
		err := stream.Play(ctx, func(ctx context.Context, p camera.Packet) {
			peer.PutPacket(webrtc.RtpPacket(p))
		})
		if err != nil {
			log.Printf("listen ended with error: %v", err)
		}

		stream.Close(ctx)
		peer.Close()
	}()

	sdpBase64, err := encodeToBase64Json(desc)
	if err != nil {
		return api.SetupWebRTCResponse{}, err
	}

	return api.SetupWebRTCResponse{
		Id:              peer.Id,
		LocalDescBase64: sdpBase64,
	}, nil
}

func (s *GymCameraService) checkCameraIdExists(cameraId int) error {
	if !(0 <= cameraId && cameraId < len(s.Units)) {
		return ErrorCameraNotFound{idx: cameraId}
	}
	return nil
}

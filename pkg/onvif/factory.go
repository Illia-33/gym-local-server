package onvif

import cam "github.com/Illia-33/gym-localserver/pkg/camera"

type OnvifCameraFactory struct {
}

func (f *OnvifCameraFactory) Create(c cam.Config) (cam.Camera, error) {
	return CreateOnvifCamera(c)
}

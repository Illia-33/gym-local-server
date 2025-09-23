package onvif

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cam "github.com/Illia-33/gym-localserver/pkg/camera"

	onvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	mediasdk "github.com/use-go/onvif/sdk/media"
)

func CreateOnvifCamera(c cam.Config) (cam.Camera, error) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    fmt.Sprintf("%s:%d", c.Ip, c.Port),
		Username: c.Login,
		Password: c.Password,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		return cam.Camera{}, err
	}

	profiles, err := mediasdk.Call_GetProfiles(context.TODO(), dev, media.GetProfiles{})
	if err != nil {
		return cam.Camera{}, err
	}

	profileToken := profiles.Profiles[0].Token

	return cam.Camera{
		Ptz: &ptzController{
			device:       dev,
			profileToken: profileToken,
		},
		Stream: &streamController{
			device:       dev,
			profileToken: profileToken,
		},
	}, nil
}

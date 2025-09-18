package onvif

import (
	"context"
	"fmt"
	"net/http"

	cam "github.com/Illia-33/gym-localserver/pkg/camera"

	onvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	mediasdk "github.com/use-go/onvif/sdk/media"
)

func CreateOnvifCamera(c cam.Config) (cam.Camera, error) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:      fmt.Sprintf("%s:%d", c.Ip, c.Port),
		Username:   c.Login,
		Password:   c.Password,
		HttpClient: new(http.Client),
	})
	if err != nil {
		return cam.Camera{}, err
	}

	profiles, err := mediasdk.Call_GetProfiles(context.TODO(), dev, media.GetProfiles{})
	if err != nil {
		return cam.Camera{}, err
	}

	return cam.Camera{
		Ptz: &ptzController{
			device:       dev,
			profileToken: profiles.Profiles[0].Token,
		},
	}, nil
}

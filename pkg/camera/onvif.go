package camera

import (
	"context"
	"fmt"
	cfg "gymlocalserver/config"
	"log"
	"net/http"
	"time"

	onvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	sdkMedia "github.com/use-go/onvif/sdk/media"
	sdkPtz "github.com/use-go/onvif/sdk/ptz"
	"github.com/use-go/onvif/xsd"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

type OnvifCamera struct {
	Device       *onvif.Device
	ProfileToken xonvif.ReferenceToken
}

type OnvifCameraFactory struct {
}

func CreateOnvifCamera(c *cfg.Camera) (*OnvifCamera, error) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:      fmt.Sprintf("%s:%d", c.Ip, c.Port),
		Username:   c.Login,
		Password:   c.Password,
		HttpClient: new(http.Client),
	})
	if err != nil {
		return nil, err
	}

	profiles, err := sdkMedia.Call_GetProfiles(context.TODO(), dev, media.GetProfiles{})
	if err != nil {
		return nil, err
	}

	return &OnvifCamera{
		Device:       dev,
		ProfileToken: profiles.Profiles[0].Token, // TODO
	}, nil
}

func (f *OnvifCameraFactory) Create(c *cfg.Camera) (PtzCamera, error) {
	return CreateOnvifCamera(c)
}

func (c *OnvifCamera) StartPtz(vel PtzVelocity, deadline time.Duration) {
	_, err := sdkPtz.Call_ContinuousMove(context.TODO(), c.Device, ptz.ContinuousMove{
		ProfileToken: c.ProfileToken,
		Velocity: xonvif.PTZSpeed{
			PanTilt: xonvif.Vector2D{
				X: vel.Pan,
				Y: vel.Tilt,
			},
			Zoom: xonvif.Vector1D{
				X: vel.Zoom,
			},
		},
		Timeout: xsd.Duration.NewDateTime(
			"0", "0", "0", "0", "0", "0",
			fmt.Sprint(deadline.Milliseconds()/1000),
		),
	})

	if err != nil {
		log.Printf("start ptz failed: %v", err) // TODO
	}
}

func (c *OnvifCamera) StopPtz() {
	_, err := sdkPtz.Call_Stop(context.TODO(), c.Device, ptz.Stop{
		ProfileToken: c.ProfileToken,
		PanTilt:      true,
		Zoom:         true,
	})

	if err != nil {
		log.Printf("stop ptz failed: %v", err) // TODO
	}
}

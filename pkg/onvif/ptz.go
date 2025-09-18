package onvif

import (
	"context"
	"fmt"
	"time"

	cam "github.com/Illia-33/gym-localserver/pkg/camera"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/ptz"
	ptzsdk "github.com/use-go/onvif/sdk/ptz"
	"github.com/use-go/onvif/xsd"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

type ptzController struct {
	device       *onvif.Device
	profileToken xonvif.ReferenceToken
}

func (c *ptzController) Start(ctx context.Context, vel cam.PtzVelocity, deadline time.Duration) error {
	_, err := ptzsdk.Call_ContinuousMove(ctx, c.device, ptz.ContinuousMove{
		ProfileToken: c.profileToken,
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

	return err
}

func (c *ptzController) Stop(ctx context.Context) error {
	_, err := ptzsdk.Call_Stop(ctx, c.device, ptz.Stop{
		ProfileToken: c.profileToken,
		PanTilt:      true,
		Zoom:         true,
	})

	return err
}

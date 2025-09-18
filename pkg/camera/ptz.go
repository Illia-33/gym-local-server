package camera

import (
	"context"
	"time"
)

type PtzVelocity struct {
	Pan  float64
	Tilt float64
	Zoom float64
}

type PtzContoller interface {
	Start(ctx context.Context, vel PtzVelocity, deadline time.Duration) error
	Stop(ctx context.Context) error
}

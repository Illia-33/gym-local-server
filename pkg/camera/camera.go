package camera

import (
	"time"
)

type PtzVelocity struct {
	Pan  float64
	Tilt float64
	Zoom float64
}

type PtzCamera interface {
	StartPtz(vel PtzVelocity, deadline time.Duration)
	StopPtz()
}

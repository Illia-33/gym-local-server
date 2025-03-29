package camera

import (
	cfg "github.com/Illia-33/gym-localserver/config"

	"fmt"
)

type PtzCameraFactory interface {
	Create(c *cfg.Camera) (PtzCamera, error)
}

var cameraFactoriesRegistrar = map[cfg.Type]PtzCameraFactory{}

func RegisterFactory(cameraType cfg.Type, f PtzCameraFactory) {
	cameraFactoriesRegistrar[cameraType] = f
}

func Create(c *cfg.Camera) (PtzCamera, error) {
	factory := cameraFactoriesRegistrar[c.Type]
	if factory == nil {
		return nil, fmt.Errorf("unsupported camera type: %v", c.Type)
	}

	return factory.Create(c)
}

package camera

import (
	"fmt"
)

var camera_factories_registrar = map[string]CameraFactory{}

func RegisterFactory(cameraType string, f CameraFactory) {
	camera_factories_registrar[cameraType] = f
}

func Create(cameraType string, c Config) (Camera, error) {
	factory := camera_factories_registrar[cameraType]
	if factory == nil {
		return Camera{}, fmt.Errorf("unsupported camera type: %v", cameraType)
	}

	return factory.Create(c)
}

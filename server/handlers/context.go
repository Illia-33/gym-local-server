package handlers

import (
	cam "gymlocalserver/camera"
	cfg "gymlocalserver/config"
	"log"
)

type Context struct {
	Cameras []cam.PtzCamera
}

func CreateContext(configs []cfg.Camera) Context {
	cameras := make([]cam.PtzCamera, 0, len(configs))
	for _, c := range configs {
		camera, err := cam.Create(&c)
		if err != nil {
			log.Printf("failed to create %s camera: %v\n", c.Label, err)
			continue
		}
		cameras = append(cameras, camera)
	}

	log.Printf("Found %d cameras", len(cameras))

	return Context{
		Cameras: cameras,
	}
}

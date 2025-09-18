package service

import "fmt"

type ErrorCameraNotFound struct {
	idx int
}

func (e ErrorCameraNotFound) Error() string {
	return fmt.Sprintf("camera with id %d not found", e.idx)
}

type ErrorNoWorkingCameras struct {
}

func (e ErrorNoWorkingCameras) Error() string {
	return "no working camera has been found"
}

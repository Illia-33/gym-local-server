package api

type CamerasRequest struct {
}

type CameraDescription struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type CamerasResponse struct {
	Cameras []CameraDescription `json:"cameras"`
}

package api

type CameraDescription struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type GetCamerasResponse struct {
	Cameras []CameraDescription `json:"cameras"`
}

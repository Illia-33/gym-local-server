package api

type Velocity struct {
	Pan  float64 `json:"pan"`
	Tilt float64 `json:"tilt"`
	Zoom float64 `json:"zoom"`
}

type StartPtzRequest struct {
	Velocity Velocity `json:"velocity"`
	Deadline string   `json:"deadline"`
}

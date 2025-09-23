package api

type SetupWebRTCRequest struct {
	OfferBase64 string `json:"offer_b64"`
}

type SetupWebRTCResponse struct {
	Id              string `json:"id"`
	LocalDescBase64 string `json:"local_desc_b64"`
}

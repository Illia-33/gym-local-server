package service

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Illia-33/gym-localserver/internal/webrtc"
	"github.com/Illia-33/gym-localserver/pkg/sdp"
)

var codec_sdp_to_webrtc = map[sdp.Codec]webrtc.Codec{
	"VP8":  webrtc.CODEC_VP8,
	"VP9":  webrtc.CODEC_VP9,
	"H264": webrtc.CODEC_H264,
	"H265": webrtc.CODEC_H265,
}

func SdpCodecToWebRTCCodec(c sdp.Codec) webrtc.Codec {
	res, exists := codec_sdp_to_webrtc[c]
	if !exists {
		return webrtc.CODEC_UNKNOWN
	}

	return res
}

func encodeToBase64Json(o any) (string, error) {
	asJson, err := json.Marshal(o)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(asJson), nil
}

func decodeFromBase64Json[DataType any](b64 string) (DataType, error) {
	var res DataType
	asJson, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(asJson, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

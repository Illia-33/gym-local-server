package webrtc

import "github.com/pion/webrtc/v4"

type Codec int

const (
	CODEC_UNKNOWN Codec = iota
	CODEC_VP8
	CODEC_VP9
	CODEC_H264
	CODEC_H265
)

func codecToMimeType(codec Codec) string {
	switch codec {
	case CODEC_VP8:
		return webrtc.MimeTypeVP8

	case CODEC_VP9:
		return webrtc.MimeTypeVP9

	case CODEC_H264:
		return webrtc.MimeTypeH264

	case CODEC_H265:
		return webrtc.MimeTypeH265
	}

	return ""
}

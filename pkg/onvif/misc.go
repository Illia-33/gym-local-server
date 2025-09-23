package onvif

import (
	"log"
	"strconv"
	"strings"

	"github.com/Illia-33/gym-localserver/pkg/sdp"
)

// func CodecFromSdp(desc *sdp.SessionDescription) (string, error) {
// 	videoTrack := findVideoTrack(desc)
// 	if videoTrack == -1 {
// 		return "", errors.New("no video track")
// 	}

// 	video := desc.MediaDescriptions[videoTrack]
// 	video.Attributes

// }

func portFromLocalAddr(localAddr string) int {
	split := strings.Split(localAddr, ":")
	if len(split) == 0 {
		log.Panicf("bad local addr: %v", localAddr)
	}

	port, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		log.Panicf("bad port in local addr: %v", localAddr)
	}

	return port
}

func buildVideoStreamUri(baseUri string, desc *sdp.SessionDescription, videoTrackIdx int) string {
	for _, attr := range desc.MediaDescriptions[videoTrackIdx].Attributes {
		if attr.Key == "control" {
			return baseUri + "/" + attr.Value
		}
	}

	return ""
}

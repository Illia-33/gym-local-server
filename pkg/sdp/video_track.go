package sdp

import "strings"

type Codec string

// Find media description index that relates to video track
// Returns -1 if video track has not been found
func (s *SessionDescription) FindVideoTrack() int {
	for i, md := range s.MediaDescriptions {
		if md.MediaName.Media != "video" {
			continue
		}

		if len(md.MediaName.Protos) < 2 {
			continue
		}

		if md.MediaName.Protos[0] != "RTP" || md.MediaName.Protos[1] != "AVP" {
			continue
		}

		return i
	}

	return -1
}

// Returns name of codec used by video track
// Returns empty string if there is no video track
func (s *SessionDescription) FindVideoCodec() Codec {
	idx := s.FindVideoTrack()
	if idx == -1 {
		return ""
	}

	for _, attr := range s.MediaDescriptions[idx].Attributes {
		if attr.Key != "rtpmap" {
			continue
		}

		// value has format "<payload type> <encoding name>/<clock rate> [/<encoding parameters>]"
		val := attr.Value
		spaceSplit := strings.Split(val, " ")
		if len(spaceSplit) < 2 {
			return ""
		}

		slashSplit := strings.Split(spaceSplit[1], "/")
		if len(slashSplit) < 2 {
			return ""
		}

		return Codec(slashSplit[0])
	}

	return ""
}

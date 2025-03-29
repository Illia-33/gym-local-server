package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Illia-33/gym-localserver/camera"
	cfg "github.com/Illia-33/gym-localserver/config"
	"github.com/Illia-33/gym-localserver/rtsp"
	"github.com/Illia-33/gym-localserver/rtsp/sdp"

	"github.com/use-go/onvif/media"
	sdk "github.com/use-go/onvif/sdk/media"
	"github.com/use-go/onvif/xsd/onvif"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findVideoTrack(desc *sdp.SessionDescription) int {
	for i, md := range desc.MediaDescriptions {
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

func buildVideoStreamUri(baseUri string, desc *sdp.SessionDescription, videoTrackIdx int) string {
	for _, attr := range desc.MediaDescriptions[videoTrackIdx].Attributes {
		if attr.Key == "control" {
			return baseUri + "/" + attr.Value
		}
	}

	return ""
}

func main() {
	port := 24567

	camera, err := camera.CreateOnvifCamera(&cfg.Camera{
		Type:     "onvif",
		Ip:       "192.168.0.160",
		Port:     8899,
		Login:    "admin",
		Password: "20240909-3",
	})

	logFatal(err)

	streamUriResponse, err := sdk.Call_GetStreamUri(context.TODO(), camera.Device, media.GetStreamUri{
		StreamSetup: onvif.StreamSetup{
			Stream: "RTP-Unicast",
			Transport: onvif.Transport{
				Protocol: "UDP",
			},
		},
		ProfileToken: camera.ProfileToken,
	})

	logFatal(err)

	uri := string(streamUriResponse.MediaUri.Uri)
	log.Printf("Uri = %s", uri)

	rtspClient, err := rtsp.CreateClientFromUri(uri)
	logFatal(err)

	optionsResponse, err := rtspClient.Options()
	logFatal(err)
	log.Printf("setup: %+v", optionsResponse)

	describeResponse, err := rtspClient.Describe()
	logFatal(err)

	desc := &describeResponse.Description

	videoTrack := findVideoTrack(desc)
	if videoTrack == -1 {
		log.Fatal("video track not found")
	}

	videoUri := buildVideoStreamUri(uri, desc, videoTrack)

	setup, err := rtspClient.Setup(videoUri, port)
	logFatal(err)

	log.Printf("setup: %+v", setup)

	receiver := rtspClient.GetPacketReceiver(port)

	// SETTING UP CONNECTION WITH REMOTE RTP PION

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6004")
	logFatal(err)

	remoteConn, err := net.DialUDP("udp", nil, udpAddr)
	logFatal(err)

	// END SETTING UP

	receiver.OnDataReceived = func(data []byte) {
		remoteConn.Write(data)
	}

	receiver.Listen()

	playResponse, err := rtspClient.Play(setup.Session)
	logFatal(err)
	log.Printf("%+v", playResponse)

	time.Sleep(time.Hour)

	receiver.Close()
	rtspClient.Close()
}

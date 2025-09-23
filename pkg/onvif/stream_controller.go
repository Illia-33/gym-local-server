package onvif

import (
	"context"
	"errors"

	"github.com/Illia-33/gym-localserver/pkg/camera"
	"github.com/Illia-33/gym-localserver/pkg/rtsp"
	"github.com/Illia-33/gym-localserver/pkg/sdp"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	mediasdk "github.com/use-go/onvif/sdk/media"
	xonvif "github.com/use-go/onvif/xsd/onvif"
)

type streamController struct {
	device       *onvif.Device
	profileToken xonvif.ReferenceToken

	cachedStreamUri string
}

func (c *streamController) Setup(ctx context.Context) (camera.Stream, error) {
	streamUri, err := c.getStreamURI(ctx)
	if err != nil {
		return nil, err
	}

	rtspClient, err := rtsp.CreateClientFromUri(streamUri)
	if err != nil {
		return nil, err
	}

	lis, err := newUdpListener()
	if err != nil {
		rtspClient.Close()
		return nil, err
	}

	port := portFromLocalAddr(lis.connection.LocalAddr().String())

	session, desc, err := c.setupRTSP(ctx, streamUri, rtspClient, port)
	if err != nil {
		rtspClient.Close()
		lis.close()
		return nil, err
	}

	return &rtpStream{
		client:    rtspClient,
		lis:       lis,
		sessionId: session,
		description: camera.StreamDescription{
			SessionDescription: desc,
		},
	}, nil
}

func (s *streamController) setupRTSP(_ context.Context, streamUri string, client *rtsp.Client, port int) (session string, desc sdp.SessionDescription, err error) {
	describeResponse, err := client.Describe()
	if err != nil {
		return "", sdp.SessionDescription{}, err
	}
	sdpDesc := &describeResponse.Description

	videoTrackIdx := sdpDesc.FindVideoTrack()
	if videoTrackIdx == -1 {
		return "", sdp.SessionDescription{}, errors.New("video track not found")
	}

	videoUri := buildVideoStreamUri(streamUri, sdpDesc, videoTrackIdx)
	setupResponse, err := client.Setup(videoUri, port)
	if err != nil {
		return "", sdp.SessionDescription{}, err
	}

	return setupResponse.Session, *sdpDesc, nil
}

func (c *streamController) getStreamURI(ctx context.Context) (string, error) {
	if c.cachedStreamUri == "" {
		streamUri, err := c.fetchStreamURI(ctx)
		if err != nil {
			return "", err
		}

		c.cachedStreamUri = streamUri
	}

	return c.cachedStreamUri, nil
}

func (c *streamController) fetchStreamURI(ctx context.Context) (string, error) {
	streamUriResponse, err := mediasdk.Call_GetStreamUri(ctx, c.device, media.GetStreamUri{
		StreamSetup: xonvif.StreamSetup{
			Stream: "RTP-Unicast",
			Transport: xonvif.Transport{
				Protocol: "UDP",
			},
		},
		ProfileToken: xonvif.ReferenceToken(c.profileToken),
	})
	if err != nil {
		return "", err
	}

	uri := string(streamUriResponse.MediaUri.Uri)
	return uri, nil
}

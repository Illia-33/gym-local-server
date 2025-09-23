package onvif

import (
	"context"
	"errors"
	"log"

	"github.com/Illia-33/gym-localserver/pkg/camera"
	"github.com/Illia-33/gym-localserver/pkg/rtsp"
)

type rtpStream struct {
	client      *rtsp.Client
	lis         udpListener
	sessionId   string
	description camera.StreamDescription
}

func (s *rtpStream) Play(ctx context.Context, handler camera.PacketHandler) error {
	if handler == nil {
		log.Panicln("packet handler must be not nil")
	}

	err := s.playRTSP(ctx)
	if err != nil {
		return err
	}

	for {
		var packet [2048]byte
		n, err := s.lis.read(packet[:])
		if err != nil {
			return err
		}

		handler(ctx, packet[:n])
	}
}

func (s *rtpStream) Describe(context.Context) (camera.StreamDescription, error) {
	return s.description, nil
}

func (s *rtpStream) Close(context.Context) error {
	return errors.Join(s.client.Close(), s.lis.close())
}

func (s *rtpStream) playRTSP(_ context.Context) error {
	_, err := s.client.Play(s.sessionId)
	return err
}

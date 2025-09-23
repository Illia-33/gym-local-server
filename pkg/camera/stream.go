package camera

import (
	"context"

	"github.com/Illia-33/gym-localserver/pkg/sdp"
)

type StreamController interface {
	Setup(context.Context) (Stream, error)
}

type Packet []byte
type PacketHandler func(context.Context, Packet)
type StreamDescription struct {
	sdp.SessionDescription
}

type Stream interface {
	Play(context.Context, PacketHandler) error
	Describe(context.Context) (StreamDescription, error)
	Close(context.Context) error
}

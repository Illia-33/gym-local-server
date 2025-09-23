package webrtc

import (
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
)

type RtpPeer struct {
	Id string

	peerConn   *webrtc.PeerConnection
	rtpSener   *webrtc.RTPSender
	videoTrack *webrtc.TrackLocalStaticRTP
	rtpChan    chan RtpPacket
}

type RtpPacket []byte

type RtpPeerConfig struct {
	Codec Codec
}

type SessionDescription struct {
	webrtc.SessionDescription
}

func NewRtpPeer(cfg RtpPeerConfig) (RtpPeer, error) {
	id := uuid.NewString()

	peerConn, err := webrtc.NewPeerConnection(
		webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		},
	)
	if err != nil {
		return RtpPeer{}, err
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: codecToMimeType(cfg.Codec),
	}, "video", id)
	if err != nil {
		return RtpPeer{}, err
	}

	rtpSender, err := peerConn.AddTrack(videoTrack)
	if err != nil {
		return RtpPeer{}, err
	}

	return RtpPeer{
		Id: id,

		peerConn:   peerConn,
		rtpSener:   rtpSender,
		videoTrack: videoTrack,
		rtpChan:    nil,
	}, nil
}

type StartConfig struct {
	Offer SessionDescription
}

func (c *RtpPeer) Start(cfg StartConfig) (*webrtc.SessionDescription, error) {
	err := c.peerConn.SetRemoteDescription(cfg.Offer.SessionDescription)
	if err != nil {
		return nil, err
	}

	answer, err := c.peerConn.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	err = c.peerConn.SetLocalDescription(answer)
	if err != nil {
		return nil, err
	}

	c.peerConn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Printf("(rtp client  %s) rtp peer: state %s", c.Id, state.String())
	})

	c.peerConn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			log.Printf("(rtp client  %s) got ICE candidate: %s %s:%d", c.Id, candidate.Foundation, candidate.Address, candidate.Port)
			c.initPacketProcessing()
		} else {
			log.Printf("(rtp client  %s) ICE candidate gathering is finished", c.Id)
		}
	})

	return c.peerConn.LocalDescription(), nil
}

func (c *RtpPeer) initPacketProcessing() {
	const rtp_chan_capacity = 1000

	rtpChan := make(chan RtpPacket, rtp_chan_capacity)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		const max_packet_send_attempts = 5
		errs := make([]error, 0, max_packet_send_attempts)

		wg.Done()

		for {
			packet, ok := <-rtpChan
			if !ok {
				log.Printf("(rtp client %s) stop packet processing due to close channel", c.Id)
				break
			}

			packetSent := false
			var err error
			for range max_packet_send_attempts {
				_, err = c.videoTrack.Write(packet)
				if err != nil {
					errs = append(errs, err)
					continue
				} else {
					packetSent = true
				}
			}

			if !packetSent {
				log.Printf("(rtp client %s) stop packet processing due to errors while sending: %v", c.Id, errors.Join(errs...))
				break
			}

			// ok, packet sent successfully
			errs = errs[:0]
		}
	}()

	wg.Wait()

	c.rtpChan = rtpChan
}

func (c *RtpPeer) PutPacket(p RtpPacket) {
	if p == nil {
		panic("packet must be non-nil")
	}
	c.rtpChan <- p
}

func (c *RtpPeer) Close() error {
	close(c.rtpChan)
	return c.peerConn.Close()
}

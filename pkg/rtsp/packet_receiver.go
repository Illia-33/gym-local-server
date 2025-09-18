package rtsp

import (
	"fmt"
	"log"
	"net"
)

type PacketReceiver struct {
	Port           int
	OnDataReceived func([]byte)

	listener *net.UDPConn
}

func (pr *PacketReceiver) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", pr.Port))
	if err != nil {
		return err
	}

	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	ln.SetReadBuffer(300000)

	pr.listener = ln
	pr.startProcessing()
	return nil
}

func (pr *PacketReceiver) startProcessing() {
	go func() {
		buffer := make([]byte, 2048)
		for {
			n, _, err := pr.listener.ReadFromUDP(buffer) // TODO check that camera send it
			if err != nil {
				log.Printf("error while receiving data: %v", err)
				continue
			}

			if pr.OnDataReceived != nil {
				pr.OnDataReceived(buffer[:n])
			}
		}
	}()
}

func (pr *PacketReceiver) Close() error {
	return pr.listener.Close()
}

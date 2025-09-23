package onvif

import (
	"net"
)

type udpListener struct {
	connection *net.UDPConn
}

func (l *udpListener) close() error {
	if l.connection == nil {
		return nil
	}

	err := l.connection.Close()
	if err != nil {
		return err
	}

	l.connection = nil
	return nil
}

func (l *udpListener) read(buffer []byte) (int, error) {
	n, _, err := l.connection.ReadFromUDP(buffer)
	return n, err
}

func newUdpListener() (udpListener, error) {
	addr, err := net.ResolveUDPAddr("udp4", ":0")
	if err != nil {
		return udpListener{}, err
	}

	sock, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return udpListener{}, err
	}

	err = sock.SetReadBuffer(262144)
	if err != nil {
		return udpListener{}, err
	}

	return udpListener{
		connection: sock,
	}, nil
}

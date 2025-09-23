package rtsp

import (
	"bufio"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/Illia-33/gym-localserver/pkg/rtsp/requests"
)

const (
	gym_local_server_user_agent string = "GymLocalServer"
)

type Client struct {
	Uri string

	socket *net.TCPConn
	cseq   int
}

func (c *Client) Options() (*requests.OptionsResponse, error) {
	request := requests.BuildOptionsRequest(requests.RequestHeader{
		Uri:       c.Uri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	})

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParseOptionsResponse(response)
}

func (c *Client) Describe() (*requests.DescribeResponse, error) {
	request := requests.BuildDescribeRequest(requests.RequestHeader{
		Uri:       c.Uri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	})
	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParseDescribeResponse(response)
}

func (c *Client) Setup(videoUri string, port int) (*requests.SetupResponse, error) {
	request := requests.BuildSetupRequest(requests.RequestHeader{
		Uri:       videoUri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	}, port)

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParseSetupResponse(response)
}

func (c *Client) Play(sessionId string) (*requests.PlayResponse, error) {
	request := requests.BuildPlayRequest(requests.RequestHeader{
		Uri:       c.Uri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	}, sessionId)

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParsePlayResponse(response)
}

func (c *Client) Pause(sessionId string) (*requests.PauseResponse, error) {
	request := requests.BuildPauseRequest(requests.RequestHeader{
		Uri:       c.Uri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	}, sessionId)

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParsePauseResponse(response)
}

func (c *Client) Teardown(sessionId string) (*requests.TeardownResponse, error) {
	request := requests.BuildTeardownRequest(requests.RequestHeader{
		Uri:       c.Uri,
		CSeq:      c.cseq,
		UserAgent: gym_local_server_user_agent,
	}, sessionId)

	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	return requests.ParseTeardownResponse(response)
}

func (c *Client) Close() error {
	return c.socket.Close()
}

func (c *Client) sendRequest(request string) ([]byte, error) {
	writer := bufio.NewWriter(c.socket)
	_, err := writer.WriteString(request)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(c.socket)
	b := make([]byte, 2048)
	packetSize, err := reader.Read(b)
	if err != nil {
		return nil, err
	}

	c.cseq++
	return b[:packetSize], nil
}

func CreateClientFromUri(uri string) (*Client, error) {
	parsedUri, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(parsedUri.Host, ":") {
		parsedUri.Host = parsedUri.Host + ":554"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", parsedUri.Host)
	if err != nil {
		return nil, err
	}

	sock, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	log.Printf("create RTSP client: %s", tcpAddr.String())

	return &Client{
		Uri:    parsedUri.String(),
		socket: sock,
		cseq:   1,
	}, nil
}

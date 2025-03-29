package requests

import (
	"bytes"
	"errors"
	"gymlocalserver/rtsp/sdp"
	"strings"
)

type DescribeResponse struct {
	ResponseHeader
	Description sdp.SessionDescription
}

func BuildDescribeRequest(header RequestHeader) string {
	var b strings.Builder
	writeHeader(&b, MethodDescribe, header)
	b.WriteString("Accept: application/sdp\r\n")
	writeTail(&b)
	return b.String()
}

func ParseDescribeResponse(response []byte) (*DescribeResponse, error) {
	paragraphs := bytes.Split(response, []byte("\r\n\r\n"))
	if len(paragraphs) < 2 {
		return nil, errors.New("there is not enough paragraphs in response")
	}

	var result DescribeResponse
	err := parseHeader(&result.ResponseHeader, paragraphs[0])
	if err != nil {
		return nil, err
	}

	rawSdp := paragraphs[1]
	err = result.Description.Unmarshal(rawSdp)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

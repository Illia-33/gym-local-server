package requests

import (
	"fmt"
	"strings"
)

type SetupResponse struct {
	ResponseHeader
}

func BuildSetupRequest(header RequestHeader, port int) string {
	var b strings.Builder
	writeHeader(&b, MethodSetup, header)
	b.WriteString(fmt.Sprintf("Transport: RTP/AVP/UDP;unicast;client_port=%d-%d\r\n", port, port+1))
	writeTail(&b)
	return b.String()
}

func ParseSetupResponse(response []byte) (*SetupResponse, error) {
	var result SetupResponse
	err := parseHeader(&result.ResponseHeader, response)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

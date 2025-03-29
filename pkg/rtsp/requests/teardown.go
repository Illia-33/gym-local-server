package requests

import (
	"fmt"
	"strings"
)

type TeardownResponse struct {
	ResponseHeader
}

func BuildTeardownRequest(header RequestHeader, sessionId string) string {
	var b strings.Builder
	writeHeader(&b, MethodPlay, header)
	b.WriteString(fmt.Sprintf("Session: %s\r\n", sessionId))
	writeTail(&b)
	return b.String()
}

func ParseTeardownResponse(response []byte) (*TeardownResponse, error) {
	var result TeardownResponse
	err := parseHeader(&result.ResponseHeader, response)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

package requests

import (
	"fmt"
	"strings"
)

type PauseResponse struct {
	ResponseHeader
}

func BuildPauseRequest(header RequestHeader, sessionId string) string {
	var b strings.Builder
	writeHeader(&b, MethodPlay, header)
	b.WriteString(fmt.Sprintf("Session: %s\r\n", sessionId))
	writeTail(&b)
	return b.String()
}

func ParsePauseResponse(response []byte) (*PauseResponse, error) {
	var result PauseResponse
	err := parseHeader(&result.ResponseHeader, response)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

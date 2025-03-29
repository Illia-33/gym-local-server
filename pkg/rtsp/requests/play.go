package requests

import (
	"fmt"
	"strings"
)

type PlayResponse struct {
	ResponseHeader
}

func BuildPlayRequest(header RequestHeader, sessionId string) string {
	var b strings.Builder
	writeHeader(&b, MethodPlay, header)
	b.WriteString(fmt.Sprintf("Session: %s\r\n", sessionId))
	writeTail(&b)
	return b.String()
}

func ParsePlayResponse(response []byte) (*PlayResponse, error) {
	var result PlayResponse
	err := parseHeader(&result.ResponseHeader, response)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

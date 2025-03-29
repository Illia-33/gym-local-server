package requests

import (
	"strings"
)

type OptionsResponse struct {
	ResponseHeader
	AvailableMethods []Method
}

func BuildOptionsRequest(header RequestHeader) string {
	var b strings.Builder
	writeHeader(&b, MethodOptions, header)
	writeTail(&b)
	return b.String()
}

func ParseOptionsResponse(response []byte) (*OptionsResponse, error) {
	var result OptionsResponse
	err := parseHeader(&result.ResponseHeader, response)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

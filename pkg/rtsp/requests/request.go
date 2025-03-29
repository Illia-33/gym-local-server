package requests

import (
	"fmt"
	"strings"
)

type RequestHeader struct {
	Uri       string
	CSeq      int
	UserAgent string
}

func writeHeader(b *strings.Builder, method Method, header RequestHeader) {
	b.WriteString(fmt.Sprintf("%s %s RTSP/1.0\r\n", method, header.Uri))
	b.WriteString(fmt.Sprintf("CSeq: %d\r\n", header.CSeq))
	b.WriteString(fmt.Sprintf("User-Agent: %s\r\n", header.UserAgent))
}

func writeTail(b *strings.Builder) {
	b.WriteString("\r\n")
}

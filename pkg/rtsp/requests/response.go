package requests

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type ResponseHeader struct {
	StatusCode      int
	CSeq            int    `rtsp:"CSeq"`
	ContentLength   int    `rtsp:"Content-Length"`
	ContentLocation string `rtsp:"Content-Location"`
	ContentType     string `rtsp:"Content-Type"`
	Session         string `rtsp:"Session"`
	LastModified    string `rtsp:"Last-Modified"`
}

func parseStatusLine(rawStatusLine []byte) (int, error) {
	split := bytes.Split(rawStatusLine, []byte(" "))
	if len(split) != 3 {
		return 0, fmt.Errorf("bad status line, must be 3 words: %s", string(rawStatusLine))
	}

	if !bytes.Equal(split[0], []byte("RTSP/1.0")) {
		return 0, fmt.Errorf("only RTSP/1.0 is supported")
	}

	return strconv.Atoi(string(split[1]))
}

func parseHeader(h *ResponseHeader, rawHeader []byte) error {
	var header ResponseHeader

	split := bytes.Split(rawHeader, []byte("\r\n"))
	statusCode, err := parseStatusLine(split[0])
	if err != nil {
		return err
	}
	header.StatusCode = statusCode

	split = split[1:]
	t := reflect.TypeOf(header)
	for _, s := range split {
		if len(s) == 0 {
			continue
		}
		colon := bytes.Index(s, []byte(": "))
		if colon == -1 {
			return fmt.Errorf("bad response header row: %s", string(s))
		}

		key := string(s[:colon])
		value := string(s[colon+2:])

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("rtsp") != key {
				continue
			}

			if key == "Content-Length" || key == "CSeq" {
				num, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("cannot convert to int: %s", string(s))
				}
				reflect.ValueOf(&header).Elem().Field(i).SetInt(int64(num))
			} else {
				reflect.ValueOf(&header).Elem().Field(i).SetString(value)
			}
			break
		}
	}

	*h = header
	return nil
}

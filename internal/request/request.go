package request

import (
	"errors"
	"io"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
}

func (rl *RequestLine) ValidMethod() bool {
	// TODO
	return true
}

func (rl *RequestLine) ValidRequestTarget() bool {
	// TODO
	return true
}

func (rl *RequestLine) ValidHTTPVersion() bool {
	return rl.HTTPVersion == "1.1"
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

var ERROR_BAD_REQUEST_LINE = errors.New("Bad request line")
var ERROR_INCOMPLETE_REQUEST_LINE = errors.New("Incomplete request line")
var ERROR_UNSUPPORTED_HTTP_VERSION = errors.New("Unsupported HTTP version")
var SEPARATOR = "\r\n"

func parseRequestLine(input string) (*RequestLine, string, error) {
	line, rest, isFound := strings.Cut(input, SEPARATOR)
	if !isFound {
		return nil, input, ERROR_INCOMPLETE_REQUEST_LINE
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, rest, ERROR_BAD_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, rest, ERROR_BAD_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HTTPVersion:   httpParts[1],
	}

	if !rl.ValidHTTPVersion() {
		return nil, rest, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return rl, rest, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(errors.New("Unable to read io.ReadAll"), err)
	}

	str := string(data)
	rl, str, err := parseRequestLine(str)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse request line"), err)
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}

package request

import (
	"bytes"
	"errors"
	"io"
)

var SEP = []byte("\r\n")

var (
	ErrNeedMoreData           = errors.New("Need more data")
	ErrInvalidRequestLine     = errors.New("Invalid request line")
	ErrUnsupportedHTTPVersion = errors.New("Unsupported HTTP version - only HTTP/1.1 is supported")
	ErrInvalidHeaders         = errors.New("Invalid headers")
	ErrHeadersTooLarge        = errors.New("Headers too large")
)

type RequestLine struct {
	Method      string
	Target      string
	HTTPVersion string
}

type Request struct {
	Line    RequestLine
	Headers map[string]string
	Body    []byte
}

type ParserState int

const (
	ParserStateRequestLine ParserState = iota
	ParserStateHeaders
	ParserStateBody
	ParserStateComplete
	ParserStateError
)

type RequestParser struct {
	State   ParserState
	Request Request

	buffer bytes.Buffer
}

func NewRequestParser() *RequestParser {
	return &RequestParser{
		State: ParserStateRequestLine,
	}
}

func (p *RequestParser) Parse(data []byte) (int, error) {
	p.buffer.Write(data)

	for {
		switch p.State {
		case ParserStateRequestLine:
			return p.parseRequestLine()
		case ParserStateComplete:
			return 0, nil
		default:
			return 0, nil
		}
	}
}

func (p *RequestParser) parseRequestLine() (int, error) {
	data := p.buffer.Bytes()

	line, _, found := bytes.Cut(data, SEP)
	n := len(line) + len(SEP)

	if !found {
		return n, ErrNeedMoreData
	}

	parts := bytes.SplitN(line, []byte(" "), 3)
	if len(parts) != 3 {
		return n, ErrInvalidRequestLine
	}

	proto, version, found := bytes.Cut(parts[2], []byte("/"))
	if !found {
		return n, ErrInvalidRequestLine
	}

	if !bytes.Equal(proto, []byte("HTTP")) {
		return n, ErrInvalidRequestLine
	}

	if !bytes.Equal(version, []byte("1.1")) {
		return n, ErrUnsupportedHTTPVersion
	}

	p.Request.Line.Method = string(parts[0])
	p.Request.Line.Target = string(parts[1])
	p.Request.Line.HTTPVersion = string(version)

	p.buffer.Next(n)
	p.State = ParserStateHeaders

	return n, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("Unable to io.ReadAll")
	}

	parser := NewRequestParser()
	_, err = parser.Parse(data)
	if err != nil {
		return nil, errors.Join(errors.New("Unable to parse"), err)
	}

	return &parser.Request, nil
}

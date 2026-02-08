package request

import (
	"bytes"
	"errors"
	"io"
)

var SEP = []byte("\r\n")

var (
	ErrInvalidRequestLine     = errors.New("Invalid request line")
	ErrUnsupportedHTTPVersion = errors.New("Unsupported HTTP version - only HTTP/1.1 is supported")
	ErrInvalidHeaders         = errors.New("Invalid headers")
	ErrHeadersTooLarge        = errors.New("Headers too large")
	ErrIncompleteRequest      = errors.New("Incomplete request")
	ErrNeedMoreData           = errors.New("Need more data")
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

func (p *RequestParser) Complete() bool {
	return p.State == ParserStateComplete
}

func (p *RequestParser) Parse(data []byte) error {
	p.buffer.Write(data)

	for {
		switch p.State {
		case ParserStateRequestLine:
			if err := p.parseRequestLine(); err != nil {
				return err
			}
		case ParserStateHeaders:
			// TODO: Implement
			p.State = ParserStateComplete
		case ParserStateBody:
			// TODO: Implement
			p.State = ParserStateComplete
		case ParserStateComplete:
			return nil
		case ParserStateError:
			return errors.New("Parser in error state")
		}
	}
}

func (p *RequestParser) parseRequestLine() error {
	data := p.buffer.Bytes()

	line, _, found := bytes.Cut(data, SEP)
	n := len(line) + len(SEP)

	if !found {
		return ErrNeedMoreData
	}

	parts := bytes.SplitN(line, []byte(" "), 3)
	if len(parts) != 3 {
		return ErrInvalidRequestLine
	}

	proto, version, found := bytes.Cut(parts[2], []byte("/"))
	if !found {
		return ErrInvalidRequestLine
	}

	if !bytes.Equal(proto, []byte("HTTP")) {
		return ErrInvalidRequestLine
	}

	if !bytes.Equal(version, []byte("1.1")) {
		return ErrUnsupportedHTTPVersion
	}

	p.Request.Line.Method = string(parts[0])
	p.Request.Line.Target = string(parts[1])
	p.Request.Line.HTTPVersion = string(version)

	p.buffer.Next(n)
	p.State = ParserStateComplete

	return nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	parser := NewRequestParser()
	buf := make([]byte, 1024)

	for !parser.Complete() {
		n, err := reader.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if err := parser.Parse(buf[:n]); err != nil && err != ErrNeedMoreData {
			return nil, err
		}
	}

	if !parser.Complete() {
		return nil, ErrIncompleteRequest
	}

	return &parser.Request, nil
}

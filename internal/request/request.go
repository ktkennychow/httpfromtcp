package request

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ParserState int

type Request struct {
	RequestLine   RequestLine
	stateOfParser ParserState
}

const (
	initialized ParserState = iota
	done
)
const bufferSize = 8
const crlf = "\r\n"

var capAlphaRegex = regexp.MustCompile(`^[A-Z]+$`)

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	numOfBytes := 0
	newLineIndex := bytes.Index(data, []byte(crlf))
	if newLineIndex == -1 {
		return nil, numOfBytes, nil
	}
	numOfBytes += len(data)
	requestLineString := string(data)[:newLineIndex]
	parts := strings.Split(requestLineString, " ")

	if len(parts) > 3 {
		return nil, numOfBytes, errors.New("request line has more than three parts")
	}
	if len(parts) < 3 {
		return nil, numOfBytes, errors.New("request line has less than three parts")
	}
	if !capAlphaRegex.MatchString(parts[0]) {
		return nil, numOfBytes, errors.New("method should contain only capital alphabetic characters")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, numOfBytes, errors.New("http version not supported")
	}

	requestLine := RequestLine{
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	return &requestLine, numOfBytes, nil
}

func (r *Request) parse(data []byte) (int, error) {
	numOfBytes := 0
	if r.stateOfParser == initialized {
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		numOfBytes += n
		r.RequestLine = *requestLine
	} else if r.stateOfParser == done {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("error: trying to read data in an unknown state")
	}
	r.stateOfParser = done
	return numOfBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := Request{stateOfParser: initialized}

	for request.stateOfParser != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			request.stateOfParser = done
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return &request, nil
}

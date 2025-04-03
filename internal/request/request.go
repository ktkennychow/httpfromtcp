package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ParserState int

type Request struct {
	Headers       headers.Headers
	RequestLine   RequestLine
	stateOfParser ParserState
}

const (
	stateInitialized ParserState = iota
	stateParsingHeaders
	stateDone
)

const bufferSize = 8
const crlf = "\r\n"

func isCapitalAlphaChar(char rune) bool {
	if char < 'A' || char > 'Z' {
		return false
	}
	return true
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	newLineIndex := bytes.Index(data, []byte(crlf))
	if newLineIndex == -1 {
		return nil, 0, nil
	}
	requestLineString := string(data)[:newLineIndex]
	parts := strings.Split(requestLineString, " ")

	if len(parts) > 3 {
		return nil, 0, errors.New("request line has more than three parts")
	}
	if len(parts) < 3 {
		return nil, 0, errors.New("request line has less than three parts")
	}
	hasOnlyCapitalAlphaChar := true
	for _, char := range parts[0] {
		if !isCapitalAlphaChar(char) {
			hasOnlyCapitalAlphaChar = false
			break
		}
	}
	if !hasOnlyCapitalAlphaChar {
		return nil, 0, errors.New("method should contain only capital alphabetic characters")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, 0, errors.New("http version not supported")
	}

	requestLine := RequestLine{
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	return &requestLine, newLineIndex + 2, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.stateOfParser {
	case stateInitialized:
		// fmt.Println("parsing request line")
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.stateOfParser = stateParsingHeaders
		return n, nil
	case stateParsingHeaders:
		// fmt.Println("parsing headers")
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.stateOfParser = stateDone
		}
		return n, nil
	case stateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
	}
	return 0, errors.New("error: trying to read data in an unknown state")
}

func (r *Request) parse(data []byte) (int, error) {
	totalNumOfBytes := 0
	for r.stateOfParser != stateDone {
		// fmt.Println("totalNumOfBytesParsed:", totalNumOfBytes)
		// fmt.Println("string passed into parseSingle", string(data[totalNumOfBytes:]))
		n, err := r.parseSingle(data[totalNumOfBytes:])
		// fmt.Println("bytes parsed in parseSingle", n)
		if err != nil {
			return 0, err
		}
		totalNumOfBytes += n
		if n == 0 {
			break
		}
	}
	return totalNumOfBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := Request{stateOfParser: stateInitialized, Headers: headers.NewHeaders()}

	for request.stateOfParser != stateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			if request.stateOfParser != stateDone {
				// e.g. missing end of headers but see end of file already
				return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.stateOfParser, numBytesRead)
			}
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

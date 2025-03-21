package request

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var capAlphaRegex = regexp.MustCompile(`^[A-Z]+$`)

func RequestFromReader(reader io.Reader) (*Request, error) {
	dat, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	newLineIndex := strings.Index(string(dat), "\r\n")
	requestLine := string(dat)[:newLineIndex]
	parts := strings.Split(requestLine, " ")

	if len(parts) > 3 {
		return nil, errors.New("request line has more than three parts")
	}
	if len(parts) < 3 {
		return nil, errors.New("request line has less than three parts")
	}
	if !capAlphaRegex.MatchString(parts[0]) {
		return nil, errors.New("method should contain only capital alphabetic characters")
	}
	if parts[2] != "HTTP/1.1" {
		return nil, errors.New("http version not supported")
	}

	request := Request{
		RequestLine: RequestLine{
			HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
			RequestTarget: parts[1],
			Method:        parts[0],
		},
	}
	return &request, nil
}

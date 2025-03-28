package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (int, bool, error) {
	numOfBytes := 0
	newLineIndex := bytes.Index(data, []byte(crlf))
	if newLineIndex == -1 {
		return numOfBytes, false, nil
	}
	if newLineIndex == 0 {
		return numOfBytes, true, nil
	}
	numOfBytes += len(data) - 2

	headersString := string(data)[:newLineIndex]
	headersString = strings.TrimSpace(headersString)
	colonIndex := strings.Index(headersString, ":")
	if string(headersString[colonIndex-1]) == " " {
		return 0, false, errors.New("invalid format: a space before colon")
	}
	parts := strings.Split(headersString, ": ")
	h[parts[0]] = parts[1]

	return numOfBytes, false, nil
}

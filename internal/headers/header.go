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

const specialChars = "!#$%&'*+-.^_`|~"

func isCapitalAlphaOrSpecialChar(char rune) bool {
	if char >= 'A' && char <= 'Z' {
		return true
	}
	if char >= 'a' && char <= 'z' {
		return true
	}
	index := strings.Index(specialChars, string(char))
	return index >= 0
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	newLineIndex := bytes.Index(data, []byte(crlf))
	if newLineIndex == -1 {
		return 0, false, nil
	}
	// extra CRLF between headers and request body
	if newLineIndex == 0 {
		return 2, true, nil
	}

	headersString := string(data)[:newLineIndex]
	headersString = strings.TrimSpace(headersString)
	colonIndex := strings.Index(headersString, ":")
	if colonIndex == -1 {
		// fmt.Println("no colon found", headersString)
		return 0, false, errors.New("invalid format: no colon found")
	}
	if colonIndex > 0 && string(headersString[colonIndex-1]) == " " {
		return 0, false, errors.New("invalid format: a space before colon")
	}
	parts := strings.Split(headersString, ": ")
	hasOnlyCapitalAlphaOrSpecialChar := true
	for _, char := range parts[0] {
		if !isCapitalAlphaOrSpecialChar(char) {
			hasOnlyCapitalAlphaOrSpecialChar = false
			break
		}
	}
	if !hasOnlyCapitalAlphaOrSpecialChar {
		return 0, false, errors.New("invalid format: header name contains a invalid character")
	}
	if h[strings.ToLower(parts[0])] != "" {
		h[strings.ToLower(parts[0])] = h[strings.ToLower(parts[0])] + ", " + parts[1]
	} else {
		h[strings.ToLower(parts[0])] = parts[1]
	}
	return newLineIndex + 2, false, nil
}

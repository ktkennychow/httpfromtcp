package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	string := "Host: localhost:42069\r\n\r\n"
	headers := NewHeaders()
	data := []byte(string)
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, len(string)-2, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	string = "   Host: localhost:42069\r\n\r\n"
	headers = NewHeaders()
	data = []byte(string)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, len(string)-2, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	string = "User-Agent: curl/7.81.0\r\n\r\n"
	headers = NewHeaders()
	headers.Parse([]byte("Host: localhost:42069\r\n\r\n"))
	data = []byte(string)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "curl/7.81.0", headers["User-Agent"])
	assert.Equal(t, len(string)-2, n)
	assert.False(t, done)

	// Test: Valid done
	string = ""
	headers = NewHeaders()
	data = []byte(string)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	string = "       Host : localhost:42069       \r\n\r\n"
	headers = NewHeaders()
	data = []byte(string)
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

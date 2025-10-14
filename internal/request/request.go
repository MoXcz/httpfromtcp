package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       State
}

type State = int

const (
	INIT = iota
	DONE
)

const bufferSize = 8

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := Request{state: INIT}
	for req.state != DONE {
		if len(buf) == readToIndex {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		read, err := reader.Read(buf[readToIndex:])
		if err == io.EOF || (read == 0 && err == nil) {
			req.state = DONE
			break
		} else if err != nil {
			return nil, err
		}
		readToIndex += read

		parsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsed == 0 {
			continue
		}
		readToIndex = parsed
	}

	return &req, nil
}

func parseRequestLine(s string) (int, error) {
	// not sure if it's sure for a request line to have a single ' ' character
	// space, but let's assume for now that it does, otherwise it will be malformed
	reqParts := strings.Split(s, "\r\n")
	// needs more data
	if len(reqParts) < 2 {
		return 0, nil
	}
	parts := strings.Split(reqParts[0], " ")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid request line")
	}
	bytesRead := len(parts[0])
	if !isMethod(parts[0]) {
		return bytesRead, fmt.Errorf("invalid request line method")
	}
	bytesRead += len(parts[1])
	if !isRequestTarget(parts[1]) {
		return bytesRead, fmt.Errorf("invalid request line request target")
	}
	bytesRead += len(parts[2])
	if !isHttpVersion(parts[2]) {
		return bytesRead, fmt.Errorf("invalid request line HTTP version")
	}
	return bytesRead, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case INIT:
		bytesRead, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		// needs more data
		if bytesRead == 0 {
			return 0, nil
		}
		reqParts := strings.Split(string(data), "\r\n")
		parts := strings.Split(reqParts[0], " ")
		r.RequestLine = RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   strings.Split(parts[2], "/")[1],
		}
		r.state = DONE
		return bytesRead, nil
	case DONE:
		return 0, fmt.Errorf("trying to read data in a done state")
	}
	return 0, fmt.Errorf("invalid state")
}

// TODO: fill missing methods
func isMethod(m string) bool {
	for b := range m {
		if b < 'A' && b > 'Z' {
			return false
		}
	}
	// something like this is not viable in a stream context
	// return m == "GET" || m == "POST"
	return true
}

// '/', '/cat', '/api/resource'
func isRequestTarget(rt string) bool {
	if rt == "/" {
		return true
	}
	parts := strings.Split(rt, "/")
	if len(parts)%2 == 0 {
		return true
	}
	return false
}

func isHttpVersion(hv string) bool {
	parts := strings.Split(hv, "/")
	if len(parts) != 2 {
		return false
	}
	http := parts[0]
	if http != "HTTP" {
		return false
	}
	version := parts[1]
	if version != "1.1" {
		return false
	}
	return true
}

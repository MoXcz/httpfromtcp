package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       reqState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type reqState = int

const (
	INIT = iota
	DONE
)

const bufferSize = 8
const clrf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := Request{state: INIT}
	for req.state != DONE {
		if len(buf) <= readToIndex {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		readBytes, err := reader.Read(buf[readToIndex:])
		if err == io.EOF || (readBytes == 0 && err == nil) {
			req.state = DONE
			break
		} else if err != nil {
			return nil, err
		}
		readToIndex += readBytes

		parsedBytes, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// 'remove' parsed data
		copy(buf, buf[parsedBytes:])
		readToIndex -= parsedBytes
	}

	return &req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// verify that data can be separated by clrf
	idx := bytes.Index(data, []byte(clrf))
	// needs more data
	if idx == -1 {
		return nil, 0, nil
	}
	// not sure if it's sure for a request line to have a single ' ' character
	// space, but let's assume for now that it does, otherwise it will be malformed
	// note that 'idx' should have the number up to, but not including, the first clrf
	parts := strings.Split(string(data[:idx]), " ")
	if len(parts) < 3 {
		return nil, 0, fmt.Errorf("invalid request line")
	}
	bytesRead := len(parts[0])
	if !isMethod(parts[0]) {
		return nil, bytesRead, fmt.Errorf("invalid request line method")
	}
	bytesRead += len(parts[1])
	if !isRequestTarget(parts[1]) {
		return nil, bytesRead, fmt.Errorf("invalid request line request target")
	}
	bytesRead += len(parts[2])
	if !isHttpVersion(parts[2]) {
		return nil, bytesRead, fmt.Errorf("invalid request line HTTP version")
	}
	version := strings.Split(parts[2], "/")[1]
	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, bytesRead, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case INIT:
		reqLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		// needs more data
		if bytesRead == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.state = DONE
		return bytesRead, nil
	case DONE:
		return 0, fmt.Errorf("trying to read data in a done state")
	}
	return 0, fmt.Errorf("invalid state")
}

// TODO: fill missing methods
func isMethod(m string) bool {
	return m == "GET" || m == "POST" || m == "OPTION"
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

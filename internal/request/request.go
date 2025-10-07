package request

import (
	"fmt"
	"io"
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// premature checking
	parts := strings.Split(string(data), "\r\n")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid format, missing CLRN")
	}

	reqLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: reqLine,
	}, nil
}

func parseRequestLine(s string) (RequestLine, error) {
	// not sure if it's sure for a request line to have a single ' ' character
	// space, but let's assume for now that it does, otherwise it will be malformed
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invallid request line")
	}
	if !isMethod(parts[0]) {
		return RequestLine{}, fmt.Errorf("invalid request line method")
	}
	method := parts[0]
	if !isRequestTarget(parts[1]) {
		return RequestLine{}, fmt.Errorf("invalid request line request target")
	}
	requestTarget := parts[1]
	if !isHttpVersion(parts[2]) {
		return RequestLine{}, fmt.Errorf("invalid request line HTTP version")
	}
	httpVersion := strings.Split(parts[2], "/")[1]
	return RequestLine{
		Method: method, RequestTarget: requestTarget, HttpVersion: httpVersion}, nil
}

// TODO: fill missing methods
func isMethod(m string) bool {
	return m == "GET" || m == "POST"
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

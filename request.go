package main

import (
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	method string
	path   string
}

func (r Request) asString() string {
	return fmt.Sprintf("%s %s HTTP/1.1", r.method, r.path)
}

func newRequest(line string) (Request, error) {
	var req Request
	requestAsArray := strings.SplitN(line, " ", 3)

	if len(requestAsArray) != 3 {
		return req, errors.New("Errors parsing the request")
	}

	method := requestAsArray[0]
	path := requestAsArray[1]

	req = Request{
		method: method,
		path:   path,
	}

	return req, nil
}

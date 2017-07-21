package main

import "fmt"

var statusCodes = map[int]string{
	200: "OK",
	400: "Bad Request",
	404: "Not Found",
	501: "Not Implemented",
}

type Response struct {
	code          int
	contentType   string
	contentLength int
	body          []byte
}

func (r Response) BuildHeaders() string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n"+
		"Content-Type: %s\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n",
		r.code, statusCodes[r.code], r.contentType, len(r.body))
}

func (r Response) BuildBody() string {
	return fmt.Sprintf("\r\n%s", string(r.body))
}

func createResponse(code int, contentType string, data []byte) Response {
	return Response{
		code:        code,
		contentType: contentType,
		body:        data,
	}
}

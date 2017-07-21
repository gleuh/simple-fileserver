package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"net/http"
	"path/filepath"
)

type HttpError struct {
	msg  string
	code int
}

func (e *HttpError) Error() string { return e.msg }

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	httpCommand, _ := bufio.NewReader(conn).ReadString('\n')
	req, err := newRequest(strings.TrimSpace(httpCommand))

	if err != nil {
		sendResponse(conn, req, createResponse(400, "", make([]byte, 0)))
	}

	if req.method != "GET" && req.method != "HEAD" {
		sendResponse(conn, req, createResponse(501, "", make([]byte, 0)))
	}

	fmt.Println("Handle: " + req.asString())

	// handle the file and build the response
	filePath, err := getFilePath(req.path)

	if err != nil {
		fmt.Println(err)
	}

	data, mimeType, err := readFile(filePath)

	response := Response{
		code:        200,
		contentType: mimeType,
		body:        data,
	}

	sendResponse(conn, req, response)
}

func sendResponse(conn net.Conn, req Request, resp Response) {
	send := resp.BuildHeaders()

	if req.method == "GET" {
		send += resp.BuildBody()
	}

	conn.Write([]byte(send))
}

func getFilePath(path string) (string, error) {
	absoluteBasePath, err := filepath.Abs("public")
	basePath := absoluteBasePath + path

	// security check: do not go upper than the public/ directory
	wantedFilePath, err := filepath.Abs(basePath)

	if !strings.HasPrefix(wantedFilePath, absoluteBasePath) {
		basePath = absoluteBasePath
	}

	fi, err := os.Stat(basePath)

	if err != nil {
		return path, err
	}

	if fi.Mode().IsDir() {
		basePath += "/index.html"
	}

	return basePath, nil
}

func readFile(filePath string) ([]byte, string, error) {
	data, err := ioutil.ReadFile(filePath)
	mimeType := http.DetectContentType(data)

	return data, mimeType, err
}

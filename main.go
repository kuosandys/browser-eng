package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
)

const (
	network = "tcp"
)

var (
	ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")
	supportedSchemes        = map[string]struct{}{"http": struct{}{}, "https": struct{}{}}
)

func request(requestURL string) (map[string]string, string, error) {
	var err error
	var headers = map[string]string{}
	var body string

	u, err := url.Parse(requestURL)
	if err != nil {
		return headers, body, err
	}

	if _, ok := supportedSchemes[u.Scheme]; !ok {
		return headers, body, ErrUnsupportedURLScheme
	}

	// port
	port := u.Port()
	if port == "" && u.Scheme == "http" {
		port = "80"
	} else if port == "" && u.Scheme == "https" {
		port = "443"
	}

	// path
	path := u.Path
	if path == "" {
		path = "/index.html"
	}

	// establish the connection
	var connection net.Conn
	if u.Scheme == "http" {
		connection, err = net.Dial(network, fmt.Sprintf("%s:%s", u.Hostname(), port))
	} else if u.Scheme == "https" {
		connection, err = tls.Dial(network, fmt.Sprintf("%s:%s", u.Hostname(), port), &tls.Config{ServerName: u.Host})
	}
	if err != nil {
		return headers, body, err
	}

	defer connection.Close()

	// send
	request := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	requestHeaders := map[string]string{
		"Host":       u.Hostname(),
		"Connection": "close",
	}
	for h, v := range requestHeaders {
		request += fmt.Sprintf("%s: %s\r\n", h, v)
	}
	request += "\r\n"

	_, err = connection.Write([]byte(request))
	if err != nil {
		return headers, body, err
	}

	// receive
	c := bufio.NewReader(connection)
	if err != nil {
		return headers, body, err
	}

	// status
	statusLine, _, err := c.ReadLine()
	if err != nil {
		return headers, body, err
	}
	s := strings.SplitN(string(statusLine), " ", 3)
	version, status, explanation := s[0], s[1], s[2]
	fmt.Println(version, status, explanation)

	// headers
	var line []byte
	for {
		line, _, err = c.ReadLine()
		if err != nil {
			return headers, body, err
		}

		if string(line) == "" {
			break
		}
		s := strings.SplitN(string(line), ":", 2)
		header, value := s[0], s[1]
		headers[strings.TrimSpace(header)] = strings.ToLower(value)
	}

	// body
	var buffer bytes.Buffer
	for {
		line, err := c.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				break
			}
			return headers, body, err
		}
		buffer.Write(line)
	}
	body = buffer.String()

	return headers, body, err
}

func show(body string) {
	var inTag bool

	for _, r := range body {
		c := string(r)
		if c == "<" {
			inTag = true
		} else if c == ">" {
			inTag = false
		} else if !inTag {
			fmt.Print(c)
		}
	}
}

func load(url string) {
	_, body, err := request(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	show(body)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please input a URL")
		os.Exit(1)
	}

	load(args[1])
}

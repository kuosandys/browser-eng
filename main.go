package main

import (
	"bufio"
	"errors"
	"fmt"
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
	supportedScheme         = "http"
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

	// establish the connection
	var connection net.Conn
	connection, err = net.Dial(network, fmt.Sprintf("%s:%s", u.Hostname()+":80"))
	if err != nil {
		return headers, body, err
	}

	defer connection.Close()

	// send
	_, err = connection.Write([]byte(fmt.Sprintf("GET /index.html HTTP/1.0\r\nHost: %s\r\n\r\n", u.Hostname())))
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
	bytes := c.Buffered()
	buffer := make([]byte, bytes)
	_, err = c.Read(buffer)
	if err != nil {
		return headers, body, err
	}
	body = string(buffer)

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

func main() {
	_, body, err := request("http://example.org")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	show(body)
}

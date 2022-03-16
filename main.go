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
	network   = "tcp"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"
)

var (
	ErrUnsupportedURLScheme = errors.New("unsupported URL scheme")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	supportedSchemes        = map[string]struct{}{"http": struct{}{}, "https": struct{}{}, "file": struct{}{}, "data": struct{}{}}
	htmlEntities            = map[string]string{"&lt;": "<", "&gt;": ">"}
)

func request(requestURL string, additionalRequestHeaders map[string]string) (map[string]string, string, error) {
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

	// handle file scheme
	if u.Scheme == "file" {
		contents, err := openLocalFile(strings.TrimPrefix(requestURL, "file://"))
		if err != nil {
			return headers, body, err
		}
		return headers, contents, nil
	}

	// handle data scheme
	if u.Scheme == "data" {
		s := strings.SplitN(strings.TrimPrefix(requestURL, "data:"), ",", 2)
		mediaType, data := s[0], s[1]
		if mediaType != "text/html" && mediaType != "" {
			return headers, data, ErrUnsupportedMediaType
		}
		return headers, data, nil
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
	{
		requestHeaders := map[string]string{
			"Host":       u.Hostname(),
			"Connection": "close",
		}
		// merge headers; overwrite if necessary
		for k, v := range additionalRequestHeaders {
			requestHeaders[k] = v
		}

		for h, v := range requestHeaders {
			request += fmt.Sprintf("%s: %s\r\n", h, v)
		}
		request += "\r\n"
	}

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

func show(writer io.Writer, body string) {
	var inTag bool
	var inBody bool
	var inEntity bool
	var tagName string
	var entityName string

	for _, r := range body {
		c := string(r)

		switch true {
		// tags
		case c == "<":
			inTag = true
			if tagName == "/body" {
				inBody = false
			}
		case c == ">":
			inTag = false
			if tagName == "body" {
				inBody = true
			}
			tagName = ""
		case inTag:
			tagName += c
		// entities
		case c == "&":
			inEntity = true
			entityName += c
		case c == ";":
			entityName += c
			character := htmlEntities[entityName]
			fmt.Fprint(writer, character)
			inEntity = false
			entityName = ""
		case inBody && inEntity:
			entityName += c
		// body
		case inBody:
			fmt.Fprint(writer, c)
		}
	}
}

func load(url string) {
	_, body, err := request(url, map[string]string{"User-Agent": userAgent})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	show(os.Stdout, body)
}

func openLocalFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), err
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please input a URL")
		os.Exit(1)
	}

	load(args[1])
}

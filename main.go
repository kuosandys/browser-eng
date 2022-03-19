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
	network      = "tcp"
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"
	maxRedirects = 20
)

var (
	ErrUnsupportedURLScheme  = errors.New("unsupported URL scheme")
	ErrUnsupportedMediaType  = errors.New("unsupported media type")
	ErrMissingLocationHeader = errors.New("missing Location header")
	ErrMaxRedirectsReached   = errors.New("max redirects reached")
	supportedSchemes         = map[string]struct{}{"http": struct{}{}, "https": struct{}{}, "file": struct{}{}, "data": struct{}{}, "view-source": struct{}{}}
	htmlEntities             = map[string]string{"&lt;": "<", "&gt;": ">"}
)

func request(u *url.URL, additionalRequestHeaders map[string]string, redirected int) (map[string]string, string, error) {
	var err error
	var headers = map[string]string{}
	var body string

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
	req := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
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
			req += fmt.Sprintf("%s: %s\r\n", h, v)
		}
		req += "\r\n"
	}

	_, err = connection.Write([]byte(req))
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
		headers[strings.TrimSpace(header)] = strings.TrimSpace(strings.ToLower(value))
	}

	// handle redirect (status 3xx)
	if string(status[0]) == "3" {
		if redirected == maxRedirects {
			return headers, body, ErrMaxRedirectsReached
		}

		location, ok := headers["Location"]
		if !ok {
			return headers, body, ErrMissingLocationHeader
		}

		// missing host and scheme
		if string(location[0]) == "/" {
			location = u.Scheme + "://" + u.Hostname() + location
		}

		redirectURL, err := url.Parse(location)
		if err != nil {
			return headers, body, err
		}

		headers, body, err = request(redirectURL, map[string]string{}, redirected+1)
		return headers, body, err
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
			if strings.Contains(tagName, "body") {
				inBody = true
			}
			tagName = ""
		case inTag:
			tagName += c
		// entities
		case c == "&":
			inEntity = true
			entityName += c
		case inEntity && c == ";":
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

func transform(body string) string {
	bodyTransformed := "<body>"

	for _, r := range body {
		c := string(r)
		switch true {
		case c == "<":
			bodyTransformed += "&lt;"
		case c == ">":
			bodyTransformed += "&gt;"
		default:
			bodyTransformed += c
		}
	}

	return bodyTransformed + "</body>"
}

func openLocalFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), err
}

func load(requestURL string) error {
	var err error
	u, err := url.Parse(requestURL)
	if err != nil {
		return err
	}

	if _, ok := supportedSchemes[u.Scheme]; !ok {
		return ErrUnsupportedURLScheme
	}

	var data string
	switch u.Scheme {
	case "http":
	case "https":
		_, data, err = request(u, map[string]string{"User-Agent": userAgent}, 0)
		if err != nil {
			return err
		}

		show(os.Stdout, data)
	case "file":
		data, err = openLocalFile(strings.TrimPrefix(requestURL, "file://"))
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, data)
	case "data":
		s := strings.SplitN(strings.TrimPrefix(requestURL, "data:"), ",", 2)
		mediaType, data := s[0], s[1]
		switch mediaType {
		case "text/html":
			show(os.Stdout, data)
		case "":
			fmt.Fprint(os.Stdout, data)
		default:
			return ErrUnsupportedMediaType
		}
	case "view-source":
		u, err = url.Parse(strings.TrimPrefix(requestURL, "view-source:"))
		if err != nil {
			return err
		}
		_, data, err = request(u, map[string]string{"User-Agent": userAgent}, 0)
		if err != nil {
			return err
		}
		show(os.Stdout, transform(data))
	}

	return nil

}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please input a URL")
		os.Exit(1)
	}

	err := load(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

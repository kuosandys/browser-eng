package requester

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"

	customErrors "github.com/kuosandys/browser-engineering/pkg/errors"
	"github.com/kuosandys/browser-engineering/pkg/parser"
)

const (
	network      = "tcp"
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"
	maxRedirects = 20
)

var (
	supportedSchemes = map[string]struct{}{"http": {}, "https": {}, "file": {}, "data": {}, "view-source": {}}
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
		path = "/"
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

	defer func(connection net.Conn) {
		err = connection.Close()
	}(connection)

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
		headers[strings.ToLower(strings.TrimSpace(header))] = strings.TrimSpace(strings.ToLower(value))
	}

	// handle redirect (status 3xx)
	if string(status[0]) == "3" {
		if redirected == maxRedirects {
			return headers, body, customErrors.ErrMaxRedirectsReached
		}

		location, ok := headers["location"]
		if !ok {
			return headers, body, customErrors.ErrMissingLocationHeader
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

func MakeRequest(requestURL string) ([]interface{}, error) {
	var err error
	var text []interface{}

	u, err := url.Parse(requestURL)
	if err != nil {
		return text, err
	}

	if _, ok := supportedSchemes[u.Scheme]; !ok {
		return text, customErrors.ErrUnsupportedURLScheme
	}

	switch u.Scheme {
	case "http":
	case "https":
		_, data, err := request(u, map[string]string{"User-Agent": userAgent}, 0)
		if err != nil {
			return text, err
		}
		text = parser.Lex(data)
	case "file":
		data, err := openLocalFile(strings.TrimPrefix(requestURL, "file://"))
		if err != nil {
			return text, err
		}
		text = parser.Lex(data)
	case "data":
		s := strings.SplitN(strings.TrimPrefix(requestURL, "data:"), ",", 2)
		mediaType, data := s[0], s[1]
		if mediaType != "text/html" && mediaType != "" {
			return text, customErrors.ErrUnsupportedMediaType
		}
		text = parser.Lex(data)
	case "view-source":
		u, err = url.Parse(strings.TrimPrefix(requestURL, "view-source:"))
		if err != nil {
			return text, err
		}
		_, data, err := request(u, map[string]string{"User-Agent": userAgent}, 0)
		if err != nil {
			return text, err
		}
		text = parser.Lex(parser.Transform(data))
	}

	return text, nil
}

func openLocalFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), err
}

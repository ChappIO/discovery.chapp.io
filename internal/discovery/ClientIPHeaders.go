package discovery

import (
	"net"
	"net/http"
	"strings"
)

type ClientIPHeaders []string

func (c *ClientIPHeaders) GetClientIP(request *http.Request) string {
	for _, headerName := range *c {
		if headerValue := request.Header.Get(headerName); headerValue != "" {
			if strings.Contains(headerValue, ":") {
				if host, _, err := net.SplitHostPort(headerValue); err == nil {
					return host
				}
			}
			return headerValue
		}
	}

	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return host
	}

	return ""
}

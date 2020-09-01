package discovery

import (
	"net"
	"net/http"
)

type ClientIPHeaders []string

func (c *ClientIPHeaders) GetClientIP(request *http.Request) string {
	for _, headerName := range *c {
		if headerValue := request.Header.Get(headerName); headerValue != "" {
			return headerValue
		}
	}

	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return host
	}

	return ""
}

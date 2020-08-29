package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
)

type Agent struct {
	ID        string `json:"id"`
	PrivateIP string `json:"private_ip"`
}

type Response struct {
	ID       string  `json:"id"`
	PublicIP string  `json:"public_ip"`
	Agents   []Agent `json:"agents"`
}

var ipHeaders = []string{"x-forwarded-for"}

func getClientIp(request *http.Request) string {
	for _, headerName := range ipHeaders {
		if headerValue := request.Header.Get(headerName); headerValue != "" {
			return headerValue
		}
	}

	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return host
	}

	return ""
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		response := Response{
			ID:       strings.TrimPrefix(request.RequestURI, "/"),
			PublicIP: getClientIp(request),
			Agents:   []Agent{},
		}
		_ = json.NewEncoder(writer).Encode(&response)
	})

	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "3000"
	}

	if err := http.ListenAndServe(":"+listenPort, nil); err != http.ErrServerClosed {
		panic(err)
	}
}

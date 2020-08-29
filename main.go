package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
)

type Response struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

func getClientIp(request *http.Request) string {
	publicIp := ""

	// default to the remote address
	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		publicIp = host
	}

	return publicIp
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		response := Response{OK: true, Data: map[string]interface{}{
			"headers":  request.Header,
			"remote":   request.RemoteAddr,
			"publicIp": getClientIp(request),
		}}
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
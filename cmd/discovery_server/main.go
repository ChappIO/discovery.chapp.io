package main

import (
	"discovery.chapp.io/internal/discovery"
	"net/http"
	"os"
)

func main() {

	server := discovery.NewServer()

	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "3000"
	}

	if err := http.ListenAndServe(":"+listenPort, server); err != http.ErrServerClosed {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Agent struct {
	AgentID        string `json:"agent_id"`
	PrivateAddress string `json:"private_address"`
}

type Response struct {
	ServiceID string   `json:"service_id"`
	PublicIP  string   `json:"public_ip"`
	Agents    []*Agent `json:"agents"`
}

var ipHeaders = []string{"x-forwarded-for"}

var knownAgents = map[string][]*Agent{}
var lock = sync.Mutex{}

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

func registerAgent(serviceId string, agentId string, privateAddress string) {
	lock.Lock()
	defer lock.Unlock()

	if len(agentId) > 36 {
		agentId = agentId[:36]
	}

	// This is a new agent! Register it.
	agents, ok := knownAgents[serviceId]
	if !ok {
		agents = make([]*Agent, 0)
	}
	for _, agent := range agents {
		if agent.AgentID == agentId {
			// we already know this agent... update it
			agent.PrivateAddress = privateAddress
			return
		}
	}
	// this is a new agent, register it
	knownAgents[serviceId] = append(agents, &Agent{
		AgentID:        agentId,
		PrivateAddress: privateAddress,
	})
}

func redirectToDocs(writer http.ResponseWriter) {
	writer.Header().Set("Location", "https://github.com/ChappIO/discovery.chapp.io")
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// allow cross domain AJAX requests
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		requestPath := strings.TrimPrefix(request.URL.Path, "/")
		if requestPath == "" {
			// someone simply opened the root url. Redirect them to docs
			redirectToDocs(writer)
			return
		}

		response := Response{
			ServiceID: requestPath,
			PublicIP:  getClientIp(request),
		}

		serviceId := response.ServiceID + "/" + response.PublicIP

		if privateAddress := request.URL.Query().Get("private_address"); privateAddress != "" {
			if _, _, err := net.SplitHostPort(privateAddress); err == nil {
				// this address is valid
				if agentId := request.URL.Query().Get("agent_id"); agentId != "" {
					registerAgent(serviceId, agentId, privateAddress)
				}
			}
		}

		response.Agents = knownAgents[serviceId]
		if response.Agents == nil {
			response.Agents = []*Agent{}
		}
		log.Printf("From %s found %d agents", serviceId, len(response.Agents))
		writer.Header().Set("Content-Type", "application/json")
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

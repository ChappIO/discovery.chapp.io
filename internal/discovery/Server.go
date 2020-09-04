package discovery

import (
	"discovery.chapp.io/internal/storage"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

type Response struct {
	ServiceID string          `json:"service_id"`
	PublicIP  string          `json:"public_ip"`
	Agents    []storage.Agent `json:"agents"`
}

type Server interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type coreServer struct {
	clientIp *ClientIPHeaders
	store    storage.Storage
}

func (s *coreServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	requestPath := strings.TrimPrefix(request.URL.Path, "/")

	if requestPath == "" {
		s.serveDocs(writer)
	} else {
		s.serveServiceId(requestPath, writer, request)
	}
}

func (s *coreServer) serveDocs(writer http.ResponseWriter) {
	writer.Header().Set("Location", "https://github.com/ChappIO/discovery.chapp.io")
	writer.WriteHeader(http.StatusPermanentRedirect)
}

func (s *coreServer) serveServiceId(serviceId string, writer http.ResponseWriter, request *http.Request) {
	clientIp := s.clientIp.GetClientIP(request)

	var result []storage.Agent

	params := request.URL.Query()
	privateAddress := params.Get("private_address")
	agentId := params.Get("agent_id")
	if privateAddress != "" && agentId != "" {
		// this is an 'add agent' request
		if _, _, err := net.SplitHostPort(privateAddress); err == nil {
			// the private address is valid
			// limit agentId length to UUIDv4 length
			if len(agentId) > 36 {
				agentId = agentId[:36]
			}
			result = s.store.Add(clientIp, serviceId, storage.Agent{
				AgentID:        agentId,
				PrivateAddress: privateAddress,
			})
		}
	}
	if result == nil {
		result = s.store.Get(clientIp, serviceId)
	}

	response := Response{
		ServiceID: serviceId,
		PublicIP:  clientIp,
		Agents:    result,
	}

	log.Printf("From %s found %d agents", serviceId, len(response.Agents))
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(&response)
}

func NewServer() Server {
	return &coreServer{
		clientIp: &ClientIPHeaders{
			"x-forwarded-for",
		},
		store: storage.InMemory(),
	}
}

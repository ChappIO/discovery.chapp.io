package storage


type Agent struct {
	AgentID        string `json:"agent_id"`
	PrivateAddress string `json:"private_address"`
}

type Storage interface {
	Add(clientId string, serviceId string, agent Agent) []Agent
	Get(client string, serviceId string) []Agent
}

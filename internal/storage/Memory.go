package storage

import (
	"sync"
)

type inMemoryStorage struct {
	data map[string][]Agent
	lock sync.Mutex
}

func (i *inMemoryStorage) Add(clientId string, serviceId string, agent Agent) []Agent {
	i.lock.Lock()
	defer i.lock.Unlock()
	key := clientId + "/" + serviceId
	agents, ok := i.data[key]
	if !ok {
		agents = []Agent{agent}
	} else {
		index := -1
		for i, elem := range agents {
			if elem.AgentID == agent.AgentID {
				index = i
				break
			}
		}
		if index == -1 {
			agents = append(agents, agent)
		} else {
			agents[index] = agent
		}
	}
	i.data[key] = agents
	return agents
}

func (i *inMemoryStorage) Get(clientId string, serviceId string) []Agent {
	key := clientId + "/" + serviceId
	agents, ok := i.data[key]
	if !ok {
		return make([]Agent, 0)
	} else {
		return agents
	}
}

func InMemory() Storage {
	return &inMemoryStorage{
		data: map[string][]Agent{},
	}
}

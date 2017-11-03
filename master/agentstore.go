package master

import (
	"io/ioutil"
	"log"

	"github.com/ricottadeploy/ricotta/security"
	yaml "gopkg.in/yaml.v2"
)

type Agent struct {
	Env   []string
	Id    string
	Roles []string
}

type AgentStore struct {
	agentsMap map[security.Fingerprint]Agent
}

func NewAgentStore() AgentStore {
	return AgentStore{}
}

func (store *AgentStore) ReadFromYamlFile(agentsFile string) {
	a, err := ioutil.ReadFile(agentsFile)
	if err != nil {
		log.Fatalf("Error parsing agents.yaml: %s", err)
	}
	err1 := yaml.Unmarshal(a, &store.agentsMap)
	if err1 != nil {
		log.Fatalf("Error parsing agents file %s: %s", agentsFile, err1)
	}
	for k := range store.agentsMap {
		if !k.Valid() {
			log.Fatalf("Error parsing agents file %s: Invalid fingerprint %s", agentsFile, k)
		}
	}
}

func (store *AgentStore) ToYaml() string {
	b, err := yaml.Marshal(store.agentsMap)
	if err != nil {
		log.Fatalf("Error while marshalling: %s", err)
	}
	return string(b)
}

func (store *AgentStore) Add(fingerprint security.Fingerprint, agent Agent) {
	if !fingerprint.Valid() {
		log.Fatalf("Invalid fingerprint %s: ", fingerprint)
	}
	store.agentsMap[fingerprint] = agent
}

func (store *AgentStore) Remove(fingerprint security.Fingerprint) {
	if !fingerprint.Valid() {
		log.Fatalf("Invalid fingerprint %s: ", fingerprint)
	}
	delete(store.agentsMap, fingerprint)
}

func (store *AgentStore) Get(fingerprint security.Fingerprint) (Agent, bool) {
	agent, ok := store.agentsMap[fingerprint]
	return agent, ok
}

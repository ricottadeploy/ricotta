package master

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Agent struct {
	Env   []string
	Id    string
	Roles []string
}

type AgentStore struct {
	agentsMap map[Fingerprint]Agent
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
	for k, _ := range store.agentsMap {
		if !k.Valid() {
			log.Fatalf("Error parsing agents file %s: Invalid fingerprint %s", agentsFile, k)
		}
	}
}

func (store *AgentStore) ToYaml() string {
	b, err := yaml.Marshal(store.agentsMap)
	if err != nil {
		log.Fatal("Error while marshalling: %s", err)
	}
	return string(b)
}

func (store *AgentStore) Add(fingerprint Fingerprint, agent Agent) {
	if !fingerprint.Valid() {
		log.Fatal("Invalid fingerprint %s: ", fingerprint)
	}
	store.agentsMap[fingerprint] = agent
}

func (store *AgentStore) Remove(fingerprint Fingerprint) {
	if !fingerprint.Valid() {
		log.Fatal("Invalid fingerprint %s: ", fingerprint)
	}
	delete(store.agentsMap, fingerprint)
}

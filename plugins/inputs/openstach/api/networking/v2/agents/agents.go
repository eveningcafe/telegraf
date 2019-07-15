package agents

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/networking/v2"
)

type Agent struct {
	ID               string
	Binary           string
	Host             string
	AdminStateUp     bool
	Configurations   interface{}
	Alive            bool
	AgentType        string
	AvailabilityZone interface{}
	Topic            string
}
func List(client *v2.NetworkClient) ([]Agent, error) {
	api, err := declareListAgent(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListAgentResponse{}
	err = json.Unmarshal([]byte(api.Response),&result)
	agents := []Agent{}
	for _, v := range result.Agents {
		agents = append(agents, Agent{
			ID: v.ID,
			Binary: v.Binary,
			Host: v.Host,
			AdminStateUp: v.AdminStateUp,
			Configurations: v.Configurations,
			Alive: v.Alive,
			AgentType: v.AgentType,
			AvailabilityZone : v.AvailabilityZone,
			Topic: v.Topic,
		})
	}
	return agents, err
}

package agents

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
)

type Agent struct {
	ID               string
	Binary           string
	Host             string
	AdminStateUp     bool
	Configurations   interface{}
	Alive            bool
	AgentType        string
	AvailabilityZone string
	Topic            string
}
func List(client *v2.NetworkClient) ([]Agent, error) {
	api, err := declareListAgent(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err!=nil {
		return nil, err
	}
	result := ListAgentResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	agents := []Agent{}
	for _, v := range result.Agents {
		zone := "unknow"
		if v.AvailabilityZone != nil{
			zone = v.AvailabilityZone.(string)
		}
		agents = append(agents, Agent{
			ID: v.ID,
			Binary: v.Binary,
			Host: v.Host,
			AdminStateUp: v.AdminStateUp,
			Configurations: v.Configurations,
			Alive: v.Alive,
			AgentType: v.AgentType,
			AvailabilityZone : zone,
			Topic: v.Topic,
		})
	}
	return agents, err
}

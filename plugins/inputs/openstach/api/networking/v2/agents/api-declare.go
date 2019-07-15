package agents

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type ListAgentRequest struct {
}

type ListAgentResponse struct {
	Agents []struct {
		Binary             string      `json:"binary"`
		Description        interface{} `json:"description"`
		AvailabilityZone   interface{} `json:"availability_zone"`
		HeartbeatTimestamp string      `json:"heartbeat_timestamp"`
		AdminStateUp       bool        `json:"admin_state_up"`
		Alive              bool        `json:"alive"`
		ID                 string      `json:"id"`
		Topic              string      `json:"topic"`
		Host               string      `json:"host"`
		AgentType          string      `json:"agent_type"`
		StartedAt          string      `json:"started_at"`
		CreatedAt          string      `json:"created_at"`
		ResourcesSynced    bool        `json:"resources_synced"`
		Configurations interface{}
	} `json:"agents"`
}

//
func declareListAgent(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListAgentRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/agents",
		HeaderRequest: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		Request: jsonBody,
	}, err
}
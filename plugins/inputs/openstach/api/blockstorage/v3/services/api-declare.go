package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type ListServiceRequest struct {
}

type ListServiceResponse struct {
	Services []struct {
		Status            string      `json:"status"`
		Binary            string      `json:"binary"`
		Zone              string      `json:"zone"`
		State             string      `json:"state"`
		UpdatedAt         string      `json:"updated_at"`
		Host              string      `json:"host"`
		DisabledReason    interface{} `json:"disabled_reason"`
		Frozen            bool        `json:"frozen,omitempty"`
		Cluster           interface{} `json:"cluster,omitempty"`
		ReplicationStatus string      `json:"replication_status,omitempty"`
		ActiveBackendID   interface{} `json:"active_backend_id,omitempty"`
	} `json:"services"`
}

//
func declareListService(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListServiceRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/os-services",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

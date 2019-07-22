package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListServiceRequest struct {
}

type ListServiceResponse struct {
	Services []struct {
		ID             int    `json:"id"`
		Binary         string `json:"binary"`
		DisabledReason string `json:"disabled_reason"`
		Host           string `json:"host"`
		State          string `json:"state"`
		Status         string `json:"status"`
		UpdatedAt      string `json:"updated_at"`
		ForcedDown     bool   `json:"forced_down"`
		Zone           string `json:"zone"`
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
package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListServiceRequest struct {
}

type ListServiceResponse struct {
	Services []Service `json:"services"`
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
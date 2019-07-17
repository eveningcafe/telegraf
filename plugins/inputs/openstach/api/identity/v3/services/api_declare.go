package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type ListServiceRequest struct {
}

type ListServiceResponse struct {
	Links struct {
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Self     string      `json:"self"`
	} `json:"links"`
	Services []struct {
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		ID          string `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"services"`
}


// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-services-detail#list-services
func declareListService(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListServiceRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/services",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}



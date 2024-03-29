package regions

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListRegionResponse struct {
	Links struct {
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Self     string      `json:"self"`
	} `json:"links"`
	Regions []struct {
		Description string `json:"description"`
		ID          string `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		ParentRegionID interface{} `json:"parent_region_id"`
	} `json:"regions"`
}
type ListRegionRequest struct {
}

type ListRegionAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  ListRegionRequest
	Response ListRegionResponse
}

// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-services-detail#list-services
func declareListRegion(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListRegionRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/regions",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

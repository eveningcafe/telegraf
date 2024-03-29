package groups

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListGroupRequest struct {
}

type ListGroupResponse struct {
	Links struct {
		Self     string      `json:"self"`
		Previous interface{} `json:"previous"`
		Next     interface{} `json:"next"`
	} `json:"links"`
	Groups []struct {
		Description string `json:"description"`
		DomainID    string `json:"domain_id"`
		ID          string `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		Name string `json:"name"`
	} `json:"groups"`
}


// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-services-detail#list-services
func declareListGroup(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListGroupRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/groups",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}



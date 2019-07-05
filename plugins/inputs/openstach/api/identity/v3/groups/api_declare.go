package groups

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base"
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
func declareListGroup(endpoint string, token string) (*base.OpenstackAPI, error) {
	req := ListGroupRequest{}
	jsonBody, err := json.Marshal(req)
	return &base.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/groups",
		HeaderRequest : map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		Request: jsonBody,
	}, err
}



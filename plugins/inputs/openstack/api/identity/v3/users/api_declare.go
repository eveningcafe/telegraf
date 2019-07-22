package users

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListUserResponse struct {
	Links struct {
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Self     string      `json:"self"`
	} `json:"links"`
	Users []struct {
		DomainID string `json:"domain_id"`
		Enabled  bool   `json:"enabled"`
		ID       string `json:"id"`
		Links    struct {
			Self string `json:"self"`
		} `json:"links"`
		Name              string      `json:"name"`
		PasswordExpiresAt interface{} `json:"password_expires_at"`
	} `json:"users"`
}

type ListUserRequest struct {
}

type ListUserAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  ListUserRequest
	Response ListUserResponse
}

// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-services-detail#list-services
func declareListUser(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListUserRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/users",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}


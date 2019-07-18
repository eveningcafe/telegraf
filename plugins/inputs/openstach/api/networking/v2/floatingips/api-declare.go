package floatingips

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type ListFloatingIpRequest struct {
}

type ListFloatingIpResponse struct {
	Floatingips []Floatingip `json:"floatingips"`
}


//
func declareListFloatingIp(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListFloatingIpRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/v2.0/floatingips",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

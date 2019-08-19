package scheduler

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListPoolRequest struct {
}

type ListPoolResponse struct {
	Pools []struct {
		Name         string `json:"name"`
		Capabilities Capabilities
	} `json:"pools"`
}


//
func declareListPool(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListPoolRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/scheduler-stats/get_pools",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody:        jsonBody,
		RequestParameter: map[string]string{
			"detail" : "true",
		},
		RequestParameterRequire: true,

	}, err
}
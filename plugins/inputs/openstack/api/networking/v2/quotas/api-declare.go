package quotas

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type DetailQuotasRequest struct {
}

type DetailQuotasResponse struct {
	Quota Quota `json:"quota"`
}

//
func declareQuotasDetail(endpoint string, token string, projectID string) (*request.OpenstackAPI, error) {
	req := DetailQuotasRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/v2.0/quotas/"+projectID+"/details/",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

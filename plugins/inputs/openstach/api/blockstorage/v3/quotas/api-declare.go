package quotas

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type DetailQuotasRequest struct {
}

type DetailQuotasResponse struct {
	QuotaSet Quota `json:"quota_set"`
}

//
func declareQuotasDetail(endpoint string, token string, projectID string) (*request.OpenstackAPI, error) {
	req := DetailQuotasRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/os-quota-sets/"+projectID,
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
		RequestParameter: map[string]string{
			"usage": "True",
		},
		RequestParameterRequire: true,
	}, err
}

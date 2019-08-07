package hypervisors
import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListHypervisorRequest struct {
}


type ListHypervisorResponse struct {
	Hypervisors [] Hypervisor `json:"hypervisors"`
	HypervisorsLinks []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"hypervisors_links"`
}

//
func declareListHypervisor(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListHypervisorRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/os-hypervisors/detail",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-OpenStack-Nova-API-Version": request.XOpenStackNovaAPIv2Version,
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

package networks

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListNetworkRequest struct {
}

type ListNetworkResponse struct {
	Networks []Network `json:"networks"`
}
//
func declareListNetwork(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListNetworkRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/v2.0/networks",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

//
type NetworkIPAvailabilitiesRequest struct{

}
type NetworkIPAvailabilitiesResponse struct {
	NetworkIPAvailabilities []IPAvailabilities`json:"network_ip_availabilities"`
}
//
func declareNetworkIPAvailabilities(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := NetworkIPAvailabilitiesRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/v2.0/network-ip-availabilities",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}
package resource_providers

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListResourceProvidersRequest struct {
}

type ListResourceProvidersResponse struct {
	ResourceProviders []ResourceProviders `json:"resource_providers"`
}

//
type GetResourceProviderInventoriesRequest struct {
}

type GetResourceProviderInventoriesResponse struct {
	ResourceProviderGeneration int                         `json:"resource_provider_generation"`
	Inventories                ResourceProviderInventories `json:"inventories"`
}

//
type GetResourceProviderUsagesRequest struct {}
type GetResourceProviderUsagesResponse struct {
	ResourceProviderGeneration int                    `json:"resource_provider_generation"`
	Usages                     ResourceProviderUsages `json:"usages"`
}

func declareListResource(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListResourceProvidersRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/resource_providers",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

//
func declareGetResourceProviderInventories(endpoint string, token string, resourceID string) (*request.OpenstackAPI, error) {
	req := GetResourceProviderInventoriesRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/resource_providers/" + resourceID + "/inventories",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

func declareGetResourceProviderUsages(endpoint string, token string, resourceID string) (*request.OpenstackAPI, error) {
	req := GetResourceProviderUsagesRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/resource_providers/" + resourceID + "/usages",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
	}, err
}

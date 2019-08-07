package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
)

// Service represents an OpenStack Service.
type Service struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
	Links       struct {
		Self string `json:"self"`
	} `json:"links"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func List(client *v3.IdentityClient) ([]Service, error) {
	api, err := declareListService(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return nil,err
	}
	result := ListServiceResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	services := []Service{}
	for _, v := range result.Services {
		services = append(services, v)
	}

	return services, err
}

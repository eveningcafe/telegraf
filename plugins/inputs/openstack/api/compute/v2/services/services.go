package services

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2"
)

type Service struct {
	ID             int    `json:"id"`
	Binary         string `json:"binary"`
	DisabledReason string `json:"disabled_reason"`
	Host           string `json:"host"`
	State          string `json:"state"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at"`
	ForcedDown     bool   `json:"forced_down"`
	Zone           string `json:"zone"`
}

func List(client *v2.ComputeClient) ([]Service, error) {
	api, err := declareListService(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return []Service{},err
	}
	if (err != nil) {
		return nil, err
	}
	result := ListServiceResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	services := []Service{}
	for _, v := range result.Services {
		services = append(services,v)
	}
	return services, err
}

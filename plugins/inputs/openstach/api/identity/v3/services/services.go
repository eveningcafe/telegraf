package services

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
)
// Service represents an OpenStack Service.
type Service struct {
	// ID is the unique ID of the service.
	ID string `json:"id"`

	// Type is the type of the service.
	Type string `json:"type"`

	// Enabled is whether or not the service is enabled.
	Enabled bool `json:"enabled"`

	// Links contains referencing links to the service.
	Links struct{ Self string `json:"self"` }
}


func List(client *v3.IdentityClient) ([]Service, error) {
	api, err := declareListService(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListServiceResponse{}
	err = json.Unmarshal([]byte(api.Response),&result)
	services := []Service{}
	for _, v := range result.Services {
		services = append(services, Service{
			ID: v.ID,
			Type: v.Type,
			Enabled: v.Enabled,
			Links: v.Links,
		})
	}

	return services, err
}
package regions

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
)

type Region struct {
	// DomainID is the domain ID the user belongs to.
	DomainID string `json:"domain_id"`

	// Enabled is whether or not the user is enabled.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the user.
	ID string `json:"id"`

	// Name is the name of the user.
	Name           string `json:"name"`
	Description    interface{}
	ParentRegionID interface{}
}

func List(client *v3.IdentityClient) ([]Region, error) {
	api, err := declareListRegion(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListRegionResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	regions := []Region{}
	for _, v := range result.Regions {
		regions = append(regions, Region{
			ID:            v.ID,
			Description:  v.Description,
			ParentRegionID:     v.ParentRegionID,
		})
	}

	return regions, err
}


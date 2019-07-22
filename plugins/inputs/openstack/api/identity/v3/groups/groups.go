package groups

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
)

// Service represents an OpenStack Service.
type Group struct {
	// ID is the unique ID of the service.
	ID          string `json:"id"`
	Description interface{}
	Name        string
	DomainID    string
}

func List(client *v3.IdentityClient) ([]Group, error) {
	api, err := declareListGroup(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListGroupResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	groups := []Group{}
	for _, v := range result.Groups {
		groups = append(groups, Group{
			ID:          v.ID,
			Description: v.Description,
			Name:        v.Name,
			DomainID:    v.DomainID,
		})
	}

	return groups, err
}

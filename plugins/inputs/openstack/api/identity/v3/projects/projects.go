package projects

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
)
// Project represents an OpenStack Identity Project.
type Project struct {
	// IsDomain indicates whether the project is a domain.
	IsDomain bool `json:"is_domain"`

	// DomainID is the domain ID the project belongs to.
	DomainID string `json:"domain_id"`

	// Enabled is whether or not the project is enabled.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the project.
	ID string `json:"id"`

	// Name is the name of the project.
	Name string `json:"name"`

}


func List(client *v3.IdentityClient) ([]Project, error) {
	api, err := declareListProject(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return []Project{},err
	}
	result := ListProjectResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	projects := []Project{}
	for _, v := range result.Projects {
		projects = append(projects, Project{
			ID: v.ID,
			Enabled: v.Enabled,
			IsDomain: v.IsDomain,
			DomainID: v.DomainID,
			Name: v.Name,
		})
	}

	return projects, err
}

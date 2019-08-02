package users

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
)

type User struct {

	// DomainID is the domain ID the user belongs to.
	DomainID string `json:"domain_id"`

	// Enabled is whether or not the user is enabled.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the user.
	ID string `json:"id"`

	// Name is the name of the user.
	Name string `json:"name"`

}

func List(client *v3.IdentityClient) ([]User, error) {
	api, err := declareListUser(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return []User{},err
	}
	result := ListUserResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	users := []User{}
	for _, v := range result.Users {
		users = append(users, User{
			ID: v.ID,
			Enabled: v.Enabled,
			DomainID: v.DomainID,
			Name: v.Name,
		})
	}

	return users, err
}
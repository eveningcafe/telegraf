package projects

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
	"io/ioutil"
	"net/http"
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

func List(client *v3.IdentityClient) ([]Project, error){
	api := declareListProject(client.Token)

	jsonBody, err := json.Marshal(api.Request)

	if err != nil {
		panic(err.Error())
	}

	httpClient := &http.Client{}
	request, err := http.NewRequest(api.Method, client.Endpoint+api.Path, bytes.NewBuffer(jsonBody))
	for k, v := range api.Header {
		request.Header.Add(k,v)
	}
	resp, err := httpClient.Do(request)
	defer resp.Body.Close()

	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("List service successful ")
	} else {
		err := errors.New("List service respond status code "+ string(resp.StatusCode))
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal([]byte(body), &api.Response)

	projects := []Project{}
	for _,v := range api.Response.Projects{
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

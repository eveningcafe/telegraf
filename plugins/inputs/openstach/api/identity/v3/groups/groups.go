package groups

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
	"io/ioutil"
	"net/http"
)

// Service represents an OpenStack Service.
type Group struct {
	// ID is the unique ID of the service.
	ID string `json:"id"`

	Description interface{}
	Name        string
	DomainID    string
}

func List(client *v3.IdentityClient) ([]Group, error){
	api := declareListGroup(client.Token)

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

	services := []Group{}
	for _,v := range api.Response.Groups{
		services = append(services, Group{
			ID: v.ID,
			Description: v.Description,
			Name: v.Name,
			DomainID: v.DomainID,
		})
	}

	return services, err
}
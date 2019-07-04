package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
	"io/ioutil"
	"net/http"
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

func List(client *v3.IdentityClient) ([]User, error){
	api := declareListUser(client.Token)

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

	users := []User{}
	for _,v := range api.Response.Users{
		users = append(users, User{
			ID: v.ID,
			Enabled: v.Enabled,
			DomainID: v.DomainID,
			Name: v.Name,
		})
	}

	return users, err
}
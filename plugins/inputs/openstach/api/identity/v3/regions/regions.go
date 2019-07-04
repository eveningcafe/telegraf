package regions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
	"io/ioutil"
	"net/http"
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
	api := declareListRegion(client.Token)

	jsonBody, err := json.Marshal(api.Request)

	if err != nil {
		panic(err.Error())
	}

	httpClient := &http.Client{}
	request, err := http.NewRequest(api.Method, client.Endpoint+api.Path, bytes.NewBuffer(jsonBody))
	for k, v := range api.Header {
		request.Header.Add(k, v)
	}
	resp, err := httpClient.Do(request)
	defer resp.Body.Close()

	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("List service successful ")
	} else {
		err := errors.New("List service respond status code " + string(resp.StatusCode))
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal([]byte(body), &api.Response)

	regions := []Region{}
	for _, v := range api.Response.Regions {
		regions = append(regions, Region{
			ID:            v.ID,
			Description:  v.Description,
			ParentRegionID:     v.ParentRegionID,
		})
	}

	return regions, err
}

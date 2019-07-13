package services

import v2 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/compute/v2"

type Service struct{
	
}
func List(client *v2.ComputeClient) ([]Service, error) {
	api, err := declareListService(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListServiceResponse{}
	err = json.Unmarshal([]byte(api.Response),&result)
	users := []Service{}
	for _, v := range result.Services {
		users = append(users, Service{
			ID: v.ID,
			Type: v.Type,
			Enabled: v.Enabled,
			Links: v.Links,
		})
	}

	return users, err
}
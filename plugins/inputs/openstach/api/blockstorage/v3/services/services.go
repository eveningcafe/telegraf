package services

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3"
)

type Service struct {
	ID                string
	Binary            string
	Status            string
	Host              string
	State             string
	UpdatedAt         string
	DisabledReason    interface{}
	ActiveBackendID   interface{}
	ReplicationStatus string
	Zone              string
}

func List(client *v3.VolumeClient) ([]Service, error) {
	api, err := declareListService(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListServiceResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	services := []Service{}
	for _, v := range result.Services {
		services = append(services, Service{
			ID: v.Binary+v.Host,
			Binary: v.Binary,
			Status: v.Status,
			Host: v.Host,
			State: v.State,
			UpdatedAt: v.UpdatedAt,
			DisabledReason: v.DisabledReason,
			ActiveBackendID: v.ActiveBackendID,
			ReplicationStatus: v.ReplicationStatus,
			Zone: v.Zone,
		})
	}
	return services, err
}

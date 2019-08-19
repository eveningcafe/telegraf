package resource_providers

import (
	"encoding/json"
	"fmt"
	v1 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/placement/v1"
)

type ResourceProviders struct {
	Generation int    `json:"generation"`
	UUID       string `json:"uuid"`
	Links      []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Name               string `json:"name"`
	ParentProviderUUID string `json:"parent_provider_uuid"`
	RootProviderUUID   string `json:"root_provider_uuid"`
}
type ResourceProviderInventories struct {
	VCPU struct {
		Total           float64     `json:"total"`
		Reserved        float64     `json:"reserved"`
		MinUnit         float64     `json:"min_unit"`
		MaxUnit         float64     `json:"max_unit"`
		StepSize        float64     `json:"step_size"`
		AllocationRatio float64 `json:"allocation_ratio"`
	} `json:"VCPU"`
	MEMORYMB struct {
		Total           float64     `json:"total"`
		Reserved        float64     `json:"reserved"`
		MinUnit         float64     `json:"min_unit"`
		MaxUnit         float64     `json:"max_unit"`
		StepSize        float64     `json:"step_size"`
		AllocationRatio float64 `json:"allocation_ratio"`
	} `json:"MEMORY_MB"`
	DISKGB struct {
		Total           float64     `json:"total"`
		Reserved        float64     `json:"reserved"`
		MinUnit         float64     `json:"min_unit"`
		MaxUnit         float64     `json:"max_unit"`
		StepSize        float64     `json:"step_size"`
		AllocationRatio float64 `json:"allocation_ratio"`
	} `json:"DISK_GB"`
}

type ResourceProviderUsages struct {
	VCPU     float64 `json:"VCPU"`
	MEMORYMB float64 `json:"MEMORY_MB"`
	DISKGB   float64 `json:"DISK_GB"`
}

func List(client *v1.PlacementClient) ([]ResourceProviders, error) {
	api, err := declareListResource(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return nil, err
	}
	result := ListResourceProvidersResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	if err !=nil{
		return nil,fmt.Errorf("unable to get resource info info: request placement api GET %s - %v",api.Path, err)
	}
	resources := []ResourceProviders{}
	for _, v := range result.ResourceProviders {
		resources = append(resources, v)
	}
	return resources, err
}

func GetInventories(client *v1.PlacementClient, resourceID string) (ResourceProviderInventories, error) {
	api, err := declareGetResourceProviderInventories(client.Endpoint, client.Token, resourceID)
	err = client.DoReuest(api)
	if err != nil {
		return ResourceProviderInventories{}, err
	}
	result := GetResourceProviderInventoriesResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	if err !=nil{
		return ResourceProviderInventories{},fmt.Errorf("unable to get resource info info: request placement api GET %s - %v",api.Path, err)
	}
	resourcesInventories := result.Inventories
	return resourcesInventories, err
}

func GetUsages(client *v1.PlacementClient, resourceID string) (ResourceProviderUsages, error) {
	api, err := declareGetResourceProviderUsages(client.Endpoint, client.Token, resourceID)
	err = client.DoReuest(api)
	if err != nil {
		return ResourceProviderUsages{}, err
	}
	result := GetResourceProviderUsagesResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	resourcesInventories := result.Usages
	//for _, v := range result.ResourceProviders {
	//	resources = append(resources, v)
	//}
	return resourcesInventories, err
}

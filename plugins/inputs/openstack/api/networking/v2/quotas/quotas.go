package quotas

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
)

type Quota struct {
	RbacPolicy struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"rbac_policy"`
	Subnetpool struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"subnetpool"`
	SecurityGroupRule struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"security_group_rule"`
	SecurityGroup struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"security_group"`
	Subnet struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"subnet"`
	Port struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"port"`
	Network struct {
		Used     int `json:"used"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"network"`
}


func Detail(client *v2.NetworkClient, projectID string) (Quota, error) {
	api, err := declareQuotasDetail(client.Endpoint, client.Token, projectID)
	err = client.DoReuest(api)
	if err != nil {
		return Quota{}, err
	}
	result := DetailQuotasResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	return result.Quota, err
}

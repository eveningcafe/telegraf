package quotas

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
)

type Quota struct {
	Network struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"network"`
	Subnet struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"subnet"`
	Subnetpool struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"subnetpool"`
	Port struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"port"`
	Router struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"router"`
	Floatingip struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"floatingip"`
	RbacPolicy struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"rbac_policy"`
	SecurityGroup struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"security_group"`
	SecurityGroupRule struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	} `json:"security_group_rule"`
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

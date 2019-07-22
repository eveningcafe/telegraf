package quotas

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2"
)

type Quota struct {
	Cores struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"cores"`
	FixedIps struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"fixed_ips"`
	FloatingIps struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"floating_ips"`
	ID                       string `json:"id"`
	InjectedFileContentBytes struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"injected_file_content_bytes"`
	InjectedFilePathBytes struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"injected_file_path_bytes"`
	InjectedFiles struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"injected_files"`
	Instances struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"instances"`
	KeyPairs struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"key_pairs"`
	MetadataItems struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"metadata_items"`
	RAM struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"ram"`
	SecurityGroupRules struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"security_group_rules"`
	SecurityGroups struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"security_groups"`
	ServerGroupMembers struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"server_group_members"`
	ServerGroups struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"server_groups"`
	Networks struct {
		InUse    int `json:"in_use"`
		Limit    int `json:"limit"`
		Reserved int `json:"reserved"`
	} `json:"networks"`
}

func Detail(client *v2.ComputeClient, projectID string) (Quota, error) {
	api, err := declareQuotasDetail(client.Endpoint, client.Token, projectID)
	err = api.DoReuest()
	if err != nil {
		return Quota{}, err
	}
	result := DetailQuotasResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	return result.QuotaSet, err
}

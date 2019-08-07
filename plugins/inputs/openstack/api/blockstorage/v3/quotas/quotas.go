package quotas

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3"
)

type Quota struct {
	ID      string `json:"id"`
	Volumes struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"volumes"`
	VolumesLvmdriver1 struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"volumes_lvmdriver-1"`
	Snapshots struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"snapshots"`
	SnapshotsLvmdriver1 struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"snapshots_lvmdriver-1"`
	Backups struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"backups"`
	Groups struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"groups"`
	PerVolumeGigabytes struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"per_volume_gigabytes"`
	Gigabytes struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"gigabytes"`
	GigabytesLvmdriver1 struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"gigabytes_lvmdriver-1"`
	BackupGigabytes struct {
		Reserved  int `json:"reserved"`
		Allocated int `json:"allocated"`
		Limit     int `json:"limit"`
		InUse     int `json:"in_use"`
	} `json:"backup_gigabytes"`
}


func Detail(client *v3.VolumeClient, projectID string) (Quota, error) {
	api, err := declareQuotasDetail(client.Endpoint, client.Token, projectID)
	err = client.DoReuest(api)
	if err != nil {
		return Quota{}, err
	}
	result := DetailQuotasResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	return result.QuotaSet, err
}

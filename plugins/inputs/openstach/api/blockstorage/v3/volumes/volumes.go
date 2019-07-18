package volumes

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3"
)

type Volumes struct {
	ID                       string
	Name                     interface{}
	ReplicationStatus        interface{}
	AvailabilityZone         string
	Status                   string
	Description              interface{}
	Size                     int
	Attachments              []interface{}
	Bootable                 string
	Multiattach              bool
	VolumeType               interface{}
	OsVolHostAttrHost        interface{} //backend ceph
	OsVolTenantAttrTenantID  string
	OsVolMigStatusAttrNameID interface{}
}
func ListVolumes(client *v3.VolumeClient) ([]Volumes, error) {
	api, err := declareListVolumes(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListVolumesResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	volumes := []Volumes{}
	for _, v := range result.Volumes {
		volumes = append(volumes, Volumes{
			ID: v.ID,
			Name: v.Name,
			Size: v.Size,
			ReplicationStatus: v.ReplicationStatus,
			AvailabilityZone: v.AvailabilityZone,
			Status: v.Status,
			Description: v.Description,
			Attachments: v.Attachments,
			Bootable: v.Bootable,
			Multiattach: v.Multiattach,
			VolumeType: v.VolumeType,
			OsVolHostAttrHost: v.OsVolHostAttrHost, // backend of the volumes,
			OsVolTenantAttrTenantID: v.OsVolTenantAttrTenantID,
			OsVolMigStatusAttrNameID: v.OsVolMigStatusAttrNameID, //None means that a migration is not currently in progress

		})
	}
	return volumes, err
}

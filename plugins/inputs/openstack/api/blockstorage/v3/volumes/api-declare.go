package volumes

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
)

type ListVolumesRequest struct {
}

type ListVolumesResponse struct {
	Volumes []struct {
		Attachments        []interface{} `json:"attachments"`
		AvailabilityZone   string        `json:"availability_zone"`
		Bootable           string        `json:"bootable"`
		ConsistencygroupID interface{}   `json:"consistencygroup_id"`
		CreatedAt          string        `json:"created_at"`
		Description        interface{}   `json:"description"`
		Encrypted          bool          `json:"encrypted"`
		ID                 string        `json:"id"`
		Links              []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
		Metadata struct {
		} `json:"metadata"`
		MigrationStatus           interface{} `json:"migration_status"`
		Multiattach               bool        `json:"multiattach"`
		Name                      interface{} `json:"name"`
		OsVolHostAttrHost         interface{} `json:"os-vol-host-attr:host"`
		OsVolMigStatusAttrMigstat interface{} `json:"os-vol-mig-status-attr:migstat"`
		OsVolMigStatusAttrNameID  interface{} `json:"os-vol-mig-status-attr:name_id"`
		OsVolTenantAttrTenantID   string      `json:"os-vol-tenant-attr:tenant_id"`
		ReplicationStatus         interface{} `json:"replication_status"`
		Size                      int         `json:"size"`
		SnapshotID                interface{} `json:"snapshot_id"`
		SourceVolid               interface{} `json:"source_volid"`
		Status                    string      `json:"status"`
		UpdatedAt                 interface{} `json:"updated_at"`
		UserID                    string      `json:"user_id"`
		VolumeType                interface{} `json:"volume_type"`
	} `json:"volumes"`
}

//
func declareListVolumes(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListVolumesRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/volumes/detail",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
		RequestParameter: map[string]string{
			"all_tenants" : "1",
		},
		RequestParameterRequire: true,
	}, err
}

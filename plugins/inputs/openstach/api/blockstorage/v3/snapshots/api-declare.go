package snapshots

import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)
type ListSnapshotsRequest struct {

}
type ListSnapshotsResponse struct {
	Snapshots []struct {
		CreatedAt   string `json:"created_at"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Metadata    struct {
			Key string `json:"key"`
		} `json:"metadata"`
		Name                                  string      `json:"name"`
		OsExtendedSnapshotAttributesProgress  string      `json:"os-extended-snapshot-attributes:progress"`
		OsExtendedSnapshotAttributesProjectID string      `json:"os-extended-snapshot-attributes:project_id"`
		Size                                  int         `json:"size"`
		Status                                string      `json:"status"`
		UpdatedAt                             interface{} `json:"updated_at"`
		VolumeID                              string      `json:"volume_id"`
	} `json:"snapshots"`
}

//
func declareListSnapshots(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListSnapshotsRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/snapshots/detail",
		RequestHeader: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token": token,
		},
		RequestBody: jsonBody,
		RequestParameter: map[string]string{
			"all_tenants": "1",
		},
		RequestParameterRequire: true,
	}, err
}

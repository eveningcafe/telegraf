package snapshots

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3"
)

type Snapshot struct {
	ID          string
	Name        string
	Size        int
	Description string
	Progress    string
	ProjectID   string
	VolumeID    string
	Status      string
}

func List(client *v3.VolumeClient) ([]Snapshot, error) {
	api, err := declareListSnapshots(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListSnapshotsResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	snapshots := []Snapshot{}
	for _, v := range result.Snapshots{
		snapshots = append(snapshots, Snapshot{
			ID: v.ID,
			Name: v.Name,
			Size: v.Size,
			Description: v.Description,
			Progress: v.OsExtendedSnapshotAttributesProgress,
			ProjectID: v.OsExtendedSnapshotAttributesProjectID,
			VolumeID: v.VolumeID,
			Status: v.Status,
		})
	}
	return snapshots, err
}

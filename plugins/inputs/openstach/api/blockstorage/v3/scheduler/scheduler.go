package scheduler

import (
	"encoding/json"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3"
)

// StoragePool represents an individual StoragePool retrieved from the
// schedulerstats API.
type StoragePool struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
}


type Capabilities struct {
	TotalCapacityGb    float64    `json:"total_capacity_gb"`
	FreeCapacityGb     float64    `json:"free_capacity_gb"`
	VolumeBackendName  string `json:"volume_backend_name"`
	//The percentage of the total capacity that is reserved for the internal use by the back end.
	ReservedPercentage int    `json:"reserved_percentage"`
	DriverVersion      string `json:"driver_version"`
	MaxOverSubscriptionRatio string `json:"max_over_subscription_ratio"`

	//The storage back end for the back-end volume. For example, iSCSI or FC.
	StorageProtocol    string `json:"storage_protocol"`
	//The quality of service (QoS) support.
	QoSSupport         bool   `json:"QoS_support"`
	BackendState             string      `json:"backend_state"`
	ReplicationEnabled       bool        `json:"replication_enabled"`
	ProvisionedCapacityGb    int         `json:"provisioned_capacity_gb"`
	AllocatedCapacityGb      int         `json:"allocated_capacity_gb"`

}

func ListPool(client *v3.VolumeClient) ([]StoragePool, error) {
	api, err := declareListPool(client.Endpoint, client.Token)
	err = api.DoReuest()
	result := ListPoolResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	pools := []StoragePool{}
	for _, v := range result.Pools {
		pools = append(pools, StoragePool{
			Name:         v.Name,
			Capabilities: v.Capabilities,
		})
	}
	return pools, err
}

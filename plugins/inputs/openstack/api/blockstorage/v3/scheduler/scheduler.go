package scheduler

import (
	"encoding/json"
	"fmt"
	v3 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3"
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
	ReservedPercentage float64    `json:"reserved_percentage"`
	DriverVersion      string `json:"driver_version"`
	MaxOverSubscriptionRatio interface{} `json:"max_over_subscription_ratio"` // from version queens, it return float, or it return
	//The storage back end for the back-end volume. For example, iSCSI or FC.
	StorageProtocol    string `json:"storage_protocol"`
	//The quality of service (QoS) support.
	QoSSupport         bool   `json:"QoS_support"`
	ReplicationEnabled       bool        `json:"replication_enabled"`
	ProvisionedCapacityGb    float64         `json:"provisioned_capacity_gb"`
	AllocatedCapacityGb      float64         `json:"allocated_capacity_gb"`

}
//type Capabilities struct {
//			FilterFunction           interface{} `json:"filter_function"`
//			VendorName               string      `json:"vendor_name"`
//			GoodnessFunction         interface{} `json:"goodness_function"`
//			Multiattach              bool        `json:"multiattach"`
//			ProvisionedCapacityGb    float64     `json:"provisioned_capacity_gb"`
//			Timestamp                string      `json:"timestamp"`
//			AllocatedCapacityGb      int         `json:"allocated_capacity_gb"`
//			VolumeBackendName        string      `json:"volume_backend_name"`
//			ThinProvisioningSupport  bool        `json:"thin_provisioning_support"`
//			FreeCapacityGb           float64     `json:"free_capacity_gb"`
//			DriverVersion            string      `json:"driver_version"`
//			TotalCapacityGb          float64     `json:"total_capacity_gb"`
//			ReservedPercentage       int         `json:"reserved_percentage"`
//			MaxOverSubscriptionRatio float64     `json:"max_over_subscription_ratio"`
//			ReplicationEnabled       bool        `json:"replication_enabled"`
//			StorageProtocol          string      `json:"storage_protocol"`
//}


func ListPool(client *v3.VolumeClient) ([]StoragePool, error) {
	api, err := declareListPool(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err!= nil{
		return nil,err
	}
	result := ListPoolResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	if err!= nil{
		return nil,fmt.Errorf("can't prase json format of respone to %s %s :%s",api.Endpoint, api.Path,err)
	}
	pools := []StoragePool{}
	for _, v := range result.Pools {
		pools = append(pools, StoragePool{
			Name:         v.Name,
			Capabilities: v.Capabilities,
		})
	}
	return pools, err
}

package hypervisors

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2"
)
// Hypervisor represents a hypervisor in the OpenStack cloud.

type Hypervisor struct {
	CPUInfo struct {
		Arch     string   `json:"arch"`
		Model    string   `json:"model"`
		Vendor   string   `json:"vendor"`
		Features []string `json:"features"`
		Topology struct {
			Cores   int `json:"cores"`
			Threads int `json:"threads"`
			Sockets int `json:"sockets"`
		} `json:"topology"`
	} `json:"cpu_info"`
	CurrentWorkload    float64    `json:"current_workload"`
	Status             string `json:"status"`
	State              string `json:"state"`
	DiskAvailableLeast float64    `json:"disk_available_least"`
	HostIP             string `json:"host_ip"`
	FreeDiskGb         float64    `json:"free_disk_gb"`
	FreeRAMMb          float64    `json:"free_ram_mb"`
	HypervisorHostname string `json:"hypervisor_hostname"`
	HypervisorType     string `json:"hypervisor_type"`
	HypervisorVersion  float64    `json:"hypervisor_version"`
	ID                 string `json:"id"`
	LocalGb            float64    `json:"local_gb"`
	LocalGbUsed        float64    `json:"local_gb_used"`
	MemoryMb           float64    `json:"memory_mb"`
	MemoryMbUsed       float64    `json:"memory_mb_used"`
	RunningVms         int    `json:"running_vms"`
	Service            struct {
		Host           string      `json:"host"`
		ID             string      `json:"id"`
		DisabledReason interface{} `json:"disabled_reason"`
	} `json:"service"`
	Vcpus     float64 `json:"vcpus"`
	VcpusUsed float64 `json:"vcpus_used"`
}

func List(client *v2.ComputeClient) ([]Hypervisor, error) {
	api, err := declareListHypervisor(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	result := ListHypervisorResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	hypervisors := []Hypervisor{}
	for _, v := range result.Hypervisors {
		hypervisors = append(hypervisors, v)}
	return hypervisors, err
}


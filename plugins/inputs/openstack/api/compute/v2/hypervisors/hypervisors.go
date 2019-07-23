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
	CurrentWorkload    int    `json:"current_workload"`
	Status             string `json:"status"`
	State              string `json:"state"`
	DiskAvailableLeast int    `json:"disk_available_least"`
	HostIP             string `json:"host_ip"`
	FreeDiskGb         int    `json:"free_disk_gb"`
	FreeRAMMb          int    `json:"free_ram_mb"`
	HypervisorHostname string `json:"hypervisor_hostname"`
	HypervisorType     string `json:"hypervisor_type"`
	HypervisorVersion  int    `json:"hypervisor_version"`
	ID                 string `json:"id"`
	LocalGb            int    `json:"local_gb"`
	LocalGbUsed        int    `json:"local_gb_used"`
	MemoryMb           int    `json:"memory_mb"`
	MemoryMbUsed       int    `json:"memory_mb_used"`
	RunningVms         int    `json:"running_vms"`
	Service            struct {
		Host           string      `json:"host"`
		ID             string      `json:"id"`
		DisabledReason interface{} `json:"disabled_reason"`
	} `json:"service"`
	Vcpus     int `json:"vcpus"`
	VcpusUsed int `json:"vcpus_used"`
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


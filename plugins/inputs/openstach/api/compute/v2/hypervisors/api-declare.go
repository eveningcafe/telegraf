package hypervisors
import (
	"encoding/json"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base/request"
)

type ListHypervisorRequest struct {
}


type ListHypervisorResponse struct {
	Hypervisors []struct {
		CPUInfo CPUInfo
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
		ID                 string    `json:"id"`
		LocalGb            int    `json:"local_gb"`
		LocalGbUsed        int    `json:"local_gb_used"`
		MemoryMb           int    `json:"memory_mb"`
		MemoryMbUsed       int    `json:"memory_mb_used"`
		RunningVms         int    `json:"running_vms"`
		Service            struct {
			Host           string      `json:"host"`
			ID             string         `json:"id"`
			DisabledReason interface{} `json:"disabled_reason"`
		} `json:"service"`
		Vcpus     int `json:"vcpus"`
		VcpusUsed int `json:"vcpus_used"`
	} `json:"hypervisors"`
	HypervisorsLinks []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"hypervisors_links"`
}
//
func declareListHypervisor(endpoint string, token string) (*request.OpenstackAPI, error) {
	req := ListHypervisorRequest{}
	jsonBody, err := json.Marshal(req)
	return &request.OpenstackAPI{
		Method:   "GET",
		Endpoint: endpoint,
		Path:     "/os-hypervisors/detail",
		HeaderRequest: map[string]string{
			"Content-Type": "application/json",
			"X-OpenStack-Nova-API-Version": request.XOpenStackNovaAPIv2Version,
			"X-Auth-Token": token,
		},
		Request: jsonBody,
	}, err
}

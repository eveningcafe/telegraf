package v2

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
)

type ComputeClient struct{
	*base.Client
}
func NewComputeV2(providerClient authenticator.ProviderClient, region string) (*ComputeClient, error) {
	var err error
	c := &ComputeClient{}
	c.Client, err  = base.NewClient( providerClient , region , "compute")
	return c, err
}

// CPUInfo represents CPU information of the hypervisor.
type CPUInfo struct {
	Vendor   string   `json:"vendor"`
	Arch     string   `json:"arch"`
	Model    string   `json:"model"`
	Features []string `json:"features"`
	Topology Topology `json:"topology"`
}

type Topology struct {
	Sockets int `json:"sockets"`
	Cores   int `json:"cores"`
	Threads int `json:"threads"`
}

// Hypervisor represents a hypervisor in the OpenStack cloud.
type Hypervisor struct {
	// A structure that contains cpu information like arch, model, vendor,
	// features and topology.
	CPUInfo CPUInfo `json:"-"`

	// The current_workload is the number of tasks the hypervisor is responsible
	// for. This will be equal or greater than the number of active VMs on the
	// system (it can be greater when VMs are being deleted and the hypervisor is
	// still cleaning up).
	CurrentWorkload int `json:"current_workload"`

	// Status of the hypervisor, either "enabled" or "disabled".
	Status string `json:"status"`

	// State of the hypervisor, either "up" or "down".
	State string `json:"state"`

	// DiskAvailableLeast is the actual free disk on this hypervisor,
	// measured in GB.
	DiskAvailableLeast int `json:"disk_available_least"`

	// HostIP is the hypervisor's IP address.
	HostIP string `json:"host_ip"`

	// FreeDiskGB is the free disk remaining on the hypervisor, measured in GB.
	FreeDiskGB int `json:"-"`

	// FreeRAMMB is the free RAM in the hypervisor, measured in MB.
	FreeRamMB int `json:"free_ram_mb"`

	// HypervisorHostname is the hostname of the hypervisor.
	HypervisorHostname string `json:"hypervisor_hostname"`

	// HypervisorType is the type of hypervisor.
	HypervisorType string `json:"hypervisor_type"`

	// HypervisorVersion is the version of the hypervisor.
	HypervisorVersion int `json:"-"`

	// ID is the unique ID of the hypervisor.
	ID string `json:"-"`

	// LocalGB is the disk space in the hypervisor, measured in GB.
	LocalGB int `json:"-"`

	// LocalGBUsed is the used disk space of the  hypervisor, measured in GB.
	LocalGBUsed int `json:"local_gb_used"`

	// MemoryMB is the total memory of the hypervisor, measured in MB.
	MemoryMB int `json:"memory_mb"`

	// MemoryMBUsed is the used memory of the hypervisor, measured in MB.
	MemoryMBUsed int `json:"memory_mb_used"`

	// RunningVMs is the The number of running vms on the hypervisor.
	RunningVMs int `json:"running_vms"`

	// Service is the service this hypervisor represents.
	//Service Service `json:"service"`

	// VCPUs is the total number of vcpus on the hypervisor.
	VCPUs int `json:"vcpus"`

	// VCPUsUsed is the number of used vcpus on the hypervisor.
	VCPUsUsed int `json:"vcpus_used"`
}




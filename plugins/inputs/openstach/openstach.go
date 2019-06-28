package openstach

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/endpoints"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/services"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/projects"
	"log"
)

const (
	// plugin is used to identify ourselves in log output
	plugin = "openstack"
)

var sampleConfig = `
  ## This is the recommended interval to poll.
  interval = '60m'
  ## The identity endpoint to authenticate against and get the
  ## service catalog from
  identity_endpoint = "https://my.openstack.cloud:5000"
  ## The domain to authenticate against when using a V3
  ## identity endpoint.  Defaults to 'default'
  domain = "default"
  ## The project to authenticate as
  project = "admin"
  ## The user to authenticate as, must have admin rights
  username = "admin"
  ## The user's password to authenticate with
  password = "Passw0rd"
`

// OpenStack is the main structure associated with a collection instance.
type OpenStack struct {
	// Configuration variables
	IdentityEndpoint string
	Domain           string
	Project          string
	Username         string
	Password         string

	// Locally cached clients
	identity *gophercloud.ServiceClient
	compute  *gophercloud.ServiceClient
	volume   *gophercloud.ServiceClient

	// Locally cached resources
	services     serviceMap
	projects     projectMap
	hypervisors  hypervisorMap
	flavors      flavorMap
	servers      serverMap
	volumes      volumeMap
	storagePools storagePoolMap
}

// SampleConfig return a sample configuration file for auto-generation and
// implements the Input interface.
func (o *OpenStack) SampleConfig() string {
	return sampleConfig
}

// Description returns a description string of the input plugin and implements
// the Input interface.
func (o *OpenStack) Description() string {
	return "Collects performance metrics from OpenStack services"
}

// initialize performs any necessary initialization functions
func (o *OpenStack) initialize() error {
	// Authenticate against Keystone and get a token provider
	provider, err := openstack.AuthenticatedClient(gophercloud.AuthOptions{
		IdentityEndpoint: o.IdentityEndpoint,
		DomainName:       o.Domain,
		TenantName:       o.Project,
		Username:         o.Username,
		Password:         o.Password,
	})
	if err != nil {
		return fmt.Errorf("Unable to authenticate OpenStack user: %v", err)
	}

	// Create required clients and attach to the OpenStack struct
	if o.identity, err = openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{}); err != nil {
		return fmt.Errorf("unable to create V3 identity client: %v", err)
	}
	if o.compute, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{}); err != nil {
		return fmt.Errorf("unable to create V2 compute client: %v", err)
	}
	if o.volume, err = openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{}); err != nil {
		return fmt.Errorf("unable to create V2 block storage client: %v", err)
	}

	// Initialize resource maps and slices
	o.services = serviceMap{}
	o.projects = projectMap{}
	o.hypervisors = hypervisorMap{}
	o.flavors = flavorMap{}
	o.servers = serverMap{}
	o.volumes = volumeMap{}
	o.storagePools = storagePoolMap{}

	return nil
}
// gatherHypervisors collects hypervisors from the OpenStack API.
func (o *OpenStack) gatherHypervisors() error {
	page, err := hypervisors.List(o.compute).AllPages()
	if err != nil {
		return fmt.Errorf("unable to list hypervisors: %v", err)
	}
	hypervisors, err := hypervisors.ExtractHypervisors(page)
	if err != nil {
		return fmt.Errorf("unable to extract hypervisors: %v", err)
	}
	for _, hypervisor := range hypervisors {
		o.hypervisors[hypervisor.ID] = hypervisor
	}
	return nil
}

// accumulateHypervisors accumulates statistics from hypervisors.
func (o *OpenStack) accumulateCompute(acc telegraf.Accumulator) {
	for _, hypervisor := range o.hypervisors {
		tags := tagMap{
			"name": hypervisor.HypervisorHostname,
		}
		fields := fieldMap{
			"memory_mb":      hypervisor.MemoryMB,
			"memory_mb_used": hypervisor.MemoryMBUsed,
			"running_vms":    hypervisor.RunningVMs,
			"vcpus":          hypervisor.VCPUs,
			"vcpus_used":     hypervisor.VCPUsUsed,
		}
		acc.AddFields("openstack_compute", fields, tags)
	}
}

func (o *OpenStack) Gather(acc telegraf.Accumulator) error {
	// Perform any required set up
	if err := o.initialize(); err != nil {
		return err
	}

	// Gather resources.  Note service harvesting must come first as the other
	// gatherers are dependant on this information.
	gatherers := map[string]func() error{
		"services":      o.gatherServices,
		"projects":      o.gatherProjects,
		"hypervisors":   o.gatherHypervisors,
		"flavors":       o.gatherFlavors,
		"servers":       o.gatherServers,
		"volumes":       o.gatherVolumes,
		"storage pools": o.gatherStoragePools,
	}
	for resources, gatherer := range gatherers {
		if err := gather(gatherer); err != nil {
			log.Println("W!", plugin, "failed to get", resources, ":", err)
		}
	}

	// Accumulate statistics
	accumulators := []func(telegraf.Accumulator){
		o.accumulateIdentity,
		o.accumulateHypervisors,
		o.accumulateServers,
		o.accumulateVolumes,
		o.accumulateStoragePools,
	}
	for _, accumulator := range accumulators {
		accumulator(acc)
	}

	return nil
}
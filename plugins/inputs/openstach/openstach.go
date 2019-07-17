package openstach

import (
	"fmt"
	"log"
	"strconv"

	//"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	blockstorage "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3"
	blockstorageScheduler "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3/scheduler"
	blockstorageServices "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3/services"
	blockstorageSnapshots "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3/snapshots"
	blockstorageVolumes "github.com/influxdata/telegraf/plugins/inputs/openstach/api/blockstorage/v3/volumes"
	compute "github.com/influxdata/telegraf/plugins/inputs/openstach/api/compute/v2"
	computeHypervisors "github.com/influxdata/telegraf/plugins/inputs/openstach/api/compute/v2/hypervisors"
	computeServices "github.com/influxdata/telegraf/plugins/inputs/openstach/api/compute/v2/services"
	identity "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/groups"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/projects"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/services"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/users"
	networking "github.com/influxdata/telegraf/plugins/inputs/openstach/api/networking/v2"
	networkingAgent "github.com/influxdata/telegraf/plugins/inputs/openstach/api/networking/v2/agents"
)

const (
	// plugin is used to identify ourselves in log output
	plugin = "openstack"
)

var sampleConfig = `
  ## This is the recommended interval to poll.
  interval = '60m'
  ## The identity endpoint to authenticate against openstack service, use v3 indentity api
  identity_endpoint = "https://my.openstack.cloud:5000"
  ## The domain project to authenticate 
  ## identity endpoint.  Defaults to 'default'
  project_domain_id = "default"
  ## The project of user to authenticate
  project = "admin"
  ## The domain user to authenticate
  user_domain_id = "default"
  ## The user to authenticate as, must have admin rights,or monitor rights
  username = "admin"
  ## The user's password to authenticate with
  password = "Passw0rd"
  ## Opesntack service type collector
  services_gather = [
    "identity"
    "volumev3"
    "network"
    "compute"
  ]
  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
`

type tagMap map[string]string
type fieldMap map[string]interface{}
type serviceMap map[string]services.Service

// projectMap maps a project id to a Project struct.
type projectMap map[string]projects.Project
type userMap map[string]users.User
type groupMap map[string]groups.Group
type hypervisorMap map[string]computeHypervisors.Hypervisor

//volumeMap maps a volume id to a volume struct.
type volumeMap map[string]blockstorageVolumes.Volumes

// storagePoolMap maps a storage pool name to a StoragePool struct.
type storagePoolMap map[string]blockstorageScheduler.StoragePool

type networkMap map[string]blockstorageScheduler.StoragePool

// OpenStack is the main structure associated with a collection instance.
type OpenStack struct {
	// Configuration variables
	IdentityEndpoint string
	ProjectDomainID  string
	UserDomainID     string
	Project          string
	Username         string
	Password         string
	Services_gather  []string
	Region           string

	// Locally cached clients
	identityClient *identity.IdentityClient
	computeClient  *compute.ComputeClient
	volumeClient   *blockstorage.VolumeClient
	networkClient  *networking.NetworkClient

	// Locally cached resources
	services     serviceMap
	users        userMap
	groups       groupMap
	projects     projectMap
	hypervisors  hypervisorMap
	storagePools storagePoolMap
	volumes      volumeMap
	networks     networkMap
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
	provider, err := authenticator.AuthenticatedClient(authenticator.AuthOption{
		AuthURL:         o.IdentityEndpoint,
		ProjectDomainId: o.ProjectDomainID,
		UserDomainId:    o.UserDomainID,
		Username:        o.Username,
		Password:        o.Password,
		Project_name:    o.Project,
	})
	if err != nil {
		return fmt.Errorf("Unable to authenticate OpenStack user: %v", err)
	}

	// Create required clients and attach to the OpenStack struct
	if o.identityClient, err = identity.NewIdentityV3(*provider, o.Region); err != nil {
		return fmt.Errorf("unable to create V3 identity client: %v", err)
	}
	if o.computeClient, err = compute.NewComputeV2(*provider, o.Region); err != nil {
		return fmt.Errorf("unable to create V2 compute client: %v", err)
	}
	if o.volumeClient, err = blockstorage.NewBlockStorageV3(*provider, o.Region); err != nil {
		return fmt.Errorf("unable to create V3 block storage client: %v", err)
	}

	// Initialize resource maps and slices
	o.services = serviceMap{}
	o.projects = projectMap{}
	o.users = userMap{}
	o.groups = groupMap{}
	o.hypervisors = hypervisorMap{}
	o.storagePools = storagePoolMap{}
	o.volumes = volumeMap{}

	return err
}

// gatherHypervisors collects hypervisors from the OpenStack API.
func (o *OpenStack) gatherHypervisors() error {
	hypervisors, err := computeHypervisors.List(o.computeClient)
	if err != nil {
		return fmt.Errorf("unable to extract hypervisors: %v", err)
	}
	for _, hypervisor := range hypervisors {
		o.hypervisors[hypervisor.ID] = hypervisor
	}
	//return err
	return nil
}

// gatherServices collects services from the OpenStack API.
func (o *OpenStack) gatherServices() error {
	services, err := services.List(o.identityClient)
	if err != nil {
		return fmt.Errorf("unable to list services: %v", err)
	}
	for _, service := range services {
		o.services[service.ID] = service
	}

	return err
}

func (o *OpenStack) gatherProjects() error {
	projects, err := projects.List(o.identityClient)
	if err != nil {
		return fmt.Errorf("unable to list project: %v", err)
	}
	for _, project := range projects {
		o.projects[project.ID] = project
	}

	return err
}
func (o *OpenStack) gatherUsers() error {
	users, err := users.List(o.identityClient)
	if err != nil {
		return fmt.Errorf("unable to list user: %v", err)
	}
	for _, user := range users {
		o.users[user.ID] = user
	}

	return err
}

func (o *OpenStack) gatherGroups() error {
	groups, err := groups.List(o.identityClient)
	if err != nil {
		return fmt.Errorf("unable to list groups: %v", err)
	}
	for _, group := range groups {
		o.groups[group.ID] = group
	}

	return err
}

// gather volume per project
func (o *OpenStack) gatherVolumes() error {
	volumes, err := blockstorageVolumes.ListVolumes(o.volumeClient)
	if err != nil {
		return fmt.Errorf("unable to list volumes: %v", err)
	}
	for _, volume := range volumes {
		o.volumes[volume.ID] = volume
	}
	return err
}

// gather storage per pool
func (o *OpenStack) gatherStoragePools() error {
	storagePools, err := blockstorageScheduler.ListPool(o.volumeClient)
	if err != nil {
		return fmt.Errorf("unable to list storage pools: %v", err)
	}
	for _, storagePool := range storagePools {
		o.storagePools[storagePool.Name] = storagePool
	}
	return err
}

//
func (o *OpenStack) accumulateComputeAgents(acc telegraf.Accumulator) {
	agents, err := computeServices.List(o.computeClient)
	if err != nil {

	} else {
		fields := fieldMap{}
		for _, agent := range agents {
			if agent.State == "up" {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			acc.AddFields("openstack_compute", fields, tagMap{
				"region":   o.Region,
				"service":  agent.Binary,
				"hostname": agent.Host,
				"status":   agent.Status,
				"zone":     agent.Zone,
			})
		}
	}
}

// accumulateHypervisors accumulates statistics from hypervisors.
func (o *OpenStack) accumulateComputeHypervisors(acc telegraf.Accumulator) {

	for _, hypervisor := range o.hypervisors {
		fields := fieldMap{
			"memory_mb_total":     hypervisor.MemoryMB,
			"memory_mb_used":      hypervisor.MemoryMBUsed,
			"running_vms":         hypervisor.RunningVMs,
			"vcpus_total":         hypervisor.VCPUs,
			"vcpus_used":          hypervisor.VCPUsUsed,
			"local_disk_avalable": hypervisor.LocalGB,
			"local_disk_usage":    hypervisor.LocalGBUsed,
		}
		acc.AddFields("openstack_compute", fields, tagMap{
			"name":   hypervisor.HypervisorHostname,
			"region": o.Region,
		})
	}
}

// accumulateIdentity accumulates statistics from the identity service.
func (o *OpenStack) accumulateIdentity(acc telegraf.Accumulator) {

	fields := fieldMap{
		"num_projects": len(o.projects),
		"num_servives": len(o.services),
		"num_users":    len(o.users),
		"num_group":    len(o.groups),
	}
	acc.AddFields("openstack_identity", fields, tagMap{})
}

//
func (o *OpenStack) accumulateNetworkAgents(acc telegraf.Accumulator) {
	agents, err := networkingAgent.List(o.networkClient)
	if err != nil {
	} else {
		fields := fieldMap{}
		for _, agent := range agents {
			if agent.Alive == true {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			var status string
			if agent.AdminStateUp == true {
				status = "enable"
			}else {
				status = "disable"
			}

			acc.AddFields("openstack_volumes", fields, tagMap{
				"region":   o.Region,
				"service":  agent.Binary,
				"hostname": agent.Host,
				"status":   status,
				"zone":     agent.AvailabilityZone,
			})
		}
	}
}

//
func (o *OpenStack) accumulateVolumeAgents(acc telegraf.Accumulator) {
	agents, err := blockstorageServices.List(o.volumeClient)
	if err != nil {
	} else {
		fields := fieldMap{}
		for _, agent := range agents {
			if agent.State == "up" {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			acc.AddFields("openstack_volumes", fields, tagMap{
				"region":   o.Region,
				"service":  agent.Binary,
				"hostname": agent.Host,
				"status":   agent.Status,
				"zone":     agent.Zone,
			})
		}
	}
}

func (o *OpenStack) accumulateVolumesPerTenant(acc telegraf.Accumulator) {
	var numSnapshots int
	snapshots, _ := blockstorageSnapshots.List(o.volumeClient)

	// add metric volumes
	for _, volume := range o.volumes {
		var volumeType interface{}
		if volume.VolumeType != nil {
			volumeType = volume.VolumeType
		} else {
			volumeType = "unknown"
		}
		// Try derive the associated project
		var volumeProject interface{}
		if p, ok := o.projects[volume.OsVolTenantAttrTenantID]; ok {
			volumeProject = p.Name
		} else {
			volumeProject = "unknown"
		}

		tags := tagMap{
			"volumeName": volume.Name.(string),
			"volumeID": volume.ID,
			"project":    volumeProject.(string),
			"volumeType":       volumeType.(string),
			"status":     volume.Status,
			"zone":       volume.AvailabilityZone,
		}

		fields := fieldMap{
			"size_gb": volume.Size,
		}
		acc.AddFields("openstack_volume", fields, tags)
		// if volumes has snapshot add metric of it
		numSnapshots = 0
		for _, snapshot := range snapshots {
			if snapshot.VolumeID == volume.ID {
				numSnapshots++
			}
		}
		if numSnapshots != 0 {
			acc.AddFields("openstack_volume", fieldMap{
				"num_snapshots": numSnapshots,
			}, tags)
		}
	}
}

// accumulateStoragePools accumulates statistics about storage pools.
func (o *OpenStack) accumulateVolumeStoragePools(acc telegraf.Accumulator) {

	for _, storagePool := range o.storagePools {
		tags := tagMap{
			"backed_state": storagePool.Capabilities.BackendState,
			"backend_name": storagePool.Capabilities.VolumeBackendName,
			"region":       o.Region,
		}
		overcommit, _ := strconv.ParseFloat(storagePool.Capabilities.MaxOverSubscriptionRatio, 64)
		fields := fieldMap{
			"total_capacity_gb":           storagePool.Capabilities.TotalCapacityGb,
			"free_capacity_gb":            storagePool.Capabilities.FreeCapacityGb,
			"allocated_capacity_gb":       storagePool.Capabilities.AllocatedCapacityGb,
			"provisioned_capacity_gb":     storagePool.Capabilities.ProvisionedCapacityGb,
			"max_over_subscription_ratio": overcommit,
		}
		acc.AddFields("openstack_storage_pool", fields, tags)
	}
}

// gather is a wrapper around library calls out to gophercloud that catches
// and recovers from panics.  Evidently if things like volumes don't exist
// then it will go down in flames.
func gather(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from crash: %v", r)
		}
	}()
	return f()
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
		"users":         o.gatherUsers,
		"group":         o.gatherGroups,
		"hypervisors":   o.gatherHypervisors,
		"storage pools": o.gatherStoragePools,
		"volumes":       o.gatherVolumes,
	}

	for resources, gatherer := range gatherers {
		if err := gather(gatherer); err != nil {
			log.Println("W!", plugin, "failed to get", resources, " : ", err)
		}
	}
	// Accumulate statistics
	accumulators := []func(telegraf.Accumulator){
		o.accumulateIdentity,
		o.accumulateComputeAgents,
		o.accumulateComputeHypervisors,
		o.accumulateNetworkAgents,
		o.accumulateVolumeAgents,
		o.accumulateVolumeStoragePools,
		o.accumulateVolumesPerTenant,
	}
	for _, accumulator := range accumulators {
		//go routine in here
		accumulator(acc)
	}

	return nil
}

// init registers a callback which creates a new OpenStack input instance.
func init() {
	inputs.Add("openstack", func() telegraf.Input {
		return &OpenStack{
			UserDomainID:    "default",
			ProjectDomainID: "default",
		}
	})
}

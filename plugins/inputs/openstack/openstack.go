package openstack

import (
	"fmt"
	"github.com/influxdata/telegraf/filter"
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/prometheus/common/log"
	"reflect"
	"strconv"
	"sync"

	//"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerstats"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	blockstorage "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3"
	blockstorageQuotas "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3/quotas"
	blockstorageScheduler "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3/scheduler"
	blockstorageServices "github.com/influxdata/telegraf/plugins/inputs/openstack/api/blockstorage/v3/services"
	compute "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2"
	computeHypervisors "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2/hypervisors"
	computeQuotas "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2/quotas"
	computeServices "github.com/influxdata/telegraf/plugins/inputs/openstack/api/compute/v2/services"
	identity "github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/authenticator"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/groups"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/projects"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/services"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/users"
	networking "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
	networkingAgent "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/agents"
	networkingNET "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/networks"
	networkingQuotas "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/quotas"
)

const (
	// plugin is used to identify ourselves in log output
	plugin = "openstack"
)

// Note: timeout for each request is 10s, change in api/base/client/NewClient( _time_out_of_http_request)

var sampleConfig = `
  ## This is the recommended interval to poll.
  #interval = '5m'
  ## Name of cloud cluster, 
  cloud =  'my_openstack'
  ## The identity endpoint to authenticate against openstack service, use v3 indentity api
  identity_endpoint = "https://my.openstack.cloud:5000"
  ## Region of openstack cluster which crawled by plugin
  region = "RegionOne"
  ## The domain project to authenticate
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
    "identity",
    "volumev3",
    "network",
    "compute",
  ]
  ## there is no way we can get overcommit cpu and ram so It must config as same as the config at nova.conf
  cpu_overcommit_ratio = 16.0
  mem_overcommit_ratio = 1.5
  ## Optional TLS Config
  tls_ca = "/etc/telegraf/openstack.crt"
  ## Optional Use TLS but skip chain & host verification
  insecure_skip_verify = false
`

type tagMap map[string]string
type fieldMap map[string]interface{}
type serviceMap map[string]services.Service

// projectMap maps a project id to a Project struct.
type projectMap map[string]projects.Project
type userMap map[string]users.User
type groupMap map[string]groups.Group
type networkMap map[string]networkingNET.Network

// OpenStack is the main structure associated with a collection instance.
type OpenStack struct {
	// Configuration variables
	IdentityEndpoint   string
	ProjectDomainID    string
	UserDomainID       string
	Project            string
	Username           string
	Password           string
	Region             string
	Cloud              string
	ServicesGather     []string
	CpuOvercommitRatio float64
	MemOvercommitRatio float64
	tls.ClientConfig

	//filter service
	filter filter.Filter
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
	accumulators []func(telegraf.Accumulator) error

	// use false when it first run
	initialized bool
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
	if o.initialized {
		return nil
	}
	tlsCfg, err := o.ClientConfig.TLSConfig()
	// Authenticate against Keystone and get a token provider
	provider, err := authenticator.AuthenticatedClient(authenticator.AuthOption{
		AuthURL:         o.IdentityEndpoint,
		ProjectDomainId: o.ProjectDomainID,
		UserDomainId:    o.UserDomainID,
		Username:        o.Username,
		Password:        o.Password,
		Project_name:    o.Project,
		TlsCfg:          tlsCfg,
	})

	if err != nil {
		return fmt.Errorf("Unable to authenticate OpenStack user: %v", err)
	}
	// Create the required clients and attach to OpenStack object, follow the configuration file, but it always have to created identity client

	if o.identityClient, err = identity.NewIdentityV3(*provider, o.Region); err != nil {
		return fmt.Errorf("unable to create V3 identity client: %v", err)
	}

	if o.filter.Match("compute") {
		if o.computeClient, err = compute.NewComputeV2(*provider, o.Region); err != nil {
			return fmt.Errorf("unable to create V2 compute client: %v", err)
		}
	}
	if o.filter.Match("volumev3") {
		if o.volumeClient, err = blockstorage.NewBlockStorageV3(*provider, o.Region); err != nil {
			return fmt.Errorf("unable to create V3 block storage client: %v", err)
		}
	}
	if o.filter.Match("network") {
		if o.networkClient, err = networking.NewNetworkClientV2(*provider, o.Region); err != nil {
			return fmt.Errorf("unable to create V3 networking client: %v", err)
		}
	}

	// Initialize resource maps and slices
	o.services = serviceMap{}
	o.projects = projectMap{}
	o.users = userMap{}
	o.groups = groupMap{}
	if err != nil {
		o.initialized = true
		return nil
	}
	return err
}

// gatherServices collects init properties for openStack
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
		// only get project in default, or configuration id
		if project.DomainID == o.ProjectDomainID {
			o.projects[project.ID] = project
		} else {
			continue
		}

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

//accumulateStoragePools accumulates statistics about nova agents.
func (o *OpenStack) accumulateComputeAgents(acc telegraf.Accumulator) error {
	agents, err := computeServices.List(o.computeClient)
	if err != nil {
		acc.AddFields("openstack_compute", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
			"cloud":  o.Cloud,
		})
	} else {
		acc.AddFields("openstack_compute", fieldMap{"api_state": 1,}, tagMap{
			"region": o.Region,
			"cloud":  o.Cloud,
		})

		fields := fieldMap{}
		for _, agent := range agents {
			if agent.State == "up" {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			acc.AddFields("openstack_compute", fields, tagMap{
				"cloud":      o.Cloud,
				"region":     o.Region,
				"service":    agent.Binary,
				"agent_host": agent.Host,
				"status":     agent.Status,
				"zone":       agent.Zone,
			})
		}
	}
	return err
}

// accumulateStoragePools accumulates statistics about nova hypervisors.
func (o *OpenStack) accumulateComputeHypervisors(acc telegraf.Accumulator) error {
	hypervisors, err := computeHypervisors.List(o.computeClient)
	if err != nil {
	} else {
		//var MemOverCommitRatio float64
		//var CpuOverCommitRatio float64
		for _, hypervisor := range hypervisors {
			//MemOverCommitRatio, err = strconv.ParseFloat(o.MemOvercomitRatio, 64)
			//CpuOverCommitRatio, err = strconv.ParseFloat(o.CpuOvercomitRatio, 64)
			fields := fieldMap{
				"memory_mb_total":      hypervisor.MemoryMb,
				"memory_mb_used":       hypervisor.MemoryMbUsed,
				"running_vms":          hypervisor.RunningVms,
				"cpus_total":           hypervisor.Vcpus,
				"cpus_used":            hypervisor.VcpusUsed,
				"cpu_overcommit_ratio": o.CpuOvercommitRatio,
				"mem_overcommit_ratio": o.MemOvercommitRatio,
				"local_disk_avalable":  hypervisor.LocalGb,
				"local_disk_usage":     hypervisor.LocalGbUsed,
			}
			acc.AddFields("openstack_compute", fields, tagMap{
				"hypervisor_host": hypervisor.HypervisorHostname,
				"cloud":           o.Cloud,
				"region":          o.Region,
			})
		}
	}
	return err
}
func (o *OpenStack) accumulateProjectQuotas(acc telegraf.Accumulator) {
	for _, p := range o.projects {
		if p.Name == "service" {
			continue
		} else {
			tags := tagMap{
				"cloud":   o.Cloud,
				"region":  o.Region,
				"project": p.Name,
			}
			fields := fieldMap{}
			getQuotas := []func(string) (fieldMap, error){}

			if o.filter.Match("compute") {
				getQuotas = append(getQuotas, o.getComputeProjectQuotas)
			}

			if o.filter.Match("volumev3") {
				getQuotas = append(getQuotas, o.getVolumeProjectQuotas)
			}
			if o.filter.Match("network") {
				getQuotas = append(getQuotas, o.getNetworkProjectQuotas)
			}

			wg := sync.WaitGroup{}
			for _, get := range getQuotas {
				wg.Add(1)
				//go routine in here
				go func(get func(string) (fieldMap, error)) {
					defer wg.Done()
					newField, err := get(p.ID)
					if err != nil {
						acc.AddError(err)
					}
					for k, v := range newField {
						fields[k] = v
					}
				}(get)
			}
			wg.Wait()
			acc.AddFields("openstack_quotas", fields, tags)
		}
	}
}

//accumulateStoragePools accumulates statistics about nova project quotas.
func (o *OpenStack) getComputeProjectQuotas(projectsID string) (fieldMap, error) {
	var err error
	var fields fieldMap
	computeQuotas, err := computeQuotas.Detail(o.computeClient, projectsID)

	if err != nil {

	} else {
		if (computeQuotas.Cores.Limit == -1) {
			computeQuotas.Cores.Limit = 9999
		}
		if (computeQuotas.RAM.Limit == -1) {
			computeQuotas.RAM.Limit = 9999
		}
		if (computeQuotas.Instances.Limit == -1) {
			computeQuotas.Instances.Limit = 9999
		}

		fields = fieldMap{
			"cpu_limit":      computeQuotas.Cores.Limit,
			"cpu_used":       computeQuotas.Cores.InUse,
			"ram_limit":      computeQuotas.RAM.Limit,
			"ram_used":       computeQuotas.RAM.InUse,
			"instance_limit": computeQuotas.Instances.Limit,
			"instance_used":  computeQuotas.Instances.InUse,
		}
	}

	return fields, err
}

// accumulateIdentity accumulates statistics from the identity service.
func (o *OpenStack) accumulateIdentity(acc telegraf.Accumulator) error {

	fields := fieldMap{
		"num_projects": len(o.projects),
		"num_servives": len(o.services),
		"num_users":    len(o.users),
		"num_group":    len(o.groups),
	}
	acc.AddFields("openstack_identity", fields, tagMap{
		"cloud":  o.Cloud,
		"region": o.Region,
	})

	return nil
}

//accumulateStoragePools accumulates statistics about neutron agent state.
func (o *OpenStack) accumulateNetworkAgents(acc telegraf.Accumulator) error {
	agents, err := networkingAgent.List(o.networkClient)
	if err != nil {
		acc.AddFields("openstack_network", fieldMap{"api_state": 0,}, tagMap{
			"cloud":  o.Cloud,
			"region": o.Region,
		})
	} else {
		acc.AddFields("openstack_network", fieldMap{"api_state": 1,}, tagMap{
			"cloud":  o.Cloud,
			"region": o.Region,
		})
		fields := fieldMap{}
		for _, agent := range agents {
			if agent.Alive == true {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			var status string
			if agent.AdminStateUp == true {
				status = "enabled"
			} else {
				status = "disabled"
			}

			acc.AddFields("openstack_network", fields, tagMap{
				"cloud":      o.Cloud,
				"region":     o.Region,
				"service":    agent.Binary,
				"agent_host": agent.Host,
				"status":     status,
				"zone":       agent.AvailabilityZone,
			})
		}
	}
	return err
}

//accumulateStoragePools accumulates statistics about neutron available IP pools.
func (o *OpenStack) accumulateNetworkIp(acc telegraf.Accumulator) error {
	// get provider network
	networks, err := networkingNET.List(o.networkClient)
	if err != nil {
		return err
	}
	providerNetworkMap := networkMap{}
	// filter only provider network
	for _, network := range networks {
		if network.ProviderPhysicalNetwork != "" {
			providerNetworkMap[network.ID] = network
		}
	}

	// get provider network pool
	ipAvail, err := networkingNET.NetworkIPAvailabilities(o.networkClient)
	if err != nil {
		err = fmt.Errorf("unable to get network info: %v", err)

	} else {
		var providerNetwork string
		for _, ipAvailNet := range ipAvail {

			if p, ok := providerNetworkMap[ipAvailNet.NetworkID]; ok {
				providerNetwork = p.Name
			} else {
				continue
			}

			// for each in loop of subnet
			for _, ipAvalSubnet := range ipAvailNet.SubnetIPAvailability {
				tags := tagMap{
					"cloud":            o.Cloud,
					"region":           o.Region,
					"network":          ipAvailNet.NetworkName,
					"subnet_cidr":      ipAvalSubnet.Cidr,
					"provider_network": providerNetwork,
				}
				// only collect ipv4
				if ipAvalSubnet.IPVersion == 4 {
					acc.AddFields("openstack_network", fieldMap{
						"ip_total": ipAvalSubnet.TotalIps.Int64(),
						"ip_used":  ipAvalSubnet.UsedIps.Int64(),
					}, tags)
				}
			}

			// overall outsite
			//acc.AddFields("openstack_network", fieldMap{
			//	"ip_used":  ipAvailNet.UsedIps.Int64(),
			//	"ip_total": ipAvailNet.TotalIps.Int64(),
			//}, tagMap{
			//	"cloud":            o.Cloud,
			//	"region":           o.Region,
			//	"network":          ipAvailNet.NetworkName,
			//	"subnet_cidr":      "all",
			//	"provider_network": providerNetwork,
			//})
		}
	}
	return err
}

//accumulateStoragePools accumulates statistics about neutron network quotas.
func (o *OpenStack) getNetworkProjectQuotas(projectID string) (fieldMap, error) {
	var err error
	var fields fieldMap

	netQuotas, err := networkingQuotas.Detail(o.networkClient, projectID)
	if err != nil {
	} else {
		// unlimit = -1
		if (netQuotas.Network.Limit == -1) {
			netQuotas.Network.Limit = 9999
		}
		if (netQuotas.SecurityGroup.Limit == -1) {
			netQuotas.SecurityGroup.Limit = 9999
		}
		if (netQuotas.Subnet.Limit == -1) {
			netQuotas.Subnet.Limit = 9999
		}
		if (netQuotas.Port.Limit == -1) {
			netQuotas.Port.Limit = 9999
		}
		// provider network model
		if (netQuotas.Floatingip.Limit <= 0) {
			netQuotas.Floatingip.Limit = 9999
			if netQuotas.Floatingip.Limit == 0 {
				netQuotas.Floatingip.Used = 0
			}

		}
		// provider network model limit = 0
		if (netQuotas.Router.Limit <= 0) {
			netQuotas.Router.Limit = 9999
			if netQuotas.Router.Limit == 0 {
				netQuotas.Router.Used = 0
			}
		}

		fields = fieldMap{
			"network_limit":       netQuotas.Network.Limit,
			"network_used":        netQuotas.Network.Used,
			"securityGroup_limit": netQuotas.SecurityGroup.Limit,
			"securityGroup_used":  netQuotas.SecurityGroup.Used,
			"securityRule_limit":  netQuotas.SecurityGroupRule.Limit,
			"securityRule_used":   netQuotas.SecurityGroupRule.Used,
			"subnet_limit":        netQuotas.Subnet.Limit,
			"subnet_used":         netQuotas.Subnet.Used,
			"port_limit":          netQuotas.Port.Limit,
			"port_used":           netQuotas.Port.Used,
			"floatingIP_limit":    netQuotas.Floatingip.Limit,
			"floatingIP_used":     netQuotas.Floatingip.Used,
		}
	}
	return fields, err
}

//accumulateStoragePools accumulates statistics about cinder volume agent.
func (o *OpenStack) accumulateVolumeAgents(acc telegraf.Accumulator) error {
	agents, err := blockstorageServices.List(o.volumeClient)
	if err != nil {
		acc.AddFields("openstack_volumes", fieldMap{"api_state": 0,}, tagMap{
			"cloud":  o.Cloud,
			"region": o.Region,
		})
	} else {
		acc.AddFields("openstack_volumes", fieldMap{"api_state": 1,}, tagMap{
			"cloud":  o.Cloud,
			"region": o.Region,
		})

		fields := fieldMap{}
		for _, agent := range agents {
			if agent.State == "up" {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			acc.AddFields("openstack_volumes", fields, tagMap{
				"cloud":      o.Cloud,
				"region":     o.Region,
				"service":    agent.Binary,
				"agent_host": agent.Host,
				"status":     agent.Status,
				"zone":       agent.Zone,
			})
		}
	}
	return err
}

// accumulateStoragePools accumulates statistics about storage pools.
func (o *OpenStack) accumulateVolumeStoragePools(acc telegraf.Accumulator) error {
	storagePools, err := blockstorageScheduler.ListPool(o.volumeClient)
	if err != nil {

	} else {
		for _, storagePool := range storagePools {
			tags := tagMap{
				"pool_name": storagePool.Name,
				"cloud":     o.Cloud,
				"region":    o.Region,
			}

			var overcommit float64
			if reflect.TypeOf(storagePool.Capabilities.MaxOverSubscriptionRatio).String() == "string" {
				overcommit, err = strconv.ParseFloat(storagePool.Capabilities.MaxOverSubscriptionRatio.(string), 64)
			}else {
				overcommit = storagePool.Capabilities.MaxOverSubscriptionRatio.(float64)
			}
			if (err != nil ){ // can't get MaxOverSubscriptionRatio, it can be auto. bypass.
			    log.Warn("bypass err: %s ",)
				fields := fieldMap{
					"total_capacity_gb":       storagePool.Capabilities.TotalCapacityGb,
					"free_capacity_gb":        storagePool.Capabilities.FreeCapacityGb,
					"allocated_capacity_gb":   float64(storagePool.Capabilities.AllocatedCapacityGb),
					"provisioned_capacity_gb": float64(storagePool.Capabilities.ProvisionedCapacityGb),
				}
				acc.AddFields("openstack_volumes", fields, tags)
			}

			fields := fieldMap{
				"total_capacity_gb":       storagePool.Capabilities.TotalCapacityGb,
				"free_capacity_gb":        storagePool.Capabilities.FreeCapacityGb,
				"allocated_capacity_gb":   float64(storagePool.Capabilities.AllocatedCapacityGb),
				"provisioned_capacity_gb": float64(storagePool.Capabilities.ProvisionedCapacityGb),
				"disk_overcommit_ratio":   overcommit,
			}
			acc.AddFields("openstack_volumes", fields, tags)
		}
	}
	return err
}

//
func (o *OpenStack) getVolumeProjectQuotas(projectID string) (fieldMap, error) {
	var err error
	var fields fieldMap
	blockstorageQuotas, err := blockstorageQuotas.Detail(o.volumeClient, projectID)
	if err != nil {

	} else {
		if (blockstorageQuotas.Volumes.Limit == -1) {
			blockstorageQuotas.Volumes.Limit = 99999
		}
		if (blockstorageQuotas.Gigabytes.Limit == -1) {
			blockstorageQuotas.Gigabytes.Limit = 99999
		}
		if (blockstorageQuotas.Snapshots.Limit == -1) {
			blockstorageQuotas.Snapshots.Limit = 99999
		}
		fields = fieldMap{
			"volumes_limit": blockstorageQuotas.Volumes.Limit,

			"volumes_inUse":    blockstorageQuotas.Volumes.InUse,
			"volumes_limit_gb": blockstorageQuotas.Gigabytes.Limit,
			"volumes_inUse_gb": blockstorageQuotas.Gigabytes.InUse,
			"snapshot_limit": blockstorageQuotas.Snapshots.Limit,
			"snapshot_inUse": blockstorageQuotas.Snapshots.InUse,
		}
	}
	return fields, err
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

//

// main func
func (o *OpenStack) Gather(acc telegraf.Accumulator) error {
	var err error
	if o.filter == nil {
		if o.filter, err = filter.Compile(o.ServicesGather); err != nil {
			return fmt.Errorf("error compiling filter: %s", err)
		}
	}
	if !o.initialized {
		err := o.initialize()
		if err != nil {
			acc.AddFields("openstack_identity", fieldMap{"api_state": 0,}, tagMap{
				"cloud":  o.Cloud,
				"region": o.Region,
			})
			return fmt.Errorf("failed to init : %s", err)
		}
	}
	// Gather resources.  Note tags harvesting must come first
	// gatherers are dependant on this information.
	gatherers := map[string]func() error{
		"services": o.gatherServices,
		"projects": o.gatherProjects,
		"users":    o.gatherUsers,
		"group":    o.gatherGroups,
	}
	for resources, gatherer := range gatherers {
		if err := gather(gatherer); err != nil {
			acc.AddFields("openstack_identity", fieldMap{"api_state": 0,}, tagMap{
				"cloud":  o.Cloud,
				"region": o.Region,
			})
			return fmt.Errorf("failed to get %s  : %s ", resources, err)
		}
	}

	// Accumulate statistics . depend on which service will cralw
	o.accumulators = []func(telegraf.Accumulator) error{}
	if o.filter.Match("identity") {
		acc.AddFields("openstack_identity", fieldMap{"api_state": 1,}, tagMap{
			"cloud":  o.Cloud,
			"region": o.Region,
		})
		acc.AddError(o.accumulateIdentity(acc))
	}

	if o.filter.Match("compute") {
		o.accumulators = append(o.accumulators,
			o.accumulateComputeAgents,
			o.accumulateComputeHypervisors,
		)
	}

	if o.filter.Match("volumev3") {
		o.accumulators = append(o.accumulators,
			o.accumulateVolumeAgents,
			o.accumulateVolumeStoragePools,
		)
	}
	if o.filter.Match("network") {
		o.accumulators = append(o.accumulators,
			o.accumulateNetworkAgents,
			o.accumulateNetworkIp,
		)
	}
	o.accumulateProjectQuotas(acc)

	wg := sync.WaitGroup{}
	for _, accumulator := range o.accumulators {
		wg.Add(1)
		//go routine in here
		go func(accumulator func(telegraf.Accumulator) error) {
			defer wg.Done()
			acc.AddError(accumulator(acc))
		}(accumulator)
	}
	wg.Wait()
	return nil
}

// init registers a callback which creates a new OpenStack input instance.
func init() {
	inputs.Add("openstack", func() telegraf.Input {
		return &OpenStack{
			Cloud:              "my_openstack",
			Region:             "RegionOne",
			UserDomainID:       "default",
			ProjectDomainID:    "default",
			ServicesGather:     []string{},
			MemOvercommitRatio: 1.5,
			CpuOvercommitRatio: 16.0,
		}
	})
}

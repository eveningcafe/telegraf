package openstack

import (
	"fmt"
	"github.com/influxdata/telegraf/filter"
	"github.com/influxdata/telegraf/internal/tls"
	"log"
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
	networkingFloatingIP "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/floatingips"
	networkingNET "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/networks"
	networkingQuotas "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2/quotas"
)

const (
	// plugin is used to identify ourselves in log output
	plugin = "openstack"
)
// Note: on debug, comment  api/base/client/NewClient( _time_out_of_http_request)

var sampleConfig = `
  ## This is the recommended interval to poll.
  interval = '20m'
  ## Region of openstack cluster which crawled by plugin
  region = "RegionOne"
  ## The identity endpoint to authenticate against openstack service, use v3 indentity api
  identity_endpoint = "https://my.openstack.cloud:5000"
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
  ## no way we can get overcomit cpu and ram so it must config as same as the config at nova.conf
  cpu_overcomit_ratio = "16"
  ram_overcomit_ratio = "1.5"
  ## Optional TLS Config
  tls_ca = "/etc/telegraf/openstack.crt"
  ## Use TLS but skip chain & host verification
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
	IdentityEndpoint string
	ProjectDomainID  string
	UserDomainID     string
	Project          string
	Username         string
	Password         string
	Region           string
	ServicesGather   []string
	CpuOvercomitRatio string
	RamOvercomitRatio string
	tls.ClientConfig

	//filter service
	filter filter.Filter

	// Locally cached clients
	identityClient *identity.IdentityClient
	computeClient  *compute.ComputeClient
	volumeClient   *blockstorage.VolumeClient
	networkClient  *networking.NetworkClient

	// Locally cached resources
	services serviceMap
	users    userMap
	groups   groupMap
	projects projectMap
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
	// Create required clients and attach to the OpenStack struct
	if o.filter.Match("identity") {
		if o.identityClient, err = identity.NewIdentityV3(*provider, o.Region); err != nil {
			err = fmt.Errorf("unable to create V3 identity client: %v", err)
		}
	}
	if o.filter.Match("compute") {
		if o.computeClient, err = compute.NewComputeV2(*provider, o.Region); err != nil {
			err = fmt.Errorf("unable to create V2 compute client: %v", err)
		}
	}
	if o.filter.Match("volumev3") {
		if o.volumeClient, err = blockstorage.NewBlockStorageV3(*provider, o.Region); err != nil {
			err = fmt.Errorf("unable to create V3 block storage client: %v", err)
		}
	}
	if o.filter.Match("network") {
		if o.networkClient, err = networking.NewNetworkClientV2(*provider, o.Region); err != nil {
			err = fmt.Errorf("unable to create V3 networking client: %v", err)
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

//
func (o *OpenStack) accumulateComputeAgents(acc telegraf.Accumulator) error{
	agents, err := computeServices.List(o.computeClient)
	if err != nil {
		acc.AddFields("openstack_compute", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
	} else {
		acc.AddFields("openstack_compute", fieldMap{"api_state": 1,}, tagMap{
			"region": o.Region,
		})

		fields := fieldMap{}
		for _, agent := range agents {
			if agent.State == "up" {
				fields["agent_state"] = 1
			} else {
				fields["agent_state"] = 0
			}
			acc.AddFields("openstack_compute", fields, tagMap{
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

// accumulateHypervisors accumulates statistics from hypervisors.
func (o *OpenStack) accumulateComputeHypervisors(acc telegraf.Accumulator) error{
	hypervisors, err := computeHypervisors.List(o.computeClient)
	if err != nil {
		acc.AddFields("openstack_compute", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
	} else {
		var RamOverCommitRatio float64
		var CpuOverCommitRatio float64
		for _, hypervisor := range hypervisors {
			RamOverCommitRatio , err = strconv.ParseFloat(o.RamOvercomitRatio, 64)
			RamOverCommit := hypervisor.MemoryMb*RamOverCommitRatio

			CpuOverCommitRatio , err = strconv.ParseFloat(o.CpuOvercomitRatio, 64)
			CpuOverCommit := hypervisor.Vcpus*CpuOverCommitRatio
			fields := fieldMap{
				"memory_mb_total":     hypervisor.MemoryMb,
				"memory_mb_used":      hypervisor.MemoryMbUsed,
				"running_vms":         hypervisor.RunningVms,
				"vcpus_total":         hypervisor.Vcpus,
				"vcpus_used":          hypervisor.VcpusUsed,
				"vcpus_total_overcommit": CpuOverCommit,
				"memory_total_overcommit": RamOverCommit,
				"local_disk_avalable": hypervisor.LocalGb,
				"local_disk_usage":    hypervisor.LocalGbUsed,
			}
			acc.AddFields("openstack_compute", fields, tagMap{
				"hypervisor_host": hypervisor.HypervisorHostname,
				"region":          o.Region,
			})
		}
	}
	return err
}

//
func (o *OpenStack) accumulateComputeProjectQuotas(acc telegraf.Accumulator) error {
	var err error
	for _, p := range o.projects {
		if p.Name == "service" {
			continue
		}
		computeQuotas, err := computeQuotas.Detail(o.computeClient, p.ID)
		if err != nil {
			acc.AddFields("openstack_compute", fieldMap{"api_state": 0,}, tagMap{
				"region": o.Region,
			})
		} else {

			if (computeQuotas.Cores.Limit == -1) {
				computeQuotas.Cores.Limit = 99999
			}
			if (computeQuotas.RAM.Limit == -1) {
				computeQuotas.RAM.Limit = 99999
			}
			if (computeQuotas.Instances.Limit == -1) {
				computeQuotas.Instances.Limit = 99999
			}

			acc.AddFields("openstack_compute", fieldMap{
				"cpu_limit":      computeQuotas.Cores.Limit,
				"cpu_used":       computeQuotas.Cores.InUse,
				"ram_limit":      computeQuotas.RAM.Limit,
				"ram_used":       computeQuotas.RAM.InUse,
				"instance_limit": computeQuotas.Instances.Limit,
				"instance_used":  computeQuotas.Instances.InUse,
			}, tagMap{
				"region":  o.Region,
				"project": p.Name,
			})
		}
	}
	return err
}

//
// accumulateIdentity accumulates statistics from the identity service.
func (o *OpenStack) accumulateIdentity(acc telegraf.Accumulator) error {

	fields := fieldMap{
		"num_projects": len(o.projects),
		"num_servives": len(o.services),
		"num_users":    len(o.users),
		"num_group":    len(o.groups),
	}
	acc.AddFields("openstack_identity", fields, tagMap{
		"region": o.Region,
	})
	return nil
}

//
func (o *OpenStack) accumulateNetworkAgents(acc telegraf.Accumulator) error{
	agents, err := networkingAgent.List(o.networkClient)
	if err != nil {
		acc.AddFields("openstack_network", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
	} else {
		acc.AddFields("openstack_network", fieldMap{"api_state": 1,}, tagMap{
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
				status = "enable"
			} else {
				status = "disable"
			}

			acc.AddFields("openstack_network", fields, tagMap{
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

//
func (o *OpenStack) accumulateNetworkFloatingIP(acc telegraf.Accumulator) error {
	floatingIps, err := networkingFloatingIP.List(o.networkClient)
	if err != nil {
		// bypass cause openstack can use provider network model
		acc.AddFields("openstack_network", fieldMap{
			"floating_ip": 0,
		}, tagMap{
			"region": o.Region,
		})

	} else {
		acc.AddFields("openstack_network", fieldMap{
			"floating_ip": len(floatingIps),
		}, tagMap{
			"region": o.Region,
		})
	}
	return nil
}

//// num_net and subnet, need tag network provider or not
//func (o *OpenStack) accumulateNetworkNET(acc telegraf.Accumulator) error {
//	networks, err := networkingNET.List(o.networkClient)
//	if err != nil {
//		acc.AddFields("openstack_network", fieldMap{"api_state": 0,}, tagMap{
//			"region": o.Region,
//		})
//		err =  fmt.Errorf("unable to get netowork info: %v", err)
//	}
//	for _, network := range networks {
//		acc.AddFields("openstack_network", fieldMap{
//			"num_subnet_in_network": len(network.Subnets),
//		}, tagMap{
//			"region":  o.Region,
//			"network": network.Name,
//		})
//	}
//	acc.AddFields("openstack_network", fieldMap{
//		"num_network": len(networks),
//	}, tagMap{
//		"region": o.Region,
//	})
//	return err
//}

//
func (o *OpenStack) accumulateNetworkIp(acc telegraf.Accumulator) error{
	// get provider network
	networks, err := networkingNET.List(o.networkClient)
	if err != nil {
		acc.AddFields("openstack_network", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
		return err
	}
	providerNetworkMap := networkMap{}
	for _, network := range networks {
		if network.ProviderPhysicalNetwork != ""{
			providerNetworkMap[network.ID] = network
		}
	}

	// get provider network pool
	ipAvail, err := networkingNET.NetworkIPAvailabilities(o.networkClient)
	if err != nil {
		acc.AddFields("openstack_network", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
		err =  fmt.Errorf("unable to get network info: %v", err)

	} else {
		var providerNetwork string
		for _, ipAvailNet := range ipAvail {

			if p, ok := providerNetworkMap[ipAvailNet.NetworkID]; ok {
				providerNetwork = p.Name
			}else {
				continue
			}

			// for each in loop of subnet
			for _, ipAvalSubnet := range ipAvailNet.SubnetIPAvailability {
				tags := tagMap{
					"region":      o.Region,
					"network":     ipAvailNet.NetworkName,
					"subnet_cidr": ipAvalSubnet.Cidr,
					"provider_network":     providerNetwork,
				}
				acc.AddFields("openstack_network", fieldMap{
					"ip_total": ipAvalSubnet.TotalIps,
					"ip_used":  ipAvailNet.UsedIps,
				}, tags)
			}

			// overall outsite
			acc.AddFields("openstack_network", fieldMap{
				"ip_used":  ipAvailNet.UsedIps,
				"ip_total": ipAvailNet.TotalIps,
			}, tagMap{
				"region":      o.Region,
				"network":     ipAvailNet.NetworkName,
				"subnet_cidr": "all",
				"provider_network":     providerNetwork,
			})

		}
	}
	return err
}

//
func (o *OpenStack) accumulateNetworkProjectQuotas(acc telegraf.Accumulator) error {
	var err error
	for _, p := range o.projects {
		if p.Name == "service" {
			continue
		}
		netQuotas, err := networkingQuotas.Detail(o.networkClient, p.ID)
		if err != nil {
		} else {

			if (netQuotas.Network.Limit == -1) {
				netQuotas.Network.Limit = 99999
			}
			if (netQuotas.SecurityGroup.Limit == -1) {
				netQuotas.SecurityGroup.Limit = 99999
			}
			if (netQuotas.Subnet.Limit == -1) {
				netQuotas.Subnet.Limit = 99999
			}
			if (netQuotas.Port.Limit == -1) {
				netQuotas.Port.Limit = 99999
			}

			acc.AddFields("openstack_network", fieldMap{
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
			}, tagMap{
				"region":  o.Region,
				"project": p.Name,
			})
		}
	}
	return err
}

//
func (o *OpenStack) accumulateVolumeAgents(acc telegraf.Accumulator) error {
	agents, err := blockstorageServices.List(o.volumeClient)
	if err != nil {
		acc.AddFields("openstack_volumes", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
	} else {
		acc.AddFields("openstack_volumes", fieldMap{"api_state": 1,}, tagMap{
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
		acc.AddFields("openstack_volumes", fieldMap{"api_state": 0,}, tagMap{
			"region": o.Region,
		})
	} else {
		for _, storagePool := range storagePools {
			tags := tagMap{
				"backend_state": storagePool.Capabilities.BackendState,
				"backend_name":  storagePool.Capabilities.VolumeBackendName,
				"region":        o.Region,
			}
			overcommit, _ := strconv.ParseFloat(storagePool.Capabilities.MaxOverSubscriptionRatio, 64)
			fields := fieldMap{
				"total_capacity_gb":           storagePool.Capabilities.TotalCapacityGb,
				"free_capacity_gb":            storagePool.Capabilities.FreeCapacityGb,
				"allocated_capacity_gb":       float64(storagePool.Capabilities.AllocatedCapacityGb),
				"provisioned_capacity_gb":     float64(storagePool.Capabilities.ProvisionedCapacityGb),
				"max_over_subscription_ratio": overcommit,
			}
			acc.AddFields("openstack_storage_pool", fields, tags)
		}
	}
	return err
}

//
func (o *OpenStack) accumulateVolumeProjectQuotas(acc telegraf.Accumulator) error {
	var err error
	for _, p := range o.projects {
		if p.Name == "service" {
			continue
		}
		blockstorageQuotas, err := blockstorageQuotas.Detail(o.volumeClient, p.ID)
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
			acc.AddFields("openstack_volumes", fieldMap{
				"volumes_limit":         blockstorageQuotas.Volumes.Limit,
				"volumes_allocated":     blockstorageQuotas.Volumes.Allocated,
				"volumes_inUse":         blockstorageQuotas.Volumes.InUse,
				"volumes_limit_gb":      blockstorageQuotas.Gigabytes.Limit,
				"volumes_inUse_gb":      blockstorageQuotas.Gigabytes.InUse,
				"volummes_allocated_gb": blockstorageQuotas.Gigabytes.Allocated,
				"snapshot_limit":        blockstorageQuotas.Snapshots.Limit,
				"snapshot_inUse":        blockstorageQuotas.Snapshots.InUse,
				"snapshot_allocated":    blockstorageQuotas.Snapshots.Allocated,
			}, tagMap{
				"region":  o.Region,
				"project": p.Name,
			})
		}
	}
	return err
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
func (o *OpenStack) doAccumelate(accumulators []func(telegraf.Accumulator) error, acc telegraf.Accumulator) error{
	errorChannel := make (chan error)
	for _, accumulator := range accumulators {
		//acc.AddError(accumulator(acc))
		//accumulator(acc)
		//go routine in here
		go func(accumulator func(telegraf.Accumulator) error) {
			err := accumulator(acc)
			errorChannel <- err
		}(accumulator)
	}
	err := <- errorChannel
	return err
}
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
				"region": o.Region,
			})
			return fmt.Errorf("failed to init : %s", err)
		}

	}

	acc.AddFields("openstack_identity", fieldMap{"api_state": 1,}, tagMap{
		"region": o.Region,
	})
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
				"region": o.Region,
			})

			log.Println("W!", plugin, "failed to get", resources, " : ", err)
		}
	}
	if o.filter.Match("identity") {
		acc.AddError(o.accumulateIdentity(acc))
	}
	// Accumulate statistics . depend on which service will cralw
	accumulators := []func(telegraf.Accumulator) error{}

	if o.filter.Match("compute") {
		accumulators = append(accumulators,
			o.accumulateComputeAgents,
			o.accumulateComputeHypervisors,
			o.accumulateComputeProjectQuotas)
	}

	o.doAccumelate(accumulators, acc)
	if acc.()

	if o.filter.Match("volumev3") {
		accumulators = append(accumulators,
			o.accumulateVolumeAgents,
			o.accumulateVolumeStoragePools,
			o.accumulateVolumeProjectQuotas)
	}
	if o.filter.Match("network") {
		o.accumulators = append(o.accumulators,
			o.accumulateNetworkAgents,
			o.accumulateNetworkFloatingIP,
			o.accumulateNetworkIp,
			o.accumulateNetworkProjectQuotas)
	}


	for _, accumulator := range o.accumulators {
		//acc.AddError(accumulator(acc))
		//accumulator(acc)
		wg.Add(1)
		//go routine in here
		go func(accumulator func(telegraf.Accumulator) error ) {
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
			UserDomainID:    "default",
			ProjectDomainID: "default",
			ServicesGather:  []string{},
			RamOvercomitRatio: "1",
			CpuOvercomitRatio: "1",
		}
	})
}

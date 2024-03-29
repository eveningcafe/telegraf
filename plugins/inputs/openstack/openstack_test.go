package openstack_test

import (
	"fmt"
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/blockstorage"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/compute"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/indentity"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/networking"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/placement"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/influxdata/telegraf/plugins/inputs/openstack"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestOpenStackCluster(t *testing.T) {
	var err error
	var fakeKeystoneListen net.Listener
	var fakeNovaListen net.Listener
	var fakePlacementListen net.Listener
	var fakeCinderListen net.Listener
	var fakeNeutronListen net.Listener
	fakeKeystoneEndpoint := "http://127.0.0.1:5000"
	fakeNovaEndpoint := "http://127.0.0.1:8774"
	fakePlacementEndpoint := "http://127.0.0.1:8778"
	fakeCinderEndpoint := "http://127.0.0.1:8776"
	fakeNeutronEndpoint := "http://127.0.0.1:9696"

	//try to listen on server which run unit test
	fakeKeystoneListen, err = net.Listen("tcp", fakeKeystoneEndpoint[7:])
	fakeNovaListen, err = net.Listen("tcp", fakeNovaEndpoint[7:])
	fakePlacementListen, err = net.Listen("tcp", fakePlacementEndpoint[7:])
	fakeCinderListen, err = net.Listen("tcp", fakeCinderEndpoint[7:])
	fakeNeutronListen, err = net.Listen("tcp", fakeNeutronEndpoint[7:])
	if err != nil {
		fmt.Println(err)
	}
	defer fakeKeystoneListen.Close()
	defer fakeNovaListen.Close()
	defer fakeCinderListen.Close()
	defer fakeNeutronListen.Close()
	defer fakePlacementListen.Close()

	if err != nil {
		log.Fatal(err)
	}
	// fake openstack api
	fakeKeystoneServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/auth/tokens" {
			if r.Method == "POST" {
				w.Header().Set("X-Subject-Token", "Special-Test-Token")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(
					indentity.CreateTokenResponseBody(
						fakeKeystoneEndpoint,
						fakeNovaEndpoint,
						fakePlacementEndpoint,
						fakeCinderEndpoint,
						fakeNeutronEndpoint)))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v3/services" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.ServiceListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v3/projects" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.ProjectListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v3/users" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.UserListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v3/groups" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.GroupListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	fakeKeystoneServer.Listener = fakeKeystoneListen
	fakeKeystoneServer.Start()
	defer fakeKeystoneServer.Close()

	fakeNovaServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/os-services" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(compute.ServiceListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/os-hypervisors/detail" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(compute.HypervisorsListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		} else if r.URL.Path == "/os-quota-sets/"+resources.ProjectId+"/detail" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(compute.QuotasListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	fakeNovaServer.Listener = fakeNovaListen
	fakeNovaServer.Start()
	defer fakeNovaServer.Close()

	// placement
	fakePlacementServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/resource_providers" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(placement.ListResourcesProviderResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/resource_providers/"+resources.HypervisorID+"/inventories" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(placement.GetResourcesInventoriesResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		} else if r.URL.Path == "/resource_providers/"+resources.HypervisorID+"/usages" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(placement.GetResourcesUsagesResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	fakePlacementServer.Listener = fakePlacementListen
	fakePlacementServer.Start()
	defer fakePlacementServer.Close()

	fakeCinderServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/os-services" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(blockstorage.ServiceListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/scheduler-stats/get_pools" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(blockstorage.StoragePoolListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/os-quota-sets/"+resources.ProjectId {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(blockstorage.QuotasListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	fakeCinderServer.Listener = fakeCinderListen
	fakeCinderServer.Start()
	defer fakeCinderServer.Close()

	fakeNeutronServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2.0/agents" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(networking.AgentsListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v2.0/floatingips" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(networking.FloatingIpsListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v2.0/networks" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(networking.NetworkListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/v2.0/network-ip-availabilities" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(networking.IpAvailabilityListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		} else if r.URL.Path == "/v2.0/quotas/"+resources.ProjectId+"/details.json" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(networking.QuotasListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	fakeNeutronServer.Listener = fakeNeutronListen
	fakeNeutronServer.Start()
	defer fakeNeutronServer.Close()

	plugin := &plugin.OpenStack{
		IdentityEndpoint: fakeKeystoneServer.URL,
		Cloud:            "my_openstack",
		Region:           "RegionOne",
		ServicesGather: []string{
			"identity",
			"volumev3",
			"network",
			"compute",
		},
		ProjectDomainID: "default",
		UserDomainID:    "default",
		ClientConfig: tls.ClientConfig{
			InsecureSkipVerify: false,
			TLSCA:              "test/resources/openstack.crt"},
	}

	var acc testutil.Accumulator
	require.NoError(t, acc.GatherError(plugin.Gather))

	//acc.AssertContainsFields(t, "openstack_identity",map[string]interface {}{"api_state":1})
	//acc.AssertContainsFields(t, "openstack_compute", map[string]interface {}{"api_state":1})
	//acc.AssertContainsFields(t, "openstack_volumes", map[string]interface {}{"api_state":1})
	//acc.AssertContainsFields(t, "openstack_network", map[string]interface {}{"api_state":1})

	iFields := map[string]interface{}{
		"num_projects": 1,
		"num_servives": 7,
		"num_users":    7,
		"num_group":    1,
	}
	iTags := map[string]string{
		"cloud":  "my_openstack",
		"region": "RegionOne",
	}
	acc.AssertContainsTaggedFields(t, "openstack_identity", iFields, iTags)

	cFields := map[string]interface{}{
		"memory_used_mb":        float64(16384),
		"mem_overcommit_ratio":  float64(1.5),
		"memory_total_mb":       float64(31785),
		"running_vms":           2,
		"cpu_total":             float64(32),
		"disk_overcommit_ratio": float64(1),
		"cpu_overcommit_ratio":  float64(16),
		"local_disk_reserved_gb":   float64(0),
		"hypervisor_workload":   float64(0),
		"local_disk_usage_gb":      float64(0),
		"local_disk_total_gb":      float64(22354),
		"cpu_used":             float64(16),
		"cpu_reserved":          float64(0),
		"memory_reserved_mb":    float64(512),
	}
	cTags := map[string]string{
		"hypervisor_host":   "compute05",
		"hypervisor_state":  "up",
		"hypervisor_status": "enabled",
		"cloud":             "my_openstack",
		"region":            "RegionOne",
	}
	acc.AssertContainsTaggedFields(t, "openstack_compute", cFields, cTags)

	qFields := map[string]interface{}{
		"network_limit":       100,
		"network_used":        1,
		"securityGroup_limit": 10,
		"securityGroup_used":  1,
		"securityRule_limit":  100,
		"securityRule_used":   6,
		"subnet_limit":        100,
		"subnet_used":         1,
		"port_limit":          500,
		"port_used":           1,
		"floatingIP_limit":    9999,
		"floatingIP_used":     0,

		"snapshot_inUse":   0,
		"volumes_inUse":    0,
		"volumes_limit_gb": 1000,
		"volumes_inUse_gb": 0,
		"volumes_limit":    10,
		"snapshot_limit":   10,

		"cpu_limit":      20,
		"cpu_used":       0,
		"ram_limit":      51200,
		"ram_used":       0,
		"instance_limit": 10,
		"instance_used":  0,
	}
	qTags := map[string]string{
		"project": "demo",
		"cloud":   "my_openstack",
		"region":  "RegionOne",
	}

	acc.AssertContainsTaggedFields(t, "openstack_quotas", qFields, qTags)

	vsFields := map[string]interface{}{
		"total_capacity_gb":       float64(125.03),
		"free_capacity_gb":        float64(125.03),
		"allocated_capacity_gb":   float64(0),
		"provisioned_capacity_gb": float64(0),
		"disk_overcommit_ratio":   float64(20),
	}
	vsTags := map[string]string{
		"pool_name": "controller@ceph#RBD",
		"region":    "RegionOne",
		"cloud":     "my_openstack",
	}

	acc.AssertContainsTaggedFields(t, "openstack_volumes", vsFields, vsTags)

	nFields := map[string]interface{}{
		"ip_used":  int64(1),
		"ip_total": int64(52),
	}
	nTags := map[string]string{
		"subnet_cidr":      "192.168.33.0/24",
		"cloud":            "my_openstack",
		"region":           "RegionOne",
		"provider_network": "provider",
		"network":          "provider",
	}

	acc.AssertContainsTaggedFields(t, "openstack_network", nFields, nTags)

}

func TestOpenstackInRealOpenstack(t *testing.T) {

	plugin := &plugin.OpenStack{
		IdentityEndpoint: "https://controller:5000/v3",
		Project:          "admin",
		UserDomainID:     "default",
		ProjectDomainID:  "default",
		Password:         "Welcome123",
		Username:         "admin",
		Cloud:            "my_openstack",
		Region:           "RegionOne",
		ServicesGather:   []string{"identity", "volumev3", "compute", "network"},
		ClientConfig: tls.ClientConfig{
			InsecureSkipVerify: false,
			TLSCA:              "test/resources/openstack.crt"},
	}
	var acc testutil.Accumulator
	err := acc.GatherError(plugin.Gather)
	require.NoError(t, err)

}

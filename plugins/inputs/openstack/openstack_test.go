package openstack_test

import (
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/blockstorage"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/compute"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/indentity"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources/networking"
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
	var fakeCinderListen net.Listener
	var fakeNeutronListen net.Listener
	fakeKeystoneEndpoint := "http://127.0.0.1:5000"
	fakeNovaEndpoint := "http://127.0.0.1:8774"
	fakeCinderEndpoint := "http://127.0.0.1:8776"
	fakeNeutronEndpoint := "http://127.0.0.1:9696"

	//try to listen on server which run unit test
	fakeKeystoneListen, err = net.Listen("tcp", fakeKeystoneEndpoint[7:])
	fakeNovaListen, err = net.Listen("tcp", fakeNovaEndpoint[7:])
	fakeCinderListen, err = net.Listen("tcp", fakeCinderEndpoint[7:])
	fakeNeutronListen, err = net.Listen("tcp", fakeNeutronEndpoint[7:])

	if err != nil {
		log.Fatal(err)
	}
	// fake openstack api
	fakeKeystoneServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/tokens" {
			if r.Method == "POST" {
				w.Header().Set("X-Subject-Token", "Special-Test-Token")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(
					indentity.CreateTokenResponseBody(
						fakeKeystoneEndpoint,
						fakeNovaEndpoint,
						fakeCinderEndpoint,
						fakeNeutronEndpoint)))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/services" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.ServiceListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/projects" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.ProjectListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/users" {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(indentity.UserListResponseBody()))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else if r.URL.Path == "/groups" {
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
		Region:           "RegionOne",
		ServicesGather:   []string{"identity", "volumev3", "compute", "network"},
		ClientConfig: tls.ClientConfig{
			InsecureSkipVerify: false,
			TLSCA:              "test/resources/openstack.crt"},
	}

	var acc testutil.Accumulator
	require.NoError(t, acc.GatherError(plugin.Gather))


	acc.AssertContainsFields(t, "openstack_identity",map[string]interface {}{"api_state":1})
	acc.AssertContainsFields(t, "openstack_compute", map[string]interface {}{"api_state":1})
	acc.AssertContainsFields(t, "openstack_volumes", map[string]interface {}{"api_state":1})
	acc.AssertContainsFields(t, "openstack_network", map[string]interface {}{"api_state":1})

	iFields := map[string]interface{}{
		"num_projects": 1,
		"num_servives": 7,
		"num_users":   7,
		"num_group":  1,
	}
	iTags := map[string]string{
		"region" : "RegionOne",
	}
	acc.AssertContainsTaggedFields(t, "openstack_identity", iFields, iTags)

	cFields := map[string]interface{}{
		"local_disk_usage": 0,
		"memory_mb_total": 7976,
		"memory_mb_used":   512,
		"running_vms":  0,
		"vcpus_total": 6,
		"vcpus_used": 0,
		"local_disk_avalable": 410,
	}
	cTags := map[string]string{
		"hostname": "compute01",
		"region" : "RegionOne",
	}
	acc.AssertContainsTaggedFields(t, "openstack_compute", cFields, cTags)

	vFields := map[string]interface{}{
		"snapshot_inUse": 0,
		"volumes_allocated": 0,
		"volumes_inUse":   0,
		"volumes_limit_gb":  1000,
		"volumes_inUse_gb": 0,
		"volummes_allocated_gb": 0,
		"volumes_limit": 10,
		"snapshot_limit": 10,
		"snapshot_allocated": 0,
	}
	vTags := map[string]string{
		"project": "demo",
		"region" : "RegionOne",
	}

	acc.AssertContainsTaggedFields(t, "openstack_volumes", vFields, vTags)
	//
	sFields := map[string]interface{}{
		"total_capacity_gb": float64(125.03),
		"free_capacity_gb": float64(125.03),
		"allocated_capacity_gb":   float64(0),
		"provisioned_capacity_gb":  float64(0),
		"max_over_subscription_ratio": float64(20),
	}
	sTags := map[string]string{
		"backend_state": "up",
		"backend_name": "RBD",
		"region" : "RegionOne",
	}

	acc.AssertContainsTaggedFields(t, "openstack_storage_pool", sFields, sTags)

	nFields := map[string]interface{}{
	   "ip_used": 1,
		"ip_total": 52,
    }
	nTags := map[string]string{
		"subnet_cidr":"all",
		"region":"RegionOne",
		"project":"unknown",
		"network":"provider",
	}

	acc.AssertContainsTaggedFields(t, "openstack_network", nFields, nTags)


}

//func TestOpenstackInReal(t *testing.T) {
//
//	plugin := &plugin.OpenStack{
//		IdentityEndpoint: "https://controller:5000/v3",
//		Project:          "admin",
//		UserDomainID:     "default",
//		ProjectDomainID:  "default",
//		Password:         "Welcome123",
//		Username:         "admin",
//		Region:           "RegionOne",
//		ServicesGather:   []string{"identity", "volumev3", "compute", "network"},
//		ClientConfig: tls.ClientConfig{
//			InsecureSkipVerify: false,
//			TLSCA:              "test/resources/openstack.crt"},
//	}
//	var acc testutil.Accumulator
//	require.NoError(t, acc.GatherError(plugin.Gather))
//}

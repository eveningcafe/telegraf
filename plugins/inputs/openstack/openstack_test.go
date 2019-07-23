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

func TestOpenstackInReal(t *testing.T) {

	plugin := &plugin.OpenStack{
		IdentityEndpoint: "https://controller:5000/v3",
		Project:          "admin",
		UserDomainID:     "default",
		ProjectDomainID:  "default",
		Password:         "Welcome123",
		Username:         "admin",
		Region:           "RegionOne",
		ServicesGather:   []string{"identity", "volumev3", "compute", "network"},
		ClientConfig: tls.ClientConfig{
			InsecureSkipVerify: false,
			TLSCA:              "test/resources/openstack.crt"},
	}
	//metricName := "openstack"
	var acc testutil.Accumulator

	require.NoError(t, acc.GatherError(plugin.Gather))
	//
	//// basic check to see if we got the right field, value and tag
	//var metric = acc.Metrics[0]
	////require.Equal(t, metric.Measurement, metricName)
	////require.Len(t, acc.Metrics[0].Fields, 1)
	////require.Equal(t, acc.Metrics[0].Fields["a"], 1.2)
}

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
	var metric = acc.Metrics[0]
	require.Equal(t, metric.Measurement, "asd")

}

//
//func TestHTTPHeaders(t *testing.T) {
//	header := "X-Special-Header"
//	headerValue := "Special-Value"
//	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.RequestBody) {
//		if r.URL.Path == "/endpoint" {
//			if r.Header.Get(header) == headerValue {
//				_, _ = w.Write([]byte(simpleJSON))
//			} else {
//				w.WriteHeader(http.StatusForbidden)
//			}
//		} else {
//			w.WriteHeader(http.StatusNotFound)
//		}
//	}))
//	defer fakeServer.Close()
//
//	url := fakeServer.URL + "/endpoint"
//	plugin := &plugin.HTTP{
//		URLs:    []string{url},
//		Headers: map[string]string{header: headerValue},
//	}
//
//	p, _ := parsers.NewParser(&parsers.Config{
//		DataFormat: "json",
//		MetricName: "metricName",
//	})
//	plugin.SetParser(p)
//
//	var acc testutil.Accumulator
//	plugin.Init()
//	require.NoError(t, acc.GatherError(plugin.Gather))
//}
//
//func TestInvalidStatusCode(t *testing.T) {
//	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.RequestBody) {
//		w.WriteHeader(http.StatusNotFound)
//	}))
//	defer fakeServer.Close()
//
//	url := fakeServer.URL + "/endpoint"
//	plugin := &plugin.HTTP{
//		URLs: []string{url},
//	}
//
//	metricName := "metricName"
//	p, _ := parsers.NewParser(&parsers.Config{
//		DataFormat: "json",
//		MetricName: metricName,
//	})
//	plugin.SetParser(p)
//
//	var acc testutil.Accumulator
//	plugin.Init()
//	require.Error(t, acc.GatherError(plugin.Gather))
//}
//
//func TestMethod(t *testing.T) {
//	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.RequestBody) {
//		if r.Method == "POST" {
//			w.WriteHeader(http.StatusOK)
//		} else {
//			w.WriteHeader(http.StatusNotFound)
//		}
//	}))
//	defer fakeServer.Close()
//
//	plugin := &plugin.HTTP{
//		URLs:   []string{fakeServer.URL},
//		Method: "POST",
//	}
//
//	p, _ := parsers.NewParser(&parsers.Config{
//		DataFormat: "json",
//		MetricName: "metricName",
//	})
//	plugin.SetParser(p)
//
//	var acc testutil.Accumulator
//	plugin.Init()
//	require.NoError(t, acc.GatherError(plugin.Gather))
//}
//
//const simpleJSON = `
//{
//    "a": 1.2
//}
//`
//
//func TestBodyAndContentEncoding(t *testing.T) {
//	ts := httptest.NewServer(http.NotFoundHandler())
//	defer ts.Close()
//
//	url := fmt.Sprintf("http://%s", ts.Listener.Addr().String())
//
//	tests := []struct {
//		name             string
//		plugin           *plugin.HTTP
//		queryHandlerFunc func(t *testing.T, w http.ResponseWriter, r *http.RequestBody)
//	}{
//		{
//			name: "no body",
//			plugin: &plugin.HTTP{
//				Method: "POST",
//				URLs:   []string{url},
//			},
//			queryHandlerFunc: func(t *testing.T, w http.ResponseWriter, r *http.RequestBody) {
//				body, err := ioutil.ReadAll(r.Body)
//				require.NoError(t, err)
//				require.Equal(t, []byte(""), body)
//				w.WriteHeader(http.StatusOK)
//			},
//		},
//		{
//			name: "post body",
//			plugin: &plugin.HTTP{
//				URLs:   []string{url},
//				Method: "POST",
//				Body:   "test",
//			},
//			queryHandlerFunc: func(t *testing.T, w http.ResponseWriter, r *http.RequestBody) {
//				body, err := ioutil.ReadAll(r.Body)
//				require.NoError(t, err)
//				require.Equal(t, []byte("test"), body)
//				w.WriteHeader(http.StatusOK)
//			},
//		},
//		{
//			name: "get method body is sent",
//			plugin: &plugin.HTTP{
//				URLs:   []string{url},
//				Method: "GET",
//				Body:   "test",
//			},
//			queryHandlerFunc: func(t *testing.T, w http.ResponseWriter, r *http.RequestBody) {
//				body, err := ioutil.ReadAll(r.Body)
//				require.NoError(t, err)
//				require.Equal(t, []byte("test"), body)
//				w.WriteHeader(http.StatusOK)
//			},
//		},
//		{
//			name: "gzip encoding",
//			plugin: &plugin.HTTP{
//				URLs:            []string{url},
//				Method:          "GET",
//				Body:            "test",
//				ContentEncoding: "gzip",
//			},
//			queryHandlerFunc: func(t *testing.T, w http.ResponseWriter, r *http.RequestBody) {
//				require.Equal(t, r.Header.Get("Content-Encoding"), "gzip")
//
//				gr, err := gzip.NewReader(r.Body)
//				require.NoError(t, err)
//				body, err := ioutil.ReadAll(gr)
//				require.NoError(t, err)
//				require.Equal(t, []byte("test"), body)
//				w.WriteHeader(http.StatusOK)
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.RequestBody) {
//				tt.queryHandlerFunc(t, w, r)
//			})
//
//			parser, err := parsers.NewParser(&parsers.Config{DataFormat: "influx"})
//			require.NoError(t, err)
//
//			tt.plugin.SetParser(parser)
//
//			var acc testutil.Accumulator
//			tt.plugin.Init()
//			err = tt.plugin.Gather(&acc)
//			require.NoError(t, err)
//		})
//	}
//}

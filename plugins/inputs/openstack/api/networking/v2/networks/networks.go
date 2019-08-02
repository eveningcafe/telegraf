package networks

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
	"math/big"
)

type Network struct {
	AdminStateUp            bool          `json:"admin_state_up"`
	AvailabilityZoneHints   []interface{} `json:"availability_zone_hints"`
	AvailabilityZones       []string      `json:"availability_zones"`
	CreatedAt               string        `json:"created_at"`
	DNSDomain               string        `json:"dns_domain"`
	ID                      string        `json:"id"`
	Ipv4AddressScope        interface{}   `json:"ipv4_address_scope"`
	Ipv6AddressScope        interface{}   `json:"ipv6_address_scope"`
	L2Adjacency             bool          `json:"l2_adjacency"`
	Mtu                     int           `json:"mtu"`
	Name                    string        `json:"name"`
	PortSecurityEnabled     bool          `json:"port_security_enabled"`
	ProjectID               string        `json:"project_id"`
	QosPolicyID             string        `json:"qos_policy_id"`
	ProviderNetworkType     string        `json:"provider:network_type"`
	ProviderPhysicalNetwork string        `json:"provider:physical_network"`
	ProviderSegmentationID  int           `json:"provider:segmentation_id"`
	RevisionNumber          int           `json:"revision_number"`
	RouterExternal          bool          `json:"router:external"`
	Shared                  bool          `json:"shared"`
	Status                  string        `json:"status"`
	Subnets                 []string      `json:"subnets"`
	Tags                    []string      `json:"tags"`
	TenantID                string        `json:"tenant_id"`
	UpdatedAt               string        `json:"updated_at"`
	VlanTransparent         bool          `json:"vlan_transparent"`
	Description             string        `json:"description"`
	IsDefault               bool          `json:"is_default"`
	Segments                []struct {
		ProviderNetworkType     string `json:"provider:network_type"`
		ProviderPhysicalNetwork string `json:"provider:physical_network"`
		ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	} `json:"segments,omitempty"`
}

//
type IPAvailabilities struct {
		NetworkID            string `json:"network_id"`
		NetworkName          string `json:"network_name"`
		SubnetIPAvailability []struct {
			Cidr       string `json:"cidr"`
			IPVersion  int    `json:"ip_version"`
			SubnetID   string `json:"subnet_id"`
			SubnetName string `json:"subnet_name"`
			TotalIps   big.Int  `json:"total_ips"`
			UsedIps    big.Int    `json:"used_ips"`
		} `json:"subnet_ip_availability"`
		ProjectID string `json:"project_id"`
		TenantID  string `json:"tenant_id"`
		TotalIps  big.Int    `json:"total_ips"`
		UsedIps   big.Int    `json:"used_ips"`
}

func List(client *v2.NetworkClient) ([]Network, error) {
	api, err := declareListNetwork(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return nil, err
	}
	result := ListNetworkResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	networks := []Network{}
	for _, v := range result.Networks{
		networks = append(networks, v)
	}
	return networks, err
}

func NetworkIPAvailabilities(client *v2.NetworkClient) ([]IPAvailabilities, error) {
	api, err := declareNetworkIPAvailabilities(client.Endpoint, client.Token)
	err = client.DoReuest(api)
	if err != nil {
		return nil, err
	}
	result := NetworkIPAvailabilitiesResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)

	ipAvail := []IPAvailabilities{}
	for _, v := range result.NetworkIPAvailabilities{
		ipAvail = append(ipAvail, v)
	}
	return ipAvail, err
}


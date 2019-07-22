package floatingips

import (
	"encoding/json"
	v2 "github.com/influxdata/telegraf/plugins/inputs/openstack/api/networking/v2"
	"time"
)

type Agent struct {
	ID               string
	Binary           string
	Host             string
	AdminStateUp     bool
	Configurations   interface{}
	Alive            bool
	AgentType        string
	AvailabilityZone string
	Topic            string
}
type Port struct {
	Status       string `json:"status"`
	Name         string `json:"name"`
	AdminStateUp bool   `json:"admin_state_up"`
	NetworkID    string `json:"network_id"`
	DeviceOwner  string `json:"device_owner"`
	MacAddress   string `json:"mac_address"`
	DeviceID     string `json:"device_id"`
}
type Floatingip struct {
	RouterID          string        `json:"router_id"`
	Description       string        `json:"description"`
	DNSDomain         string        `json:"dns_domain"`
	DNSName           string        `json:"dns_name"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	RevisionNumber    int           `json:"revision_number"`
	ProjectID         string        `json:"project_id"`
	TenantID          string        `json:"tenant_id"`
	FloatingNetworkID string        `json:"floating_network_id"`
	FixedIPAddress    string        `json:"fixed_ip_address"`
	FloatingIPAddress string        `json:"floating_ip_address"`
	PortID            string        `json:"port_id"`
	ID                string        `json:"id"`
	Status            string        `json:"status"`
	PortDetails       Port          `json:"port_details,omitempty"`
	Tags              []string      `json:"tags"`
	PortForwardings   []interface{} `json:"port_forwardings"`
}

func List(client *v2.NetworkClient) ([]Floatingip, error) {
	api, err := declareListFloatingIp(client.Endpoint, client.Token)
	err = api.DoReuest()
	if err!=nil {
		return nil, err
	}
	result := ListFloatingIpResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody),&result)
	floatingips := []Floatingip{}
	for _, v := range result.Floatingips {
		floatingips = append(floatingips, v)
	}
	return floatingips, err
}

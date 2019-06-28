package openstach

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
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
type OpenStackConfig struct {
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
func (o *OpenStackConfig) SampleConfig() string {
	return sampleConfig
}

// Description returns a description string of the input plugin and implements
// the Input interface.
func (o *OpenStackConfig) Description() string {
	return "Collects performance metrics from OpenStack services"
}
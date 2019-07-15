package v2

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
)

type NetworkClient struct{
	*base.Client
}
func NewNetworkClientV3(providerClient authenticator.ProviderClient, region string) (*NetworkClient, error) {
	var err error
	c := &NetworkClient{}
	c.ServiceType = "network"
	c.Client, err  = base.NewClient( providerClient , region, c.ServiceType )
	return c, err
}
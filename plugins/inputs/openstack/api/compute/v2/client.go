package v2

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/authenticator"
)

type ComputeClient struct{
	*base.Client
}
func NewComputeV2(providerClient authenticator.ProviderClient, region string) (*ComputeClient, error) {
	var err error
	c := &ComputeClient{}
	c.Client, err  = base.NewClient( providerClient , region , "compute")
	return c, err
}




package v1

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/authenticator"
)

type PlacementClient struct{
	*base.Client
}
func NewPlacementClientV1(providerClient authenticator.ProviderClient, region string) (*PlacementClient, error) {
	var err error
	c := &PlacementClient{}
	c.Client, err  = base.NewClient( providerClient , region, "placement" )
	return c, err
}

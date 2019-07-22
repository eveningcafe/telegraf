package v3

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/authenticator"
)

type VolumeClient struct {
	*base.Client
}

func NewBlockStorageV3(providerClient authenticator.ProviderClient, region string) (*VolumeClient, error) {
	var err error
	c := &VolumeClient{}
	c.Client, err  = base.NewClient( providerClient , region , "volumev3")
	return c, err
}
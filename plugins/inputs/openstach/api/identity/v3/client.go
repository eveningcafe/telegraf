package v3

import (
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/base"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
)

type IdentityClient struct {
	*base.Client
}

func NewIdentityV3(providerClient authenticator.ProviderClient, region string) (*IdentityClient, error) {
	var err error
	c := &IdentityClient{}
	c.Client, err  = base.NewClient( providerClient , region , "identity")
	return c, err
}

package v3

import (
	"errors"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
)

type IdentityClient struct {
	Token       string
	Endpoint    string
	ServiceType string
	Region      string
}

func NewIdentityV3(providerClient authenticator.ProviderClient, region string) (*IdentityClient, error) {
	c := new(IdentityClient)
	c.ServiceType = "identity"
	c.Token = providerClient.Token
	for _, ca := range providerClient.Catalog {
		if ca.Type == c.ServiceType {
			for _, e := range ca.Endpoints {
				if e.Interface == "public" && e.Region == region {
					c.Endpoint = e.URL
					c.Region = e.Region
				}
			}
		}
	}
	if c.Endpoint == "" {
		return nil, errors.New("no service " + c.ServiceType + " avalable on region " + region)
	}
	return c, nil
}

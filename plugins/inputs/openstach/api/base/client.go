package base

import (
	"errors"
	"github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"
)

type Client struct{
	Token       string
	Endpoint    string
	ServiceType string
	Region      string
}
func NewClient(providerClient authenticator.ProviderClient, region string, serviceType string) (*Client, error) {
	c := new(Client)
	c.Token = providerClient.Token
	c.ServiceType = serviceType
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
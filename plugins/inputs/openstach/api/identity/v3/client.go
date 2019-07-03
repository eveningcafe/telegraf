package v3

import "github.com/influxdata/telegraf/plugins/inputs/openstach/api/identity/v3/authenticator"

type IdentityClient struct {
	Token         string
	Endpoint      string
	ServiceType   string
}

func NewIdentityV3(token string, catalog []authenticator.Catalog) (*IdentityClient, error){
	c := new(IdentityClient)
	c.ServiceType = "identity"
	c.Token = token
	for _, ca := range catalog {
		if ca.Type==c.ServiceType {
			for _,e :=range ca.Endpoints {
				if e.Interface =="public"{
					c.Endpoint = e.URL
				}

			}
		}
	}

	return c, nil
}

package base

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/base/request"
	"github.com/influxdata/telegraf/plugins/inputs/openstack/api/identity/v3/authenticator"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeout  = 10

type Client struct {
	Token       string
	Endpoint    string
	ServiceType string
	Region      string
	HTTPClient  *http.Client
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
	if serviceType == "identity" && !strings.Contains(c.Endpoint , "v3") {
		c.Endpoint = c.Endpoint+"/v3"
	}
	c.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: providerClient.TlsCfg,
		},
		Timeout: time.Second * timeout,
		//Timeout: time.Duration(10), // ondebug comment it
	}

	return c, nil
}
// change ResponseBody arr of openstack,
func (c *Client) DoReuest(api *request.OpenstackAPI) (error) {
	var request *http.Request
	var err error
	if(api.RequestBodyRequire == true){
		request, err = http.NewRequest(api.Method, api.Endpoint+api.Path, bytes.NewBuffer(api.RequestBody))
	}else{
		request, err = http.NewRequest(api.Method, api.Endpoint+api.Path, nil)
	}

	if(api.RequestParameterRequire == true){
		q := request.URL.Query()
		for k, v := range api.RequestParameter {
			q.Add(k,v)
		}
		request.URL.RawQuery = q.Encode()
	}
	if (err != nil ){
		return errors.New("bad request to "+request.URL.Path+" fail to normalized input")
	}
	for k, v := range api.RequestHeader {
		request.Header.Add(k,v)
	}
	// it may keep tls in section
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	} else {
		err = errors.New(fmt.Sprintf("RequestBody to "+request.URL.Path+" Respond status code %d", resp.StatusCode))
		return err
	}
	api.ResponseHeader = resp.Header
	api.ResponseBody, err = ioutil.ReadAll(resp.Body)
	return  err
}
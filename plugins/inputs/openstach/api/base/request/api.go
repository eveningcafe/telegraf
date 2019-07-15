package request

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// request type struct depent on format data in, respone is json data in responebody

const (
	XOpenStackNovaAPIv2Version  = "2.53"
)
type OpenstackAPI struct {
	Endpoint string
	Path     string
	Method   string
	HeaderRequest   map[string]string
	HeaderResponse  map[string][]string
	Request  []byte
	RequestBodyRequire bool // default no requestbody
	Response []byte
}
// change Response arr of openstack,
func (o *OpenstackAPI) DoReuest() (error) {
	var request *http.Request
	var err error
	httpClient := &http.Client{}
	if(o.RequestBodyRequire == true){
		request, err = http.NewRequest(o.Method, o.Endpoint+o.Path, bytes.NewBuffer(o.Request))
	}else{
		request, err = http.NewRequest(o.Method, o.Endpoint+o.Path, nil)
	}
	for k, v := range o.HeaderRequest {
		request.Header.Add(k,v)
	}
	resp, err := httpClient.Do(request)
	o.HeaderResponse = resp.Header
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	} else {
		err = errors.New("Request to "+request.URL.Path+"Respond status code "+ string(resp.StatusCode))
	}

	o.Response, err = ioutil.ReadAll(resp.Body)
	return  err
}
//type ListGroupAPI interface {
//	Path     string
//	Method   string
//	Header   map[string]string
//	Request  interface{}
//	Response interface{}
//	declare()
//}



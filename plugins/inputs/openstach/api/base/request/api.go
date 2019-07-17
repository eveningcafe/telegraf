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
	Endpoint           string
	Path               string
	Method             string
	RequestHeader      map[string]string
	ResponseHeader     map[string][]string
	RequestParameter   map[string]string
	RequestParameterRequire bool // default no requestparameter
	RequestBody        []byte
	RequestBodyRequire bool // default no requestbody
	ResponseBody       []byte
	
}
// change ResponseBody arr of openstack,
func (o *OpenstackAPI) DoReuest() (error) {
	var request *http.Request
	var err error
	httpClient := &http.Client{}
	if(o.RequestBodyRequire == true){
		request, err = http.NewRequest(o.Method, o.Endpoint+o.Path, bytes.NewBuffer(o.RequestBody))
	}else{
		request, err = http.NewRequest(o.Method, o.Endpoint+o.Path, nil)
	}

	if(o.RequestParameterRequire == true){
		q := request.URL.Query()
		for k, v := range o.RequestParameter {
			q.Add(k,v)
		}
		request.URL.RawQuery = q.Encode()
	}
	if (err != nil ){
		return errors.New("bad request to "+request.URL.Path+" fail to normalized input")
	}
	for k, v := range o.RequestHeader {
		request.Header.Add(k,v)
	}
	resp, err := httpClient.Do(request)
	o.ResponseHeader = resp.Header
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	} else {
		err = errors.New("RequestBody to "+request.URL.Path+"Respond status code "+ string(resp.StatusCode))
		return err
	}
	
	o.ResponseBody, err = ioutil.ReadAll(resp.Body)
	return  err
}




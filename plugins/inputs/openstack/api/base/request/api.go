package request

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

package groups


type ListGroupRequest struct {
}

type ListGroupResponse struct {
	Links struct {
		Self     string      `json:"self"`
		Previous interface{} `json:"previous"`
		Next     interface{} `json:"next"`
	} `json:"links"`
	Groups []struct {
		Description string `json:"description"`
		DomainID    string `json:"domain_id"`
		ID          string `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		Name string `json:"name"`
	} `json:"groups"`
}

type ListGroupAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  ListGroupRequest
	Response ListGroupResponse
}

// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-services-detail#list-services
func declareListGroup(token string) *ListGroupAPI{
	a:= new(ListGroupAPI)
	a.Path = "/groups"
	a.Method = "GET"
	a.Header = map[string]string{
		"Content-Type": "application/json",
		"X-Auth-Token": token,
	}
	return a
}
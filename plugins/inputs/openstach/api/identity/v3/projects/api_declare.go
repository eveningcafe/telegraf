package projects

type ListProjectRequest struct {
}

type ListProjectResponse struct {
	Links struct {
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Self     string      `json:"self"`
	} `json:"links"`
	Projects []struct {
		IsDomain    bool        `json:"is_domain"`
		DomainID    string      `json:"domain_id"`
		Enabled     bool        `json:"enabled"`
		ID          string      `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		Name     string        `json:"name"`
		ParentID interface{}   `json:"parent_id"`
		Tags     []interface{} `json:"tags"`
	} `json:"projects"`
}

type ListProjectAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  ListProjectRequest
	Response ListProjectResponse
}

// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-projects-detail
func declareListProject(token string) *ListProjectAPI{
	a:= new(ListProjectAPI)
	a.Path = "/projects"
	a.Method = "GET"
	a.Header = map[string]string{
		"Content-Type": "application/json",
		"X-Auth-Token": token,
	}
	return a
}




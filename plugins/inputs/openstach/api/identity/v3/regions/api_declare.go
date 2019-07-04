package regions

type ListRegionResponse struct {
	Links struct {
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
		Self     string      `json:"self"`
	} `json:"links"`
	Regions []struct {
		Description string `json:"description"`
		ID          string `json:"id"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		ParentRegionID interface{} `json:"parent_region_id"`
	} `json:"regions"`
}
type ListRegionRequest struct {
}

type ListRegionAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  ListRegionRequest
	Response ListRegionResponse
}

// https://developer.openstack.org/api-ref/identity/v3/?expanded=list-projects-detail
func declareListRegion(token string) *ListRegionAPI {
	a := new(ListRegionAPI)
	a.Path = "/regions"
	a.Method = "GET"
	a.Header = map[string]string{
		"Content-Type": "application/json",
		"X-Auth-Token": token,
	}
	return a
}

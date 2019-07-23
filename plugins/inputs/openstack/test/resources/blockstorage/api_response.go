package blockstorage

type ServiceListResponseBody struct {
	JsonBody string
}
func (r *ServiceListResponseBody) NewBody( )  {
	r.JsonBody = `

`
}
package authenticator

import (
	"encoding/json"
	"net/textproto"
)


type AuthOption struct {
	AuthURL           string
	ProjectDomainId string
	UserDomainId    string
	Username          string
	Password          string
	Project_name      string
}


type ProviderClient struct {
	Token       string
	UserID      string
	UserName    string
	ProjectID   string
	ProjectName string
	Catalog     []Catalog
	authURL     string
}


func AuthenticatedClient(options AuthOption) (*ProviderClient, error){
	api, err := declareCreateToken(options.AuthURL,options.Username,options.Password,options.Project_name,options.UserDomainId,options.ProjectDomainId)
	err = api.DoReuest()
	result := CreateTokenResponse{}
	err = json.Unmarshal([]byte(api.ResponseBody), &result)
	p := ProviderClient{
		UserName: result.Token.User.Name,
		UserID: result.Token.User.ID,
		ProjectName: result.Token.Project.Name,
		ProjectID: result.Token.Project.ID,
		authURL: options.AuthURL,
		Catalog: result.Token.Catalog,
		Token: textproto.MIMEHeader(api.ResponseHeader).Get("X-Subject-Token"),
	}
	return &p, err
}
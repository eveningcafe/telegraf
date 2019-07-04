package authenticator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (p *ProviderClient) GetCatalog() {

	//client := &http.Client{
	//}
	//request, err := http.NewRequest("POST", p.authURL+"/auth/tokens",nil)
	//request.Header.Set("Content-Type", "application/json")
	//resp, err := client.Do(request)
	//defer resp.Body.Close()
	//
	//if err != nil {
	//	panic(err.Error())
	//}
	//if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	//	fmt.Println("Authentication successful ")
	//} else {
	//	err := errors.New("Respond status code "+ string(resp.StatusCode))
	//	panic(err.Error())
	//}
}


func AuthenticatedClient(options AuthOption) (*ProviderClient, error){
	api := declareCreateToken(options.Username,options.Password,options.Project_name,options.UserDomainId,options.ProjectDomainId)

	jsonBody, err := json.Marshal(api.Request)

	if err != nil {
		panic(err.Error())
	}

	client := &http.Client{}
	request, err := http.NewRequest(api.Method, options.AuthURL+api.Path, bytes.NewBuffer(jsonBody))
	for k, v := range api.Header {
		request.Header.Add(k,v)
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("Authentication successful ")
	} else {
		err := errors.New("Respond status code "+ string(resp.StatusCode))
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal([]byte(body), &api.Response)
	p := ProviderClient{
		UserName: api.Response.Token.User.Name,
		UserID: api.Response.Token.User.ID,
		ProjectName: api.Response.Token.Project.Name,
		ProjectID: api.Response.Token.Project.ID,
		authURL: options.AuthURL,
		Catalog: api.Response.Token.Catalog,
		Token: resp.Header.Get("X-Subject-Token"),
	}

	return &p, err
}
//func main(){
//	provider, err := AuthenticatedClient(AuthOption{
//		authURL:           "http://controller:5000/v3",
//		method:            "password",
//		project_domain_id: "default",
//		user_domain_id:    "default",
//		username:          "admin",
//		password:          "Welcome123",
//		project_name:      "admin",
//	} )
//	fmt.Println(*provider)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
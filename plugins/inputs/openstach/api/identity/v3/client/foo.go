package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
type AuthOption struct {
	authURL string
	method string
	project_domain_id string
	user_domain_id string
	username string
	password string
	project_name string
	
}

type ResponeGetToken struct {
	Token struct {
		Methods []string `json:"methods"`
		User    struct {
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID                string      `json:"id"`
			Name              string      `json:"name"`
			PasswordExpiresAt interface{} `json:"password_expires_at"`
		} `json:"user"`
		AuditIds  []string  `json:"audit_ids"`
		ExpiresAt time.Time `json:"expires_at"`
		IssuedAt  time.Time `json:"issued_at"`
		Project   struct {
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		IsDomain bool `json:"is_domain"`
		Roles    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
		Catalog Catalog
	}
}

type Catalog struct {
	Endpoints []struct {
		ID        string `json:"id"`;
		Interface string `json:"interface"`;
		RegionID  string `json:"region_id"`;
		URL       string `json:"url"`;
		Region    string `json:"region"`
	} `json:"endpoints"`;
	ID   string `json:"id"`;
	Type string `json:"type"`;
	Name string `json:"name"`
}
type ProviderClient struct {
	Token       string
	UserID      string
	UserName    string
	ProjectID   string
	ProjectName string
	Catalog Catalog
}

func AuthenticatedClient(options AuthOption) (*ProviderClient, error){
	auth := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{options.method},
				"password": map[string]interface{}{
					"user": map[string]interface{}{
						"name": options.username,
						"domain": map[string]interface{}{
							"id": options.user_domain_id,
						},
						"password": options.password,
					},
				},
			},
			"scope": map[string]interface{}{
				"project": map[string]interface{}{
					"domain": map[string]interface{}{
						"id": options.project_domain_id,
					},
					"name": options.project_name,
				},
			},
		},
	}

	jsonData, err := json.Marshal(auth)

	//fmt.Println(string(jsonData))
	//log.Fatalln(err)

	if err != nil {
		panic(err.Error())
	}

	client := &http.Client{
	}
	request, err := http.NewRequest("POST", options.authURL+"/auth/tokens", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
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

	//var result map[string]interface{}

	//json.NewDecoder(resp.Body).Decode(&result)

	//log.Println(result)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var respData ResponeGetToken
	err = json.Unmarshal([]byte(body), &respData)
	p := ProviderClient{
		UserName: respData.Token.User.Name,
		UserID: respData.Token.User.ID,
		ProjectName: respData.Token.Project.Name,
		ProjectID: respData.Token.Project.ID,
		Token: resp.Header.Get("X-Subject-Token"),
		Catalog: respData.Token.Catalog,
	}
	fmt.Println(respData)

	return &p, err
}
func main(){
	provider, err := AuthenticatedClient(AuthOption{
		authURL:           "http://controller:5000/v3",
		method:            "password",
		project_domain_id: "default",
		user_domain_id:    "default",
		username:          "admin",
		password:          "Welcome123",
		project_name:      "admin",
	} )
	fmt.Println(*provider)
	if err != nil {
		log.Fatalln(err)
	}
}
package authenticator

import "time"

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
type CreateTokenRequest struct {
	Auth struct {
		Identity struct {
			Methods  []string `json:"methods"`
			Password struct {
				User struct {
					Name   string `json:"name"`
					Domain struct {
						ID string `json:"id"`
					} `json:"domain"`
					Password string `json:"password"`
				} `json:"user"`
			} `json:"password"`
		} `json:"identity"`
		Scope struct {
			Project struct {
				Name   string `json:"name"`
				Domain struct {
					ID string `json:"id"`
				} `json:"domain"`
			} `json:"project"`
		} `json:"scope"`
	} `json:"auth"`
}
type CreateTokenResponse struct {
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
		Catalog []Catalog
	}
}

type CreateTokenAPI struct {
	Path     string
	Method   string
	Header   map[string]string
	Request  CreateTokenRequest
	Response CreateTokenResponse
}

func declareGETToken(userName string, password string, project string, userDomainID string, projectDomainID string) *CreateTokenAPI {
	a := new(CreateTokenAPI)
	a.Path = "/auth/tokens"
	a.Method = "POST"
	a.Header = map[string]string{
		"Content-Type": "application/json",
	}
	a.Request = CreateTokenRequest{}
	a.Request.Auth.Identity.Methods                  = []string{"password"}
	a.Request.Auth.Identity.Password.User.Password   = password
	a.Request.Auth.Identity.Password.User.Domain.ID  = userDomainID
	a.Request.Auth.Identity.Password.User.Name       = userName
	a.Request.Auth.Scope.Project.Name                = project
	a.Request.Auth.Scope.Project.Domain.ID           = projectDomainID
	return a
}

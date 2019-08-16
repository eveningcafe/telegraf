package vcs_mgr_ceph

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type VCSCeph struct {
	MgrPath   string
	Username  string
	Password  string
	SessionID string
}

var VCSCephConfig = `
  ## ceph mgr dashboard endpoint
  mgr_path = "https://192.168.28.12:8443" 
  ## username of mgr dasboard
  username = "admin"
  password = "123qweA@"
`

func (c *VCSCeph) SampleConfig() string {
	return VCSCephConfig
}

func (c *VCSCeph) Description() string {
	return "crawl ceph dashboard to get monitor metric "
}

func AuthenticatedClient(usename string, password string) ( string, error) {
	//api, err := declareCreateToken(options.AuthURL,options.Username,options.Password,options.Project_name,options.UserDomainId,options.ProjectDomainId)
	//httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		TLSClientConfig: options.TlsCfg,
	//	},
	//	//Timeout: time.Duration(5), // ondebug comment it
	//}
	//var request *http.Request
	//if !strings.Contains(api.Endpoint , "v3") {
	//	api.Endpoint = api.Endpoint+"/v3"
	//}
	//request, err = http.NewRequest(api.Method, api.Endpoint+api.Path, bytes.NewBuffer(api.RequestBody))
	//if (err != nil) {
	//	return nil, errors.New("bad request to " + request.URL.Path + " fail to normalized input")
	//}
	//for k, v := range api.RequestHeader {
	//	request.Header.Add(k, v)
	//}
	//resp, err := httpClient.Do(request)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	//} else {
	//	err = errors.New(fmt.Sprintf("RequestBody to " + request.URL.Path + "Respond status code %d", resp.StatusCode))
	//	return nil, err
	//}
	//api.ResponseHeader = resp.Header
	//api.ResponseBody, err = ioutil.ReadAll(resp.Body)
	//result := CreateTokenResponse{}
	//err = json.Unmarshal([]byte(api.ResponseBody), &result)
	//p := ProviderClient{
	//	UserName: result.Token.User.Name,
	//	UserID: result.Token.User.ID,
	//	ProjectName: result.Token.Project.Name,
	//	ProjectID: result.Token.Project.ID,
	//	authURL: options.AuthURL,
	//	Catalog: result.Token.Catalog,
	//	Token: textproto.MIMEHeader(api.ResponseHeader).Get("X-Subject-Token"),
	//	TlsCfg: options.TlsCfg,
	//}
	return "", nil
}
func (c *VCSCeph) Init() (string, error) {
	sessionid, err := AuthenticatedClient("admin", "123qweA@")

	return sessionid, err
}
func (c *VCSCeph) Gather(acc telegraf.Accumulator) error {
	sessionId, err := c.Init()
	fmt.Println(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	inputs.Add("trig", func() telegraf.Input { return &VCSCeph{} })
}

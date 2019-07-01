package main

import (
	"bytes"
	"fmt"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/internal/tls"
	"github.com/influxdata/telegraf/plugins/parsers"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseURL string
	token   string
	servicesList []string
	endPointsList []string
}

type HTTP struct {
	URL    string
	Method string
	Body    *bytes.Buffer

	Headers            map[string]string
	InsecureSkipVerify bool

	// HTTP Basic Auth Credentials
	tls.ClientConfig

	Timeout internal.Duration

	client *http.Client

	// The parser will automatically be set by Telegraf core code because
	// this plugin implements the ParserInput interface (i.e. the SetParser method)
	parser parsers.Parser

}

func (h *HTTP) Init() error {
	tlsCfg, err := h.ClientConfig.TLSConfig()
	if err != nil {
		return err
	}

	h.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
		Timeout: h.Timeout.Duration,
	}
	return nil
}


func InitClient(identityEndpoint string, domain string, project string, username string, password string, tsl tls.ClientConfig) (*Client, error){
	header := make(map[string]string)
	header["Content-Type"]= "application/json"

	var authData = []byte(`{
                  "auth": {
                     "identity": {
                         "methods": [
                            "password"
                          ],
                         "password": {
                            "user": {
                                "name": `+username+`,
                                "domain": {"id": `+domain+`},
                                "password": ` + password+`
                             }
                          }
                      },
                     "scope": {
                         "project": {
                            "name": `+project+`,
                            "domain": {"id": `+domain+`}
                            }
                       }
                  }
               }`)


	h := HTTP{
		URL: identityEndpoint + "/auth/tokens",
		Method: "GET",
		Headers: header,
		InsecureSkipVerify: true,
		Body: bytes.NewBuffer(authData),


	}
	h.Init()

	request, err := http.NewRequest(h.Method, h.URL, h.Body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(request)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received status code %d (%s), expected %d (%s)",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
			http.StatusOK,
			http.StatusText(http.StatusOK))
	}

	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(b)

	return nil, nil
}

func main() {

}
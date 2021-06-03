package servicenow

import (
	"net/http"
)

type ClientAuth struct {
	Instance string
	User     string
	Password string
}

type Client struct {
	*http.Client
	*ClientAuth
}

func initServiceNowClient(a *ClientAuth) (*Client, error) {

	sn := new(Client)
	sn.Client = &http.Client{}
	sn.ClientAuth = a

	return sn, nil
}

func (c *Client) setBasicAuth(r *http.Request) {
	r.SetBasicAuth(c.ClientAuth.User, c.ClientAuth.Password)
}

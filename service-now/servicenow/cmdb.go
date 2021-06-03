package servicenow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetCMDBInstance(class, sysID string) error {

	u := fmt.Sprintf("%s/now/cmdb/instance/%s/%s", c.ClientAuth.Instance, class, sysID)
	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	r.Header.Set("Accept", "application/json")
	c.setBasicAuth(r)

	resp, err := c.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
}

func (c *Client) CreateCMDBInstance(class string, payload []byte) (*CreateCMDBInstance_Response, error) {

	u := fmt.Sprintf("%s/now/cmdb/instance/%s", c.ClientAuth.Instance, class)
	r, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	r.Header.Set("Content-Type", "application/json")
	c.setBasicAuth(r)

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	out := new(CreateCMDBInstance_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) UpdateCMDBInstance(class, sysID string, payload []byte) (*UpdateCMDBInstance_Response, error) {

	u := fmt.Sprintf("%s/now/cmdb/instance/%s/%s", c.ClientAuth.Instance, class, sysID)
	r, err := http.NewRequest(http.MethodPut, u, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	r.Header.Set("Content-Type", "application/json")
	c.setBasicAuth(r)

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	out := new(UpdateCMDBInstance_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

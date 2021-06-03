package servicenow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetIncident(sysID string) (*GetIncident_Response, error) {

	u := fmt.Sprintf("%s/api/now/table/incident/%s", c.ClientAuth.Instance, sysID)
	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	r.Header.Set("Accept", "application/json")
	c.setBasicAuth(r)

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	out := new(GetIncident_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) CreateIncident() (*CreateIncident_Response, error) {

	u := fmt.Sprintf("%s/api/now/table/incident", c.ClientAuth.Instance)
	r, err := http.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	r.Header.Set("Accept", "application/json")
	c.setBasicAuth(r)

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	out := new(CreateIncident_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) UpdateIncident(sysID string, payload []byte) (*UpdateIncident_Response, error) {

	u := fmt.Sprintf("%s/api/now/table/incident/%s?sysparm_exclude_reference_link=true", c.ClientAuth.Instance, sysID)
	r, err := http.NewRequest(http.MethodPut, u, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	r.Header.Set("Accept", "application/json")
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

	out := new(UpdateIncident_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

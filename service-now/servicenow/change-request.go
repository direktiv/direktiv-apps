package servicenow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetChangeRequest(sysID string) (*ChangeRequest_Response, error) {

	u := fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", c.ClientAuth.Instance, sysID)
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

	out := new(ChangeRequest_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) ApproveChangeRequest(sysID string) (*ChangeRequest_Response, error) {

	u := fmt.Sprintf("%s/api/sn_chg_rest/change/%s/approvals", c.ClientAuth.Instance, sysID)
	r, err := http.NewRequest(http.MethodPatch, u, strings.NewReader(fmt.Sprintf(`{"state": "approved"}`)))
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

	out := new(ChangeRequest_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) RejectChangeRequest(sysID string) (*ChangeRequest_Response, error) {

	u := fmt.Sprintf("%s/api/sn_chg_rest/change/%s/approvals", c.ClientAuth.Instance, sysID)
	r, err := http.NewRequest(http.MethodPatch, u, strings.NewReader(fmt.Sprintf(`{"state": "rejected"}`)))
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

	out := new(ChangeRequest_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) CreateNormalChangeRequest() (*ChangeRequest_Response, error) {

	u := fmt.Sprintf("%s/api/sn_chg_rest/change/normal", c.ClientAuth.Instance)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	out := new(ChangeRequest_Response)
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

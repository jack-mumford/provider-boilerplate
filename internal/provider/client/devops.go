package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetDevs - Returns list of devs (no auth required)
func (c *Client) GetDevs() ([]Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	items := []Dev{}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetOps - Returns list of ops (no auth required)
func (c *Client) GetOps() ([]Ops, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/op", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	items := []Ops{}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetDevOps - Returns list of devops (no auth required)
func (c *Client) GetDevOps() ([]DevOps, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devops", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	items := []DevOps{}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

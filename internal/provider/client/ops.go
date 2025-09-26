package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateOps(ops Ops) (*Ops, error) {
	rb, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ops", c.endpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newOps := Ops{}
	err = json.Unmarshal(body, &newOps)
	if err != nil {
		return nil, err
	}

	return &newOps, nil
}

func (c *Client) UpdateOps(opsID string, ops Ops) (*Ops, error) {
	rb, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/ops/id/%s", c.endpoint, opsID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedOps := Ops{}
	err = json.Unmarshal(body, &updatedOps)
	if err != nil {
		return nil, err
	}

	return &updatedOps, nil
}

func (c *Client) DeleteOps(opsID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/ops/%s", c.endpoint, opsID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

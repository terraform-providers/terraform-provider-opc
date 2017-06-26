package storage

import (
	"net/http"
)

func (c *Client) createResource(name string, requestHeaders interface{}) error {
	_, err := c.executeRequest("PUT", name, requestHeaders)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) updateResource(name string, requestHeaders interface{}) error {
	_, err := c.executeRequest("PUT", name, requestHeaders)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getResource(name string, responseBody interface{}) (*http.Response, error) {
	rsp, err := c.executeRequest("GET", name, nil)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) deleteResource(name string) error {
	_, err := c.executeRequest("DELETE", name, nil)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

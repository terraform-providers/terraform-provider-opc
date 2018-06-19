package storage

import (
	"io"
	"net/http"
)

func (c *Client) createResource(name string, requestHeaders interface{}) error {
	return c.createResourceBody(name, requestHeaders, nil)
}

func (c *Client) createResourceBody(name string, requestHeaders interface{}, body io.ReadSeeker) error {
	_, err := c.executeRequestBody("PUT", name, requestHeaders, body)
	return err
}

func (c *Client) updateResource(name string, requestHeaders interface{}) error {
	_, err := c.executeRequest("PUT", name, requestHeaders)
	return err
}

func (c *Client) getResource(name string, responseBody interface{}) (*http.Response, error) {
	return c.getResourceHeaders(name, responseBody, nil)
}

func (c *Client) getResourceHeaders(
	name string,
	responseBody interface{},
	requestHeaders interface{},
) (*http.Response, error) {
	rsp, err := c.executeRequest("GET", name, requestHeaders)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) deleteResource(name string) error {
	_, err := c.executeRequest("DELETE", name, nil)
	return err
}

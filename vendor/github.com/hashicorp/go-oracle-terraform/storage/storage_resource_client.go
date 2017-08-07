package storage

import (
	"io"
	"net/http"
)

func (c *StorageClient) createResource(name string, requestHeaders interface{}) error {
	return c.createResourceBody(name, requestHeaders, nil)
}

func (c *StorageClient) createResourceBody(name string, requestHeaders interface{}, body io.ReadSeeker) error {
	_, err := c.executeRequestBody("PUT", name, requestHeaders, body)
	if err != nil {
		return err
	}

	return nil
}

func (c *StorageClient) updateResource(name string, requestHeaders interface{}) error {
	_, err := c.executeRequest("PUT", name, requestHeaders)
	if err != nil {
		return err
	}

	return nil
}

func (c *StorageClient) getResource(name string, responseBody interface{}) (*http.Response, error) {
	return c.getResourceHeaders(name, responseBody, nil)
}

func (c *StorageClient) getResourceHeaders(
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

func (c *StorageClient) deleteResource(name string) error {
	_, err := c.executeRequest("DELETE", name, nil)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

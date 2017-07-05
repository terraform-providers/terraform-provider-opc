package storage

import (
	"net/http"
)

func (c *StorageClient) createResource(name string, requestHeaders interface{}) error {
	_, err := c.executeRequest("PUT", name, requestHeaders)
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
	rsp, err := c.executeRequest("GET", name, nil)
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

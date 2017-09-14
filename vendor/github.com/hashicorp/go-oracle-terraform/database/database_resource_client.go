package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/mitchellh/mapstructure"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*DatabaseClient
	ContainerPath    string
	ResourceRootPath string
}

func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath), requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("PUT", c.getObjectPath(c.ResourceRootPath, name), requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) getResource(name string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

// This is only used for deleting service instances. DELETE requests have a `nil` body.
func (c *ResourceClient) deleteResource(name string, backups bool) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}

	// Set deleteBackup
	if backups {
		objectPath = fmt.Sprintf("%s?deleteBackup=true", objectPath)
	}

	_, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		if v, ok := err.(*opc.OracleError); ok {
			if v.StatusCode == 404 {
				// Object can't be found, doesn't exist, no error
				return nil
			}
			return fmt.Errorf("Error on delete (%d): %s", v.StatusCode, v.Message)
		}

		// Otherwise, something went wrong.
		return err
	}

	// No errors and no response body to write
	return nil
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	c.client.DebugLogString(fmt.Sprintf("HTTP Resp (%d): %s", resp.StatusCode, buf.String()))
	// JSON decode response into interface
	var tmp interface{}
	dcd := json.NewDecoder(buf)
	if err := dcd.Decode(&tmp); err != nil {
		return fmt.Errorf("%+v", resp)
		return err
	}

	// Use mapstructure to weakly decode into the resulting interface
	msdcd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           iface,
		TagName:          "json",
	})
	if err != nil {
		return err
	}

	if err := msdcd.Decode(tmp); err != nil {
		return err
	}
	return nil
}

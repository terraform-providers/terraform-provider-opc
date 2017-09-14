package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// The UtilityClient which extends the UtilityResourceClient.
// This is purely because utility resources (SSH Keys + Access Rules) include the service
// instance name in the URL path for managing these resources, so we cannot use the same
// resource client that the service instance uses. We're still using the same OPC client, just
// abstracting the path helper functions to make life a little easier.
type UtilityClient struct {
	UtilityResourceClient
}

// UtilityResourceClient is a client to manage resources on an already created service instance
type UtilityResourceClient struct {
	*DatabaseClient
	ContainerPath     string
	ResourceRootPath  string
	ServiceInstanceID string
}

func (c *UtilityResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath), requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *UtilityResourceClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("PUT", c.getObjectPath(c.ResourceRootPath, name), requestBody)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *UtilityResourceClient) getResource(name string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.getContainerPath(c.ContainerPath)
	}

	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *UtilityResourceClient) deleteResource(name string, body interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("DELETE", objectPath, body)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

func (c *UtilityResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
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

func (c *UtilityResourceClient) getContainerPath(root string) string {
	// /paas/api/v1.1/instancemgmt/{identityDomainId}/services/dbaas/instances/{serviceId}/accessrules
	return fmt.Sprintf(root, *c.client.IdentityDomain, c.ServiceInstanceID)
}

func (c *UtilityResourceClient) getObjectPath(root, name string) string {
	// /paas/api/v1.1/instancemgmt/{identityDomainId}/services/dbaas/instances/{serviceId}/accessrules/{ruleName}
	return fmt.Sprintf(root, *c.client.IdentityDomain, c.ServiceInstanceID, name)
}

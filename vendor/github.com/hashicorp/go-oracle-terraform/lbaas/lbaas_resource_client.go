package lbaas

import (
	"fmt"
)

/*
 * The LBaaSResourceClient is the general client used for the majority of the Load Balancer
 * Service child resources (Listeners, Origin Servier Pools and Policies) which have the common URI
 * format https://{api_endpoint}/{lb_name}/{lb_region}/{resource_type}/{resource_name}?{projection}
 *
 * For SSL Certificates use the SSLCertificateClient
 * For the Load Balancer Service Instance use the LoadBalancerResourceClient
 */

// LBaaSResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type LBaaSResourceClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
	Projection       string
	Accept           string
	ContentType      string
}

// executes the Create requests to the LBaaS API
func (c *LBaaSResourceClient) createResource(lbRegion, lbName string, requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath, lbRegion, lbName), c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the LBaaS API
func (c *LBaaSResourceClient) updateResource(lbRegion, lbName, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequest("PUT", objectPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the LBaaS API specific to updating the Origin Server Pool
// which has a different update style using POST + with an PATCH Method override
func (c *LBaaSResourceClient) updateOriginServerPool(lbRegion, lbName, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequestWithMethodOverride("POST", "PATCH", objectPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Get requests to the LBaaS API
func (c *LBaaSResourceClient) getResource(lbRegion, lbName, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	queryParams := ""
	if c.Projection != "" {
		queryParams = fmt.Sprintf("?projection=%s" + c.Projection)
	}
	resp, err := c.executeRequest("GET", objectPath+queryParams, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Delete requests to the LBaaS API
func (c *LBaaSResourceClient) deleteResource(lbRegion, lbName, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequest("DELETE", objectPath, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

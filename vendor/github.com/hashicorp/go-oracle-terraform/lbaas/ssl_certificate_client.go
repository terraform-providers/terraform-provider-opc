package lbaas

import (
	"fmt"
	"strings"
)

const (
	sslCertificateContainerPath = "/certs"
	sslCertificaetResourcePath  = "/certs/%s"
)

// Supported ContentTypes for the SSL Certificate API requests
const ContentTypeServerCertificateJSON = "application/vnd.com.oracle.oracloud.lbaas.ServerCertificate+json"
const ContentTypeTrustedCertificateJSON = "application/vnd.com.oracle.oracloud.lbaas.TrustedCertificate+json"

// SSLCertificateClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type SSLCertificateClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
	Accept           string
	ContentType      string
}

// SSLCertificateClient returns an ServiceInstanceClient which is used to access the
// Load Balancer API
func (c *Client) SSLCertificateClient() *SSLCertificateClient {
	return &SSLCertificateClient{
		Client:           c,
		ContainerPath:    sslCertificateContainerPath,
		ResourceRootPath: sslCertificaetResourcePath,
		Accept: strings.Join([]string{
			ContentTypeServerCertificateJSON,
			ContentTypeTrustedCertificateJSON,
		}, ","),
		// ContentType cannot be generally set for the SSLCertificateCleint, instead it is set on each
		// Create request based on the Type of the Certificate being created.
	}
}

func (c *SSLCertificateClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, name)
}

// executes the Create requests to the Load Balancer API
func (c *SSLCertificateClient) createResource(requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.ContainerPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Get requests to the Load Balancer API
func (c *SSLCertificateClient) getResource(name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequest("GET", objectPath, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the Load Balancer API
func (c *SSLCertificateClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequestWithMethodOverride("POST", "PATCH", objectPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Delete requests to the Load Balancer API
func (c *SSLCertificateClient) deleteResource(name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequest("DELETE", objectPath, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

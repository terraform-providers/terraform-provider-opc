package storage

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const STR_ACCOUNT = "/Storage-%s"
const STR_USERNAME = "/Storage-%s:%s"
const AUTH_HEADER = "X-Auth-Token"
const STR_QUALIFIED_NAME = "%s%s/%s"

// Client represents an authenticated compute client, with compute credentials and an api client.
type StorageClient struct {
	client      *client.Client
	authToken   *string
	tokenIssued time.Time
}

func NewStorageClient(c *opc.Config) (*StorageClient, error) {
	storageClient := &StorageClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(fmt.Sprintf("https://%s.storage.oraclecloud.com", *client.IdentityDomain))
	if err != nil {
		return nil, err
	}
	client.APIEndpoint = endpoint
	storageClient.client = client

	if err := storageClient.getAuthenticationToken(); err != nil {
		return nil, err
	}

	return storageClient, nil
}

func (c *StorageClient) executeRequest(method, path string, headers interface{}) (*http.Response, error) {
	req, err := c.client.BuildRequest(method, path, nil)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers.(map[string]string) {
			req.Header.Add(k, v)
		}
	}

	// If we have an authentication token, let's authenticate, refreshing cookie if need be
	if c.authToken != nil {
		if time.Since(c.tokenIssued).Minutes() > 25 {
			if err := c.getAuthenticationToken(); err != nil {
				return nil, err
			}
		}
		req.Header.Add(AUTH_HEADER, *c.authToken)
	}

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *StorageClient) getUserName() string {
	return fmt.Sprintf(STR_USERNAME, *c.client.IdentityDomain, *c.client.UserName)
}

func (c *StorageClient) getAccount() string {
	return fmt.Sprintf(STR_ACCOUNT, *c.client.IdentityDomain)
}

// GetQualifiedName returns the fully-qualified name of a storage object, e.g. /v1/{account}/{name}
func (c *StorageClient) getQualifiedName(version string, name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "/Storage-") || strings.HasPrefix(name, "v1/") {
		return name
	}
	return fmt.Sprintf(STR_QUALIFIED_NAME, version, c.getAccount(), name)
}

// GetUnqualifiedName returns the unqualified name of a Storage object, e.g. the {name} part of /v1/{account}/{name}
func (c *StorageClient) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if !strings.Contains(name, "/") {
		return name
	}

	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[len(nameParts)-1:], "/")
}

func (c *StorageClient) unqualify(names ...*string) {
	for _, name := range names {
		*name = c.getUnqualifiedName(*name)
	}
}

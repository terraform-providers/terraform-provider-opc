package database

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const DB_ACCOUNT = "/Database-%s"
const DB_USERNAME = "/Database-%s:%s"
const AUTH_HEADER = "Authorization"
const TENANT_HEADER = "X-ID-TENANT-NAME"
const DB_QUALIFIED_NAME = "%s%s/%s"

// Client represents an authenticated database client, with compute credentials and an api client.
type DatabaseClient struct {
	client     *client.Client
	authHeader *string
}

func NewDatabaseClient(c *opc.Config) (*DatabaseClient, error) {
	databaseClient := &DatabaseClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	databaseClient.client = client

	databaseClient.authHeader = databaseClient.getAuthenticationHeader()

	return databaseClient, nil
}

func (c *DatabaseClient) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Req (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		// Debug the body for database services
		debugReqString = fmt.Sprintf("%s:\n %+v", debugReqString, string(reqBody))
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)

	// Set the authentication headers
	req.Header.Add(AUTH_HEADER, *c.authHeader)
	req.Header.Add(TENANT_HEADER, *c.client.IdentityDomain)
	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *DatabaseClient) getUserName() string {
	return fmt.Sprintf(DB_USERNAME, *c.client.IdentityDomain, *c.client.UserName)
}

func (c *DatabaseClient) getAccount() string {
	return fmt.Sprintf(DB_ACCOUNT, *c.client.IdentityDomain)
}

// GetQualifiedName returns the fully-qualified name of a database object, e.g. /v1/{account}/{name}
func (c *DatabaseClient) getQualifiedName(version string, name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "/Database-") || strings.HasPrefix(name, "v1/") {
		return name
	}
	return fmt.Sprintf(DB_QUALIFIED_NAME, version, c.getAccount(), name)
}

// GetUnqualifiedName returns the unqualified name of a Database object, e.g. the {name} part of /v1/{account}/{name}
func (c *DatabaseClient) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if !strings.Contains(name, "/") {
		return name
	}

	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[len(nameParts)-1:], "/")
}

func (c *DatabaseClient) unqualify(names ...*string) {
	for _, name := range names {
		*name = c.getUnqualifiedName(*name)
	}
}

func (c *DatabaseClient) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *DatabaseClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}

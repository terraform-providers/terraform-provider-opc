package storage

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const strAccount = "/Storage-%s"
const strUsername = "/Storage-%s:%s"
const authHeader = "X-Auth-Token"
const strQualifiedName = "%s%s/%s"
const apiVersion = "v1"

// Client represents an authenticated storage client, with storage credentials and an api client.
type Client struct {
	client      *client.Client
	authToken   *string
	tokenIssued time.Time
}

// NewStorageClient returns an authenticate storage client
func NewStorageClient(c *opc.Config) (*Client, error) {
	sClient := &Client{}
	opcClient, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	sClient.client = opcClient

	if err := sClient.getAuthenticationToken(); err != nil {
		return nil, err
	}

	return sClient, nil
}

// Execute a request with a nil body
func (c *Client) executeRequest(method, path string, headers interface{}) (*http.Response, error) {
	return c.executeRequestBody(method, path, headers, nil)
}

// Execute a request with a body supplied. The body can be nil for the request.
// Does not marshal the body into json to create the request
func (c *Client) executeRequestBody(method, path string, headers interface{}, body io.ReadSeeker) (*http.Response, error) {
	req, err := c.client.BuildNonJSONRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("%s (%s) %s", req.Method, req.URL, req.Proto)
	var debugHeaders []string

	if headers != nil {
		for k, v := range headers.(map[string]string) {
			debugHeaders = append(debugHeaders,
				fmt.Sprintf("%v: %v\n", strings.ToLower(k), v))

			req.Header.Add(k, v)
		}
		debugReqString = fmt.Sprintf("%s\n%s", debugReqString, debugHeaders)
	}

	if !strings.Contains(path, "/auth/") {
		c.client.DebugLogString(debugReqString)
	}

	// If we have an authentication token, let's authenticate, refreshing cookie if need be
	if c.authToken != nil {
		if time.Since(c.tokenIssued).Minutes() > 25 {
			if err = c.getAuthenticationToken(); err != nil {
				return nil, err
			}
		}
		req.Header.Add(authHeader, *c.authToken)
	}

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) getUserName() string {
	return fmt.Sprintf(strUsername, *c.client.IdentityDomain, *c.client.UserName)
}

func (c *Client) getAccount() string {
	return fmt.Sprintf(strAccount, *c.client.IdentityDomain)
}

// GetQualifiedName returns the fully-qualified name of a storage object, e.g. /v1/{account}/{name}
func (c *Client) getQualifiedName(name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "/Storage-") || strings.HasPrefix(name, apiVersion+"/") {
		return name
	}
	return fmt.Sprintf(strQualifiedName, apiVersion, c.getAccount(), name)
}

// GetUnqualifiedName returns the unqualified name of a Storage object, e.g. the {name} part of /v1/{account}/{name}
func (c *Client) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if !strings.Contains(name, "/") {
		return name
	}

	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[len(nameParts)-1:], "/")
}

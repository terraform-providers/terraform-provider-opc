package storage

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

const STR_ACCOUNT = "/Storage-%s"
const STR_USERNAME = "/Storage-%s:%s"
const DEFAULT_MAX_RETRIES = 1
const AUTH_HEADER = "X-Auth-Token"
const STR_QUALIFIED_NAME = "%s%s/%s"

// Client represents an authenticated compute client, with compute credentials and an api client.
type Client struct {
	identityDomain *string
	userName       *string
	password       *string
	apiEndpoint    *url.URL
	httpClient     *http.Client
	authToken      *string
	tokenIssued    time.Time
	maxRetries     *int
	logger         opc.Logger
	loglevel       opc.LogLevelType
}

func NewStorageClient(c *opc.Config) (*Client, error) {
	var err error
	// First create a client
	client := &Client{
		identityDomain: c.IdentityDomain,
		userName:       c.Username,
		password:       c.Password,
		apiEndpoint:    c.APIEndpoint,
		httpClient:     c.HTTPClient,
		maxRetries:     c.MaxRetries,
		loglevel:       c.LogLevel,
	}

	client.apiEndpoint, err = url.Parse(fmt.Sprintf("https://%s.storage.oraclecloud.com", *client.identityDomain))
	if err != nil {
		return nil, err
	}

	// Setup logger; defaults to stdout
	if c.Logger == nil {
		client.logger = opc.NewDefaultLogger()
	} else {
		client.logger = c.Logger
	}

	// If LogLevel was not set to something different,
	// double check for env var
	if c.LogLevel == 0 {
		client.loglevel = opc.LogLevel()
	}

	if err := client.getAuthenticationToken(); err != nil {
		return nil, err
	}

	// Default max retries if unset
	if c.MaxRetries == nil {
		client.maxRetries = opc.Int(DEFAULT_MAX_RETRIES)
	}

	// Protect against any nil http client
	if c.HTTPClient == nil {
		return nil, fmt.Errorf("No HTTP client specified in config")
	}

	return client, nil
}

func (c *Client) executeRequest(method, path string, headers interface{}) (*http.Response, error) {
	// Parse URL Path
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(method, c.formatURL(urlPath), nil)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers.(map[string]string) {
			req.Header.Add(k, v)
		}
	}

	// If we have an authentication toekn, let's authenticate, refreshing token if need be
	if c.authToken != nil {
		if time.Since(c.tokenIssued).Minutes() > 25 {
			if err := c.getAuthenticationToken(); err != nil {
				return nil, err
			}
		}
		req.Header.Add(AUTH_HEADER, *c.authToken)
	}

	// Execute request with supplied client
	resp, err := c.retryRequest(req)
	//resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}

	oracleErr := &opc.OracleError{
		StatusCode: resp.StatusCode,
	}
	// Even though the returned body will be in json form, it's undocumented what
	// fields are actually returned. Once we get documentation of the actual
	// error fields that are possible to be returned we can have stricter error types.
	if resp.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		oracleErr.Message = buf.String()
	}

	return nil, oracleErr
}

// Allow retrying the request until it either returns no error,
// or we exceed the number of max retries
func (c *Client) retryRequest(req *http.Request) (*http.Response, error) {
	// Double check maxRetries is not nil
	var retries int
	if c.maxRetries == nil {
		retries = DEFAULT_MAX_RETRIES
	} else {
		retries = *c.maxRetries
	}

	var statusCode int
	var errMessage string

	for i := 0; i < retries; i++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
			return resp, nil
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		errMessage = buf.String()
		statusCode = resp.StatusCode
		c.debugLogString(fmt.Sprintf("Encountered HTTP (%d) Error: %s", statusCode, errMessage))
		c.debugLogString(fmt.Sprintf("%d/%d retries left", i+1, retries))
	}

	oracleErr := &opc.OracleError{
		StatusCode: statusCode,
		Message:    errMessage,
	}

	// We ran out of retries to make, return the error and response
	return nil, oracleErr
}

func (c *Client) formatURL(path *url.URL) string {
	return c.apiEndpoint.ResolveReference(path).String()
}

func (c *Client) getUserName() string {
	return fmt.Sprintf(STR_USERNAME, *c.identityDomain, *c.userName)
}

func (c *Client) getAccount() string {
	return fmt.Sprintf(STR_ACCOUNT, *c.identityDomain)
}

// From compute_client
// GetObjectName returns the fully-qualified name of an OPC object, e.g. /identity-domain/user@email/{name}
func (c *Client) getQualifiedName(version string, name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "/Storage-") || strings.HasPrefix(name, "v1/") {
		return name
	}
	return fmt.Sprintf(STR_QUALIFIED_NAME, version, c.getAccount(), name)
}

// GetUnqualifiedName returns the unqualified name of an OPC object, e.g. the {name} part of /identity-domain/user@email/{name}
func (c *Client) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if strings.HasPrefix(name, "/oracle") {
		return name
	}
	if !strings.Contains(name, "/") {
		return name
	}

	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[3:], "/")
}

func (c *Client) unqualify(names ...*string) {
	for _, name := range names {
		*name = c.getUnqualifiedName(*name)
	}
}

func (c *Client) unqualifyUrl(url *string) {
	var validID = regexp.MustCompile(`(\/(Compute[^\/\s]+))(\/[^\/\s]+)(\/[^\/\s]+)`)
	name := validID.FindString(*url)
	*url = c.getUnqualifiedName(name)
}

// Retry function
func (c *Client) waitFor(description string, timeoutSeconds int, test func() (bool, error)) error {
	tick := time.Tick(1 * time.Second)

	for i := 0; i < timeoutSeconds; i++ {
		select {
		case <-tick:
			completed, err := test()
			c.debugLogString(fmt.Sprintf("Waiting for %s (%d/%ds)", description, i, timeoutSeconds))
			if err != nil || completed {
				return err
			}
		}
	}
	return fmt.Errorf("Timeout waiting for %s", description)
}

// Used to determine if the checked resource was found or not.
func WasNotFoundError(e error) bool {
	err, ok := e.(*opc.OracleError)
	if ok {
		return err.StatusCode == 404
	}
	return false
}

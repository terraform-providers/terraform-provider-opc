package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

const DEFAULT_MAX_RETRIES = 1

// Client represents an authenticated compute client, with compute credentials and an api client.
type Client struct {
	IdentityDomain *string
	UserName       *string
	Password       *string
	apiEndpoint    *url.URL
	httpClient     *http.Client
	authCookie     *http.Cookie
	cookieIssued   time.Time
	maxRetries     *int
	logger         opc.Logger
	loglevel       opc.LogLevelType
}

func NewClient(c *opc.Config) (*Client, error) {
	// First create a client
	client := &Client{
		IdentityDomain: c.IdentityDomain,
		UserName:       c.Username,
		Password:       c.Password,
		apiEndpoint:    c.APIEndpoint,
		httpClient:     c.HTTPClient,
		maxRetries:     c.MaxRetries,
		loglevel:       c.LogLevel,
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

// This builds an http request.
// After calling this you need to add the authentication. Header/Cookie/etc
// Then call ExecuteRequest and pass in the return value of this method
// It is split up to add additional authentication that is Oracle API dependent.
func (c *Client) BuildRequest(method, path string, body interface{}) (*http.Request, error) {
	// Parse URL Path
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Marshall request body
	var requestBody io.ReadSeeker
	var marshaled []byte
	if body != nil {
		marshaled, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	// Create request
	req, err := http.NewRequest(method, c.formatURL(urlPath), requestBody)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// This method executes the http.Request from the BuildRequest method.
// It is split up to add additional authentication that is Oracle API dependent.
func (c *Client) ExecuteRequest(req *http.Request) (*http.Response, error) {
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
		c.DebugLogString(fmt.Sprintf("Encountered HTTP (%d) Error: %s", statusCode, errMessage))
		c.DebugLogString(fmt.Sprintf("%d/%d retries left", i+1, retries))
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

// Retry function
func (c *Client) WaitFor(description string, timeoutSeconds int, test func() (bool, error)) error {
	tick := time.Tick(1 * time.Second)

	for i := 0; i < timeoutSeconds; i++ {
		select {
		case <-tick:
			completed, err := test()
			c.DebugLogString(fmt.Sprintf("Waiting for %s (%d/%ds)", description, i, timeoutSeconds))
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

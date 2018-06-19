package storage

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Header Constants
const (
	hContainerRead              = "X-Container-Read"
	hContainerWrite             = "X-Container-Write"
	hTempURLKey                 = "X-Container-Meta-Temp-Url-Key"
	hTempURLKey2                = "X-Container-Meta-Temp-Url-Key-2"
	hAccessControlAllowOrigin   = "X-Container-Meta-Access-Control-Allow-Origin"
	hAccessControlExposeHeaders = "X-Container-Meta-Access-Control-Expose-Headers"
	hAccessControlMaxAge        = "X-Container-Meta-Access-Control-Max-Age"
	hQuotaBytes                 = "X-Container-Meta-Quota-Bytes"
	hQuotaCount                 = "X-Container-Meta-Quota-Count"
	hPolicyGeoreplication       = "X-Container-Meta-Policy-Georeplication"

	hMetaPrefix       = "X-Container-Meta-"
	hRemoveMetaPrefix = "X-Remove-Container-Meta-"
)

// All X-Container-Meta-* attributes that are explictly declared in the
// Container data types. Use to distigush standard from customer attributes
var explicitMetaHeaders = []string{
	hTempURLKey,
	hTempURLKey2,
	hAccessControlAllowOrigin,
	hAccessControlExposeHeaders,
	hAccessControlMaxAge,
	hQuotaBytes,
	hQuotaCount,
	hPolicyGeoreplication,
}

// Determine if a given header is a standard attribute or custom header
func (c *Client) isCustomHeader(header string) bool {
	for _, v := range explicitMetaHeaders {
		if v == header {
			return false
		}
	}
	return true
}

// Container describes an existing Container.
type Container struct {
	// The name of the Container
	Name string
	// A container access control list (ACL) that grants read access.
	ReadACLs []string
	// A container access control list (ACL) that grants write access
	WriteACLs []string
	// The secret key value for temporary URLs.
	PrimaryKey string
	// The second secret key value for temporary URLs.
	SecondaryKey string
	// List of origins to be allowed to make cross-origin Requests.
	AllowedOrigins []string
	// List of headers exposed to the user agent (e.g. browser) in the actual request response.
	ExposedHeaders []string
	// Maximum age in seconds for the origin to hold the preflight results.
	MaxAge int
	// Maximum size of the container, in bytes
	QuotaBytes int
	// Maximum object count of the container
	QuotaCount int
	// Map of custom Container X-Container-Meta-{name} name value pairs
	CustomMetadata map[string]string
	// Georeplication Policy (undocumented)
	GeoreplicationPolicy []string
}

// CreateContainerInput defines an Container to be created.
type CreateContainerInput struct {
	// The unique name for the container. The container name must be from 1 to 256 characters long and can
	// start with any character and contain any pattern. Character set must be UTF-8. The container name
	// cannot contain a slash (/) character because this character delimits the container and object name.
	// For example, /account/container/object.
	// Required
	Name string `json:"name"`
	// Sets a container access control list (ACL) that grants read access.
	// Optional
	ReadACLs []string
	// Sets a container access control list (ACL) that grants read access.
	// Optional
	WriteACLs []string
	// Sets a secret key value for temporary URLs.
	// Optional
	PrimaryKey string
	// Sets a second secret key value for temporary URLs.
	// Optional
	SecondaryKey string
	// Sets the list of origins allowed to make cross-origin requests.
	// Optional
	AllowedOrigins []string
	// List of headers exposed to the user agent (e.g. browser) in the actual request response.
	// Optional
	ExposedHeaders []string
	// Sets the maximum age in seconds for the origin to hold the preflight results.
	// Optional
	MaxAge int
	// Sets the Maximum size of the container, in bytes
	// Optional
	QuotaBytes int
	// Sets the Maximum object count of the container
	// Optional
	QuotaCount int
	// Map of custom Container X-Container-Meta-{name} name value pairs
	// Optional
	CustomMetadata map[string]string
	// Georeplication Policy (undocumented)
	// GeoreplicationPolicy []string
}

// CreateContainer creates a new Container with the given name, key and enabled flag.
func (c *Client) CreateContainer(input *CreateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	input.Name = c.getQualifiedName(input.Name)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(input.ReadACLs) > 0 {
		headers[hContainerRead] = strings.Join(input.ReadACLs, ",")
	}
	if len(input.WriteACLs) > 0 {
		headers[hContainerWrite] = strings.Join(input.WriteACLs, ",")
	}

	headers[hTempURLKey] = input.PrimaryKey
	headers[hTempURLKey2] = input.SecondaryKey
	headers[hAccessControlAllowOrigin] = strings.Join(input.AllowedOrigins, " ")
	headers[hAccessControlExposeHeaders] = strings.Join(input.ExposedHeaders, " ")
	// headers[hPolicyGeoreplication] = strings.Join(input.GeoreplicationPolicy, " ")

	if input.MaxAge != 0 {
		headers[hAccessControlMaxAge] = strconv.Itoa(input.MaxAge)
	}
	if input.QuotaBytes != 0 {
		headers[hQuotaBytes] = strconv.Itoa(input.QuotaBytes)
	}
	if input.QuotaCount != 0 {
		headers[hQuotaCount] = strconv.Itoa(input.QuotaCount)
	}

	if len(input.CustomMetadata) > 0 {
		// add a header entry for each custom metadata item
		// X-Container-Meta-{name}: value
		for name, value := range input.CustomMetadata {
			header := fmt.Sprintf("%s%s", hMetaPrefix, name)
			if c.isCustomHeader(header) {
				headers[header] = value
			}
		}
	}

	if err := c.createResource(input.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: input.Name,
	}

	return c.GetContainer(&getInput)
}

// DeleteContainerInput describes the container to delete
type DeleteContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:"name"`
}

// DeleteContainer deletes the Container with the given name.
func (c *Client) DeleteContainer(input *DeleteContainerInput) error {
	input.Name = c.getQualifiedName(input.Name)
	return c.deleteResource(input.Name)
}

// GetContainerInput describes the container to get
type GetContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:"name"`
}

// GetContainer retrieves the Container with the given name.
func (c *Client) GetContainer(input *GetContainerInput) (*Container, error) {
	var container Container
	input.Name = c.getQualifiedName(input.Name)

	rsp, err := c.getResource(input.Name, &container)
	if err != nil {
		return nil, err
	}
	// The response doesn't come back with the name so we need to set it from the Input Name
	container.Name = c.getUnqualifiedName(input.Name)
	return c.success(rsp, &container)
}

// UpdateContainerInput defines an Container to be updated
type UpdateContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:"name"`
	// Updates a container access control list (ACL) that grants read access.
	// Optional
	ReadACLs []string
	// Updates a container access control list (ACL) that grants write access.
	// Optional
	WriteACLs []string
	// Updates the secret key value for temporary URLs.
	// Optional
	PrimaryKey string
	// Update the second secret key value for temporary URLs.
	// Optional
	SecondaryKey string
	// Updates the list of origins allowed to make cross-origin requests.
	// Optional
	AllowedOrigins []string
	// List of headers exposed to the user agent (e.g. browser) in the actual request response.
	// Optional
	ExposedHeaders []string
	// Updates the maximum age in seconds for the origin to hold the preflight results.
	// Optional
	MaxAge int
	// Updates the Maximum size of the container, in bytes
	// Optional
	QuotaBytes int
	// Updates the Maximum object count of the container
	// Optional
	QuotaCount int
	// Updates custom Container X-Container-Meta-{name} name value pairs
	// Optional
	CustomMetadata map[string]string
	// Remove custom Container X-Container-Meta-{name} headers
	// Optional
	RemoveCustomMetadata []string
	// Georeplication Policy (undocumented)
	// GeoreplicationPolicy []string
}

// Set an X-Container-Meta-{name} header with the value provided
// or if the value is empty set the X-Remove-Container-Meta-{name} header
func (c *Client) updateOrRemoveStringValue(headers map[string]string, header, value string) {
	if value == "" {
		headers[strings.Replace(header, hMetaPrefix, hRemoveMetaPrefix, 1)] = ""
	} else {
		headers[header] = value
	}
}

// Set an X-Container-Meta-{name} header with the value provided
// or if the value is 0 set the X-Remove-Container-Meta-{name} header
func (c *Client) updateOrRemoveIntValue(headers map[string]string, header string, value int) {
	if value == 0 {
		headers[strings.Replace(header, hMetaPrefix, hRemoveMetaPrefix, 1)] = ""
	} else {
		headers[header] = strconv.Itoa(value)
	}
}

// UpdateContainer updates the key and enabled flag of the Container with the given name.
func (c *Client) UpdateContainer(input *UpdateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(input.ReadACLs) > 0 {
		headers[hContainerRead] = strings.Join(input.ReadACLs, ",")
	}
	if len(input.WriteACLs) > 0 {
		headers[hContainerWrite] = strings.Join(input.WriteACLs, ",")
	}

	c.updateOrRemoveStringValue(headers, hTempURLKey, input.PrimaryKey)
	c.updateOrRemoveStringValue(headers, hTempURLKey2, input.SecondaryKey)
	c.updateOrRemoveStringValue(headers, hAccessControlAllowOrigin, strings.Join(input.AllowedOrigins, " "))
	c.updateOrRemoveStringValue(headers, hAccessControlExposeHeaders, strings.Join(input.ExposedHeaders, " "))
	c.updateOrRemoveIntValue(headers, hAccessControlMaxAge, input.MaxAge)
	c.updateOrRemoveIntValue(headers, hQuotaBytes, input.QuotaBytes)
	c.updateOrRemoveIntValue(headers, hQuotaCount, input.QuotaCount)
	// c.updateOrRemove(headers, hPolicyGeoreplication, strings.Join(input.GeoreplicationPolicy, " "))

	if len(input.CustomMetadata) > 0 {
		// add a header entry for each custom metadata item
		// X-Container-Meta-{name}: value
		for name, value := range input.CustomMetadata {
			header := fmt.Sprintf("%s%s", hMetaPrefix, name)
			if c.isCustomHeader(header) {
				headers[header] = value
			}
		}
	}

	if len(input.RemoveCustomMetadata) > 0 {
		// add a special header entry for each custom metadata item to be removed
		// X-Remove-Container-Meta-{name}: value
		for _, name := range input.RemoveCustomMetadata {
			header := fmt.Sprintf("%s%s", hMetaPrefix, name)
			if c.isCustomHeader(header) {
				// change to remove header
				header = fmt.Sprintf("%s%s", hRemoveMetaPrefix, name)
				headers[header] = ""
			}
		}
	}

	input.Name = c.getQualifiedName(input.Name)
	if err := c.updateResource(input.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: input.Name,
	}
	return c.GetContainer(&getInput)
}

func (c *Client) success(rsp *http.Response, container *Container) (*Container, error) {
	var (
		err        error
		maxAge     int
		quotaBytes int
		quotaCount int
	)

	container.ReadACLs = strings.Split(rsp.Header.Get(hContainerRead), ",")
	container.WriteACLs = strings.Split(rsp.Header.Get(hContainerWrite), ",")
	container.PrimaryKey = rsp.Header.Get(hTempURLKey)
	container.SecondaryKey = rsp.Header.Get(hTempURLKey2)
	container.AllowedOrigins = strings.Split(rsp.Header.Get(hAccessControlAllowOrigin), " ")
	container.ExposedHeaders = strings.Split(rsp.Header.Get(hAccessControlExposeHeaders), " ")
	container.GeoreplicationPolicy = strings.Split(rsp.Header.Get(hPolicyGeoreplication), " ")

	if maxAge, err = strconv.Atoi(rsp.Header.Get(hAccessControlMaxAge)); err == nil {
		container.MaxAge = maxAge
	}
	if quotaBytes, err = strconv.Atoi(rsp.Header.Get(hQuotaBytes)); err == nil {
		container.QuotaBytes = quotaBytes
	}
	if quotaCount, err = strconv.Atoi(rsp.Header.Get(hQuotaCount)); err == nil {
		container.QuotaCount = quotaCount
	}

	container.CustomMetadata = make(map[string]string)
	for header, value := range rsp.Header {
		if strings.HasPrefix(header, hMetaPrefix) && c.isCustomHeader(header) {
			name := strings.TrimPrefix(header, hMetaPrefix)
			container.CustomMetadata[name] = strings.Join(value, " ")
		}
	}

	if err != nil {
		return nil, err
	}

	return container, nil
}

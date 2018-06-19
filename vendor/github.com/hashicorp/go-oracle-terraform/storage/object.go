//- Object Resource + Data Source
//-
//- Satisfies Create, Read, Delete.
//- Object Metadata should be handled in a separate resource
//- Can only replace objects, so no Update method, use ForceNew in Terraform

package storage

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// ObjectClient details the parameters needed for a storage object client
type ObjectClient struct {
	Client
}

// Objects returns an object client
func (c *Client) Objects() *ObjectClient {
	return &ObjectClient{
		Client: Client{
			client:      c.client,
			authToken:   c.authToken,
			tokenIssued: c.tokenIssued,
		},
	}
}

// Header Constants
const (
	hAcceptRanges       = "Accept-Ranges"
	hContentDisposition = "Content-Disposition"
	hContentEncoding    = "Content-Encoding"
	hContentLength      = "Content-Length"
	hContentType        = "Content-Type"
	hCopyFrom           = "X-Copy-From"
	hDate               = "Date"
	hDeleteAt           = "X-Delete-At"
	hETag               = "ETag"
	hLastModified       = "Last-Modified"
	hNewest             = "X-Newest"
	hObjectManifest     = "X-Object-Manifest"
	hRange              = "Range"
	hTimestamp          = "X-Timestamp"
	hTransactionID      = "X-Trans-Id"
	hTransferEncoding   = "Transfer-Encoding"
	hMetadataPrefix     = "X-Object-Meta-"
)

// ObjectInfo describes an existing object
// Optional values may not be passed in as response headers
// TODO: Add query parameters if needed
type ObjectInfo struct {
	// ID is the container name + "/" object name for convenience
	ID string
	// Name of the object
	Name string
	// Type of ranges the object accepts
	AcceptRanges string
	// Name of the container
	Container string
	// Optional: Specifies the override behavior for the browser
	ContentDisposition string
	// Optional: Content's Encoding header
	ContentEncoding string
	// Length of the object in bytes
	ContentLength int
	// Type of the content
	ContentType string
	// Date of the transaction in ISO 8601 format.
	// Null value means the token never expires
	Date string
	// For objects smaller than 5GB, MD5 checksum of the object content.
	// Otherwise MD5 sum of the concatenated string of MD5 sums and ETAGS
	// for each segment of the manifest. Enclosed in double-quote characters
	Etag string
	// Date and time when the object was created/modified. ISO 8601.
	LastModified string
	// Optional: Date+Time in EPOCH that the object will be deleted.
	DeleteAt int
	// Optional: The dynamic large object manifest object.
	ObjectManifest string
	// Optional: The map of object metadata name values pairs for X-Object-Meta-{name}
	ObjectMetadata map[string]string
	// Date and time in UNIX EPOCH when the account, container, _or_ object
	// was initially created as a current version.
	Timestamp string
	// Transaction ID of the request - Used for bug reports to service providers
	TransactionID string
}

// CreateObjectInput struct for a Create Method to create a storage object
// TODO: Add query parameters if needed
type CreateObjectInput struct {
	// Name of the object.
	// Required
	Name string
	// Body of the request to use. Accepts an io.ReadSeeker, so options are open to
	// the downstream consumer
	// Required
	Body io.ReadSeeker
	// Name of the container to place the object
	// Required
	Container string
	// Override the behavior of the browser.
	// Optional
	ContentDisposition string
	// Set the content-encoding metadata
	// Optional
	ContentEncoding string
	// Changes the MIME type for the object
	// Optional - Defaults to 'text/plain'
	ContentType string
	// Specify the `container/object` to copy from. Must be UTF-8 encoded
	// and the name of the container and object must be URL-encoded
	// Optional
	CopyFrom string
	// Specify the date and time in UNIX Epoch time stamp format when the system
	// removes the object
	DeleteAt int
	// Specify the map of object metadata name values pairs for X-Object-Meta-{name}
	ObjectMetadata map[string]string
	// MD5 checksum value of the request body. Unquoted
	// Strongly recommended, not required.
	ETag string
	// TODO: If-None-Match.

	// Sets the transfer encoding. Can only be "chunked" or nil.
	// Requires content-length to be 0 if set.
	// Optional
	TransferEncoding string
	// TODO: X-Object-Meta-{name}
}

// CreateObject creates a new Object inside of a container.
func (c *ObjectClient) CreateObject(input *CreateObjectInput) (*ObjectInfo, error) {
	headers := make(map[string]string)

	name := c.getQualifiedName(fmt.Sprintf("%s/%s", input.Container, input.Name))

	if input.ContentDisposition != "" {
		headers[hContentDisposition] = input.ContentDisposition
	}
	if input.ContentEncoding != "" {
		headers[hContentEncoding] = input.ContentEncoding
	}
	if input.ContentType != "" {
		headers[hContentType] = input.ContentType
	}
	if input.ETag != "" {
		headers[hETag] = input.ETag
	}
	if input.TransferEncoding != "" {
		headers[hTransferEncoding] = input.TransferEncoding
	}
	if input.CopyFrom != "" {
		headers[hCopyFrom] = input.CopyFrom
	}
	if input.DeleteAt != 0 {
		headers[hDeleteAt] = fmt.Sprintf("%d", input.DeleteAt)
	}
	if len(input.ObjectMetadata) > 0 {
		// add a header entry for each metadata item
		// X-Object-Meta-{name}: value
		for name, value := range input.ObjectMetadata {
			header := fmt.Sprintf("%s%s", hMetadataPrefix, name)
			headers[header] = value
		}
	}

	if input.Body == nil && input.CopyFrom == "" {
		return nil, fmt.Errorf("Body cannot be nil")
	}

	if err := c.createResourceBody(name, headers, input.Body); err != nil {
		return nil, err
	}

	getInput := &GetObjectInput{
		Name:      input.Name,
		Container: input.Container,
	}

	return c.GetObject(getInput)
}

// GetObjectInput details on a storage object
// TODO: Add query parameters if needed
type GetObjectInput struct {
	// TODO If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since
	// If we actually want to support these

	// ID of the object (container/object)
	// Optional - Either ID or Name + Container are required
	ID string

	// Name of the object to get details on
	// Optional - Either ID or Name + Container are required
	Name string
	// Name of the container
	// Optional - Either ID or Name + Container are required
	Container string
	// Range of data to receive. Must be specified via a byte range:
	// bytes=-5; bytes=10-15. Accept the entire string here, as multiple ranges
	// can be specified with a comma delimiter
	// Optional
	Range string
	// If set to true, Object Storage queries all replicas to return the most recent one.
	// If you omit this header, Object Storage responds faster after it finds one valid replica.
	// Because setting this header to true is more expensive for the back end, use it only when
	// it is absolutely needed.
	// Optional
	Newest bool
}

// GetObject accepts a input struct, returns an info struct
func (c *ObjectClient) GetObject(input *GetObjectInput) (*ObjectInfo, error) {
	var object ObjectInfo
	headers := make(map[string]string)

	name, err := c.getIdentifier(input.ID, input.Container, input.Name)
	if err != nil {
		return nil, err
	}

	// Build request headers
	headers[hRange] = input.Range
	headers[hNewest] = fmt.Sprintf("%t", input.Newest)

	resp, err := c.getResourceHeaders(name, &object, headers)
	if err != nil {
		return nil, err
	}

	// Set Name, container, and ID. Not returned from API
	if input.ID != "" {
		parts := strings.Split(input.ID, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Unknown ID specified: %s", input.ID)
		}
		object.ID = input.ID
		object.Container = parts[0]
		object.Name = parts[1]
	} else {
		// Already checked for Nil container and name above
		object.ID = fmt.Sprintf("%s/%s", input.Container, input.Name)
		object.Name = input.Name
		object.Container = input.Container
	}

	return c.success(resp, &object)
}

// DeleteObjectInput struct for deleting objects
// TODO: Add query parameters if needed
type DeleteObjectInput struct {
	// ID is the container name + "/" + object name
	// Optional - Either ID or Name + Container are required
	ID string
	// Name of the Object to delete
	// Optional - Either ID or Name + Container are required
	Name string
	// Name of the container
	// Optional - Either ID or Name + Container are required
	Container string
}

// DeleteObject will delete the supplied object
func (c *ObjectClient) DeleteObject(input *DeleteObjectInput) error {
	name, err := c.getIdentifier(input.ID, input.Container, input.Name)
	if err != nil {
		return err
	}

	return c.deleteResource(c.getQualifiedName(name))
}

func (c *ObjectClient) success(resp *http.Response, object *ObjectInfo) (*ObjectInfo, error) {
	var err error
	// Translate response headers into object info struct
	object.AcceptRanges = resp.Header.Get(hAcceptRanges)
	object.ContentDisposition = resp.Header.Get(hContentDisposition)
	object.ContentEncoding = resp.Header.Get(hContentEncoding)
	object.ContentType = resp.Header.Get(hContentType)
	object.Date = resp.Header.Get(hDate)
	object.Etag = resp.Header.Get(hETag)
	object.LastModified = resp.Header.Get(hLastModified)
	object.ObjectManifest = resp.Header.Get(hObjectManifest)
	object.Timestamp = resp.Header.Get(hTimestamp)
	object.TransactionID = resp.Header.Get(hTransactionID)

	if v := resp.Header.Get(hContentLength); v != "" {
		object.ContentLength, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	if v := resp.Header.Get(hDeleteAt); v != "" {
		object.DeleteAt, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	object.ObjectMetadata = make(map[string]string)
	for header, value := range resp.Header {
		if strings.HasPrefix(header, hMetadataPrefix) {
			name := strings.TrimPrefix(header, hMetadataPrefix)
			object.ObjectMetadata[name] = strings.Join(value, " ")
		}
	}

	return object, nil
}

func (c *ObjectClient) getIdentifier(id, container, name string) (string, error) {
	var result string
	if id != "" {
		result = id
	} else {
		if container == "" && name == "" {
			return "", fmt.Errorf("Either ID or Name and Container must be set during DELETE")
		}
		result = fmt.Sprintf("%s/%s", container, name)
	}

	return c.getQualifiedName(result), nil
}

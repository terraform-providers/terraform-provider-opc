package storage

import (
	"net/http"
	"strconv"
	"strings"
)

const CONTAINER_VERSION = "v1"

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
	// Maximum age in seconds for the origin to hold the preflight results.
	MaxAge int
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
	// Sets the maximum age in seconds for the origin to hold the preflight results.
	// Optional
	MaxAge int
}

// CreateContainer creates a new Container with the given name, key and enabled flag.
func (c *StorageClient) CreateContainer(input *CreateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	input.Name = c.getQualifiedName(CONTAINER_VERSION, input.Name)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(input.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(input.ReadACLs, ",")
	}
	if len(input.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(input.WriteACLs, ",")
	}

	headers["X-Container-Meta-Temp-URL-Key"] = input.PrimaryKey
	headers["X-Container-Meta-Temp-URL-Key-2"] = input.SecondaryKey
	headers["X-Container-Meta-Access-Control-Expose-Headers"] = strings.Join(input.AllowedOrigins, " ")
	headers["X-Container-Meta-Access-Control-Max-Age"] = strconv.Itoa(input.MaxAge)

	if err := c.createResource(input.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: input.Name,
	}

	return c.GetContainer(&getInput)
}

// DeleteKeyInput describes the container to delete
type DeleteContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:name`
}

// DeleteContainer deletes the Container with the given name.
func (c *StorageClient) DeleteContainer(input *DeleteContainerInput) error {
	input.Name = c.getQualifiedName(CONTAINER_VERSION, input.Name)
	return c.deleteResource(input.Name)
}

// GetContainerInput describes the container to get
type GetContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:name`
}

// GetContainer retrieves the Container with the given name.
func (c *StorageClient) GetContainer(input *GetContainerInput) (*Container, error) {
	var container Container
	input.Name = c.getQualifiedName(CONTAINER_VERSION, input.Name)

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
	// Updates the maximum age in seconds for the origin to hold the preflight results.
	// Optional
	MaxAge int
}

// UpdateContainer updates the key and enabled flag of the Container with the given name.
func (c *StorageClient) UpdateContainer(input *UpdateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(input.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(input.ReadACLs, ",")
	}
	if len(input.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(input.WriteACLs, ",")
	}

	headers["X-Container-Meta-Temp-URL-Key"] = input.PrimaryKey
	headers["X-Container-Meta-Temp-URL-Key-2"] = input.SecondaryKey
	headers["X-Container-Meta-Access-Control-Expose-Headers"] = strings.Join(input.AllowedOrigins, " ")
	headers["X-Container-Meta-Access-Control-Max-Age"] = strconv.Itoa(input.MaxAge)

	input.Name = c.getQualifiedName(CONTAINER_VERSION, input.Name)
	if err := c.updateResource(input.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: input.Name,
	}
	return c.GetContainer(&getInput)
}

func (c *StorageClient) success(rsp *http.Response, container *Container) (*Container, error) {
	var err error
	container.ReadACLs = strings.Split(rsp.Header.Get("X-Container-Read"), ",")
	container.WriteACLs = strings.Split(rsp.Header.Get("X-Container-Write"), ",")
	container.PrimaryKey = rsp.Header.Get("X-Container-Meta-Temp-URL-Key")
	container.SecondaryKey = rsp.Header.Get("X-Container-Meta-Temp-URL-Key-2")
	container.AllowedOrigins = strings.Split(rsp.Header.Get("X-Container-Meta-Access-Control-Expose-Headers"), " ")
	container.MaxAge, err = strconv.Atoi(rsp.Header.Get("X-Container-Meta-Access-Control-Max-Age"))
	if err != nil {
		return nil, err
	}

	return container, nil
}

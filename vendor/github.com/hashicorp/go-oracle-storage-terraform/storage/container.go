package storage

import (
	"net/http"
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
	ReadACLs []string
	// Sets a container access control list (ACL) that grants read access.
	WriteACLs []string
}

// CreateContainer creates a new Container with the given name, key and enabled flag.
func (c *Client) CreateContainer(createInput *CreateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	createInput.Name = c.getQualifiedName(CONTAINER_VERSION, createInput.Name)

	if len(createInput.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(createInput.ReadACLs, ",")
	}
	if len(createInput.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(createInput.WriteACLs, ",")
	}

	if err := c.createResource(createInput.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: createInput.Name,
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
func (c *Client) DeleteContainer(deleteInput *DeleteContainerInput) error {
	deleteInput.Name = c.getQualifiedName(CONTAINER_VERSION, deleteInput.Name)
	return c.deleteResource(deleteInput.Name)
}

// GetContainerInput describes the container to get
type GetContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:name`
}

// GetContainer retrieves the Container with the given name.
func (c *Client) GetContainer(getInput *GetContainerInput) (*Container, error) {
	var (
		container Container
		rsp       *http.Response
		err       error
	)
	getInput.Name = c.getQualifiedName(CONTAINER_VERSION, getInput.Name)

	if rsp, err = c.getResource(getInput.Name, &container); err != nil {
		return nil, err
	}
	return c.success(rsp, &container)
}

// UpdateContainerInput defines an Container to be updated
type UpdateContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:"name"`
	// Updates a container access control list (ACL) that grants read access.
	ReadACLs []string
	// Updates a container access control list (ACL) that grants write access.
	WriteACLs []string
}

// UpdateContainer updates the key and enabled flag of the Container with the given name.
func (c *Client) UpdateContainer(updateInput *UpdateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	if len(updateInput.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(updateInput.ReadACLs, ",")
	}
	if len(updateInput.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(updateInput.WriteACLs, ",")
	}

	updateInput.Name = c.getQualifiedName(CONTAINER_VERSION, updateInput.Name)
	if err := c.updateResource(updateInput.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: updateInput.Name,
	}
	return c.GetContainer(&getInput)
}

func (c *Client) success(rsp *http.Response, container *Container) (*Container, error) {
	c.unqualify(&container.Name)
	container.ReadACLs = strings.Split(rsp.Header.Get("X-Container-Read"), ",")
	container.WriteACLs = strings.Split(rsp.Header.Get("X-Container-Write"), ",")

	return container, nil
}

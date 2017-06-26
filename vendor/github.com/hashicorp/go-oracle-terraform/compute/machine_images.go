package compute

// MachineImagesClient is a client for the MachineImage functions of the Compute API.
type MachineImagesClient struct {
	ResourceClient
}

// MachineImages obtains an MachineImagesClient which can be used to access to the
// MachineImage functions of the Compute API
func (c *ComputeClient) MachineImages() *MachineImagesClient {
	return &MachineImagesClient{
		ResourceClient: ResourceClient{
			ComputeClient:       c,
			ResourceDescription: "MachineImage",
			ContainerPath:       "/machineimage/",
			ResourceRootPath:    "/machineimage",
		}}
}

// DeleteMachineImageInput describes the snapshot to delete
type DeleteMachineImageInput struct {
	// The name of the MachineImage
	Name string `json:name`
}

// DeleteMachineImage deletes the MachineImage with the given name.
func (c *MachineImagesClient) DeleteMachineImage(deleteInput *DeleteMachineImageInput) error {
	return c.deleteResource(deleteInput.Name)
}

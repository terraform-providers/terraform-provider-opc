// Manages SSH Keys for a DBaaS Service Instance.
// SSH Keys can currently only be created and information fetched. They cannot
// be updated, or deleted via the API. So each interaction requires a ForceNew
// and "deleting" the resource simply removes the resource from state
// (in the provider), as there is no Delete method on the ssh key resource's
// API.

package database

import (
	"fmt"
	"strings"
	"time"
)

// API URI Paths for Container and Root objects.
// API Docs state to always set the 'credentialName' attribute in the
// resource's path to 'vmspublickey'. Hard-coding for now, but a better solution
// may need to be solved for SSH Keys if this value ever needs to be changed.
// See: https://docs.oracle.com/en/cloud/paas/database-dbaas-cloud/csdbr/op-paas-api-v1.1-instancemgmt-%7BidentityDomainId%7D-services-dbaas-instances-%7BserviceId%7D-credentials-crednames-%7BcredentialName%7D-post.html
// for more details on 'credentialName'.

// The Root Path, for SSH Keys also does not expect a conventional 'Name'
// path parameter, and the only difference between creating an SSH Key
// and viewing the summary of the current SSH Key is a PUT vs GET HTTP Method.
const (
	DBSSHKeyContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/credentials/crednames/vmspublickey"
	DBSSHKeyRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/credentials/%s"
	DBSSHKeyName          = "vmspublickey"
)

// Default timeout value for Create
// In testing this is anywhere between 10-20s depending on if it's a new SSH Key
// or if it's an "updated" ssh key.
const WaitForSSHKeyTimeout = time.Duration(30 * time.Second)

// SSHKeys returns a UtilityClient for managing SSHKeys for a DBaaS Service Instance
func (c *DatabaseClient) SSHKeys() *UtilityClient {
	return &UtilityClient{
		UtilityResourceClient: UtilityResourceClient{
			DatabaseClient:   c,
			ContainerPath:    DBSSHKeyContainerPath,
			ResourceRootPath: DBSSHKeyRootPath,
		},
	}
}

// SSHKeyInfo holds all the known information for a single SSH Key
type SSHKeyInfo struct {
	// This will likely always be "DB".
	ComponentType string `json:"componentType"`
	// The fully qualified name of the SSH Key object in OPC Cloud Storage
	// where the SSH public key value is stored.
	ComputeKeyName string `json:"computeKeyName"`
	// Should always be 'vmspublickey'
	CredName string `json:"credName"`
	// Should always be SSH
	CredType string `json:"credType"`
	// The string description of the key
	Description string `json:"description"`
	// Note: The API supplies us with the 'identityDomain' key here,
	// but since we are required to supply this during the API request
	// and it's already a "known" value, we don't return the value here as well.

	// The message returned from the last update of the SSH Key.
	LastUpdateMessage string `json:"lastUpdateMessage"`
	// Status of the last update of the SSH key
	LastUpdateStatus string `json:"lastUpdateStatus"`
	// Date and time of the last update of the SSH Key
	LastUpdateTime string `json:"lastUpdateTime"`
	// The value "opc"
	OsUserName string `json:"osUserName"`
	// The value "SERVICE"
	ParentType string `json:"parentType"`
	// The Value of the SSH Public Key with any slashes (/) it contains
	// preceded by backslashes: \/.
	PublicKey string `json:"publicKey"`
	// The name of the DatabaseCloudService Instance
	ServiceName string `json:"serviceName"`
	// The value "DBaaS"
	ServiceType string `json:"serviceType"`
}

// CreateSSHKeyInput defines the necessary input parameters to create an
// SSH Key for a DBaaS Service Instance
type CreateSSHKeyInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// The value of the SSH public key, with any slashes (/) it contains preceded by
	// backslashes, as in \/.
	// Required
	PublicKey string `json:"public-key"`
	// Time to wait for an ssh key to be ready
	Timeout time.Duration `json:"-"`
}

// Creates an SSH Key with the supplied input struct.
func (c *UtilityClient) CreateSSHKey(input *CreateSSHKeyInput) (*SSHKeyInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var sshKey SSHKeyInfo
	if err := c.createResource(input, &sshKey); err != nil {
		return nil, err
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForSSHKeyTimeout
	}

	// Can leave ServiceInstanceID nil here, it will be the same as the current client's
	result, err := c.WaitForSSHKeyReady(&GetSSHKeyInput{}, timeout)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetSSHKeyInput defines all of the necessary input parameters to retrieve the necessary
// information on a specific SSH Key.
type GetSSHKeyInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
}

// Get's information on a single SSH Key
func (c *UtilityClient) GetSSHKey(input *GetSSHKeyInput) (*SSHKeyInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Name has to be populated in this case as the Container path and the Root path are completely
	// separate paths. Otherwise, with a nil name, the Container path would have been used, which
	// would effectively return a '200 OK' for each request, but only return the summary for an SSH Key
	// instead of details
	var sshKey SSHKeyInfo
	if err := c.getResource(DBSSHKeyName, &sshKey); err != nil {
		return nil, err
	}

	return &sshKey, nil
}

// No Delete, or Update currently.
// TODO: Add Delete and Update for SSH Keys when they are available in the API.

func (c *UtilityClient) WaitForSSHKeyReady(input *GetSSHKeyInput, timeout time.Duration) (*SSHKeyInfo, error) {
	var info *SSHKeyInfo
	var getErr error
	err := c.client.WaitFor("sshkey to be ready", timeout, func() (bool, error) {
		info, getErr = c.GetSSHKey(input)
		if getErr != nil {
			return false, getErr
		}
		if info != nil {
			c.client.DebugLogString(fmt.Sprintf("SSH Key Status: %s", strings.ToLower(info.LastUpdateStatus)))
			success := strings.ToLower(info.LastUpdateStatus) == "success"
			return success, nil
		}
		// Not found, wait
		c.client.DebugLogString(fmt.Sprintf("SSH Key not found, waiting"))
		return false, nil
	})
	return info, err
}

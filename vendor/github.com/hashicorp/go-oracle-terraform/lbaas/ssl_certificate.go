package lbaas

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForSSLCertificateReadyPollInterval = 1 * time.Second
const waitForSSLCertificateReadyTimeout = 5 * time.Minute
const waitForSSLCertificateDeletePollInterval = 1 * time.Second
const waitForSSLCertificateDeleteTimeout = 5 * time.Minute

type SSLCertificateInfo struct {
	Name             string     `json:"name"`
	Certificate      string     `json:"certificate"`
	CertificateChain string     `json:"certificate_chain"`
	State            LBaaSState `json:"state"`
	Trusted          bool       `json:"trusted"`
	URI              string     `json:"uri"`
}

type CreateSSLCertificateInput struct {
	Name             string `json:"name"`
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificate_chain,omitempty"`
	PrivateKey       string `json:"private_key,omitempty"`
	Trusted          bool   `json:"trusted"`
}

type UpdateSSLCertificateInput struct {
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key,omitempty"`
}

// CreateSSLCertificate creates a new SSL certificate
func (c *SSLCertificateClient) CreateSSLCertificate(input *CreateSSLCertificateInput) (*SSLCertificateInfo, error) {

	if input.Trusted {
		c.ContentType = ContentTypeTrustedCertificateJSON
	} else {
		c.ContentType = ContentTypeServerCertificateJSON
	}

	var info SSLCertificateInfo
	if err := c.createResource(&input, &info); err != nil {
		return nil, err
	}

	createdStates := []LBaaSState{LBaaSStateCreated, LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateCreationFailed, LBaaSStateDeletionInProgress, LBaaSStateDeleted, LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkSSLCertificateState(&info, createdStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForSSLCertificateState(input.Name, createdStates, erroredStates, c.PollInterval, c.Timeout, &info)
	return &info, err
}

// DeleteSSLCertificate deletes the SSL certificate with the specified name
func (c *SSLCertificateClient) DeleteSSLCertificate(name string) (*SSLCertificateInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForSSLCertificateDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForSSLCertificateDeleteTimeout
	}

	var info SSLCertificateInfo
	if err := c.deleteResource(name, &info); err != nil {
		return nil, err
	}

	deletedStates := []LBaaSState{LBaaSStateDeleted}
	erroredStates := []LBaaSState{LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	deleted, err := c.checkSSLCertificateState(&info, deletedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if deleted {
		return &info, nil
	}
	// else poll till deleted
	err = c.WaitForSSLCertificateState(name, deletedStates, erroredStates, c.PollInterval, c.Timeout, &info)
	if err != nil && client.WasNotFoundError(err) {
		// resource could not be found, thus deleted
		return nil, nil
	}
	return &info, nil
}

// GetSSLCertificate fetch the SSL Certificate details
func (c *SSLCertificateClient) GetSSLCertificate(name string) (*SSLCertificateInfo, error) {
	var info SSLCertificateInfo
	if err := c.getResource(name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// WaitForSSLCertificateState waits for the resource to be in one of a set of desired states
func (c *SSLCertificateClient) WaitForSSLCertificateState(name string, desiredStates, errorStates []LBaaSState, pollInterval, timeoutSeconds time.Duration, info *SSLCertificateInfo) error {

	var getErr error
	err := c.client.WaitFor("SSL Certificate status update", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetSSLCertificate(name)
		if getErr != nil {
			return false, getErr
		}

		return c.checkSSLCertificateState(info, desiredStates, errorStates)
	})
	return err
}

// check the State, returns in desired state (true), not ready yet (false) or errored state (error)
func (c *SSLCertificateClient) checkSSLCertificateState(info *SSLCertificateInfo, desiredStates, errorStates []LBaaSState) (bool, error) {

	c.client.DebugLogString(fmt.Sprintf("SSL Certificate %v state is %v", info.Name, info.State))

	state := LBaaSState(info.State)

	if isStateInLBaaSStates(state, desiredStates) {
		// we're good, return okay
		return true, nil
	}
	if isStateInLBaaSStates(state, errorStates) {
		// not good, return error
		return false, fmt.Errorf("SSL Certificate %v in errored state %v", info.Name, info.State)
	}
	// not ready lifecycleTimeout
	return false, nil
}

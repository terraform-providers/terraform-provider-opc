package lbaas

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForPolicyReadyPollInterval = 1 * time.Second
const waitForPolicyReadyTimeout = 5 * time.Minute
const waitForPolicyDeletePollInterval = 1 * time.Second
const waitForPolicyDeleteTimeout = 5 * time.Minute

const (
	policyContainerPath = "/vlbrs/%s/%s/policies"
	policyResourcePath  = "/vlbrs/%s/%s/policies/%s"
)

// Policy Specific ContentTypes for API REquests. Each type of Policy has its own ContentType
const ContentTypeAppCookieSticinessPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.AppCookieStickinessPolicy+json"
const ContentTypeCloudGatePolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.CloudGatePolicy+json"
const ContentTypeLBCookieStickinessPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.LBCookieStickinessPolicy+json"
const ContentTypeLoadBalancingMechanismPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.LoadBalancingMechanismPolicy+json"
const ContentTypeRateLimitingRequestPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.RateLimitingRequestPolicy+json"
const ContentTypeRedirectPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.RedirectPolicy+json"
const ContentTypeResourceAccessControlPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.ResourceAccessControlPolicy+json"
const ContentTypeSetRequestHeaderPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.SetRequestHeaderPolicy+json"
const ContentTypeSSLNegotiationPolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.SSLNegotiationPolicy+json"
const ContentTypeTrustedCertificatePolicyJSON = "application/vnd.com.oracle.oracloud.lbaas.TrustedCertPolicy+json"

// PolicyClient is a client for the Load Balancer Policy resources.
type PolicyClient struct {
	LBaaSResourceClient
}

// PolicyClient returns an PolicyClient which is used to access the
// Load Balancer Policy API
func (c *Client) PolicyClient() *PolicyClient {

	return &PolicyClient{
		LBaaSResourceClient: LBaaSResourceClient{
			Client:           c,
			ContainerPath:    policyContainerPath,
			ResourceRootPath: policyResourcePath,
			// Accept all Policy Content Types
			Accept: strings.Join([]string{
				ContentTypeAppCookieSticinessPolicyJSON,
				ContentTypeCloudGatePolicyJSON,
				ContentTypeLBCookieStickinessPolicyJSON,
				ContentTypeLoadBalancingMechanismPolicyJSON,
				ContentTypeRateLimitingRequestPolicyJSON,
				ContentTypeResourceAccessControlPolicyJSON,
				ContentTypeRedirectPolicyJSON,
				ContentTypeSSLNegotiationPolicyJSON,
				ContentTypeSetRequestHeaderPolicyJSON,
				ContentTypeTrustedCertificatePolicyJSON,
			}, ","),
			// ContentType cannot be generally set for the PolicyClient, instead it is set on each
			// Create or Update request based on the Type of the Policy being created/updated.
		},
	}
}

type PolicyInfo struct {
	Name  string     `json:"name,omitempty"`
	State LBaaSState `json:"state,omitempty"`
	Type  string     `json:"type,omitempty"`
	URI   string     `json:"uri,omitempty"`

	// ApplicationCookieStickinessPolicy
	AppCookieName string `json:"app_cookie_name,omitempty"`

	// CloudGatePolicy
	CloudGateApplication                string `json:"cloudgate_application,omitempty"`
	CloudGatePolicyName                 string `json:"cloudgate_policy_name,omitempty"`
	IdentityServiceInstanceGuid         string `json:"identity_service_instance_guid,omitempty"`
	VirtualHostnameForPolicyAttribution string `json:"virtual_hostname_for_policy_attribution,omitempty"`

	// LoadBalancerCookieStickinessPolicy
	CookieExpirationPeriod int `json:"cookie_expiration_period,omitempty"`

	// LoadBalancingMechanismPolicy
	LoadBalancingMechanism string `json:"load_balancing_mechanism,omitempty"`

	// RateLimitingRequestPolicy
	BurstSize                   int    `json:"burst_size,omitempty"`
	DoNotDelayExcessiveRequests bool   `json:"do_not_delay_excessive_requests,omitempty"`
	HttpStatusErrorCode         int    `json:"http_status_error_code,omitempty"`
	LogLevel                    string `json:"log_level,omitempty"`
	RateLimitingCriteria        string `json:"rate_limiting_criteria,omitempty"`
	RequestsPerSecond           int    `json:"requests_per_second,omitempty"`
	StorageSize                 int    `json:"storage_size_in_mb,omitempty"`
	Zone                        string `json:"zone,omitempty"`

	// RedirectPolicy
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ResponseCode int    `json:"response_code,omitempty"`

	// ResourceAccessControlPolicy
	Disposition      string   `json:"disposition,omitempty"`
	DeniedClients    []string `json:"denied_clients,omitempty"`
	PermittedClients []string `json:"permitted_clients,omitempty"`

	// SetRequestHeaderPolicy
	HeaderName                 string   `json:"header_name,omitempty"`
	Value                      string   `json:"value,omitempty"`
	ActionWhenHeaderExists     string   `json:"action_when_hdr_exists,omitempty"`
	ActionWhenHeaderValueIs    []string `json:"action_when_hdr_value_is,omitempty"`
	ActionWhenHeaderValueIsNot []string `json:"action_when_hdr_value_is_not,omitempty"`

	// SSLNegotiationPolicy
	Port                  int      `json:"port,omitempty"`
	ServerOrderPreference string   `json:"server_order_preference,omitempty"`
	SSLProtocol           []string `json:"ssl_protocol,omitempty"`
	SSLCiphers            []string `json:"ssl_ciphers,omitempty"`

	// TrustedCertificatePolicy
	TrustedCertificate string `json:"cert,omitempty"`
}

type CreatePolicyInput struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	ApplicationCookieStickinessPolicyInfo
	CloudGatePolicyInfo
	LoadBalancerCookieStickinessPolicyInfo
	LoadBalancingMechanismPolicyInfo
	RateLimitingRequestPolicyInfo
	RedirectPolicyInfo
	ResourceAccessControlPolicyInfo
	SetRequestHeaderPolicyInfo
	SSLNegotiationPolicyInfo
	TrustedCertificatePolicyInfo
}

type UpdatePolicyInput struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	ApplicationCookieStickinessPolicyInfo
	CloudGatePolicyInfo
	LoadBalancerCookieStickinessPolicyInfo
	LoadBalancingMechanismPolicyInfo
	RateLimitingRequestPolicyInfo
	RedirectPolicyInfo
	ResourceAccessControlPolicyInfo
	SetRequestHeaderPolicyInfo
	SSLNegotiationPolicyInfo
	TrustedCertificatePolicyInfo
}

type ApplicationCookieStickinessPolicyInfo struct {
	AppCookieName string `json:"app_cookie_name,omitempty"`
}

type CloudGatePolicyInfo struct {
	CloudGateApplication                string `json:"cloudgate_application,omitempty"`
	CloudGatePolicyName                 string `json:"cloudgate_policy_name,omitempty"`
	IdentityServiceInstanceGuid         string `json:"identity_service_instance_guid,omitempty"`
	VirtualHostnameForPolicyAttribution string `json:"virtual_hostname_for_policy_attribution,omitempty"`
}

type LoadBalancerCookieStickinessPolicyInfo struct {
	CookieExpirationPeriod int `json:"cookie_expiration_period,omitempty"`
}

type LoadBalancingMechanismPolicyInfo struct {
	LoadBalancingMechanism string `json:"load_balancing_mechanism,omitempty"`
}

type RateLimitingRequestPolicyInfo struct {
	BurstSize                   int    `json:"burst_size,omitempty"`
	DoNotDelayExcessiveRequests bool   `json:"do_not_delay_excessive_requests,omitempty"`
	HttpStatusErrorCode         int    `json:"http_status_error_code,omitempty"`
	LogLevel                    string `json:"log_level,omitempty"`
	RateLimitingCriteria        string `json:"rate_limiting_criteria,omitempty"`
	RequestsPerSecond           int    `json:"requests_per_second,omitempty"`
	StorageSize                 int    `json:"storage_size_in_mb,omitempty"`
	Zone                        string `json:"zone,omitempty"`
}

type RedirectPolicyInfo struct {
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ResponseCode int    `json:"response_code,omitempty"`
}

type ResourceAccessControlPolicyInfo struct {
	Disposition      string   `json:"disposition,omitempty"`
	DeniedClients    []string `json:"denied_clients,omitempty"`
	PermittedClients []string `json:"permitted_clients,omitempty"`
}

// use pointer for nilable fields on update
type ResourceAccessControlPolicyUpdate struct {
	Disposition      string    `json:"disposition,omitempty"`
	DeniedClients    *[]string `json:"denied_clients,omitempty"`
	PermittedClients *[]string `json:"permitted_clients,omitempty"`
}

type SetRequestHeaderPolicyInfo struct {
	HeaderName                 string   `json:"header_name,omitempty"`
	Value                      string   `json:"value,omitempty"`
	ActionWhenHeaderExists     string   `json:"action_when_hdr_exists,omitempty"`
	ActionWhenHeaderValueIs    []string `json:"action_when_hdr_value_is,omitempty"`
	ActionWhenHeaderValueIsNot []string `json:"action_when_hdr_value_is_not,omitempty"`
}

// use pointer for niable fields on update
type SetRequestHeaderPolicyUpdate struct {
	HeaderName                 string    `json:"header_name,omitempty"`
	Value                      *string   `json:"value,omitempty"`
	ActionWhenHeaderExists     string    `json:"action_when_hdr_exists,omitempty"`
	ActionWhenHeaderValueIs    *[]string `json:"action_when_hdr_value_is,omitempty"`
	ActionWhenHeaderValueIsNot *[]string `json:"action_when_hdr_value_is_not,omitempty"`
}

type SSLNegotiationPolicyInfo struct {
	Port                  int      `json:"port,omitempty"`
	ServerOrderPreference string   `json:"server_order_preference,omitempty"`
	SSLProtocol           []string `json:"ssl_protocol,omitempty"`
	SSLCiphers            []string `json:"ssl_ciphers,omitempty"`
}

type TrustedCertificatePolicyInfo struct {
	TrustedCertificate string `json:"cert,omitempty"`
}

// CreatePolicy creates a new listener
func (c *PolicyClient) CreatePolicy(lb LoadBalancerContext, input *CreatePolicyInput) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyReadyTimeout
	}

	info := &PolicyInfo{}
	// set the content type based on the policy type
	if input.Type == "" {
		return nil, fmt.Errorf("Policy type for %s is not set", input.Name)
	}
	c.ContentType = fmt.Sprintf("application/vnd.com.oracle.oracloud.lbaas.%s+json", input.Type)
	if err := c.createResource(lb.Region, lb.Name, &input, info); err != nil {
		return nil, err
	}

	createdStates := []LBaaSState{LBaaSStateCreated, LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateCreationFailed, LBaaSStateDeletionInProgress, LBaaSStateDeleted, LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkPolicyState(info, createdStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return info, nil
	}
	// else poll till ready
	info, err = c.WaitForPolicyState(lb, input.Name, createdStates, erroredStates, c.PollInterval, c.Timeout)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// DeletePolicy deletes the listener with the specified input
func (c *PolicyClient) DeletePolicy(lb LoadBalancerContext, name string) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyDeleteTimeout
	}

	info := &PolicyInfo{}
	if err := c.deleteResource(lb.Region, lb.Name, name, info); err != nil {
		return nil, err
	}

	deletedStates := []LBaaSState{LBaaSStateDeleted}
	erroredStates := []LBaaSState{LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	deleted, err := c.checkPolicyState(info, deletedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if deleted {
		return info, nil
	}
	// else poll till deleted
	info, err = c.WaitForPolicyState(lb, name, deletedStates, erroredStates, c.PollInterval, c.Timeout)
	if err != nil && client.WasNotFoundError(err) {
		// resource could not be found, thus deleted
		return nil, nil
	}
	return info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) GetPolicy(lb LoadBalancerContext, name string) (*PolicyInfo, error) {

	var info PolicyInfo
	if err := c.getResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) UpdatePolicy(lb LoadBalancerContext, name, policyType string, input *UpdatePolicyInput) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyReadyTimeout
	}

	c.ContentType = c.getContentTypeForPolicyType(policyType)
	info := &PolicyInfo{}
	if err := c.updateResource(lb.Region, lb.Name, name, &input, &info); err != nil {
		return nil, err
	}

	updatedStates := []LBaaSState{LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateModificaitonFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkPolicyState(info, updatedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return info, nil
	}
	// else poll till ready
	info, err = c.WaitForPolicyState(lb, name, updatedStates, erroredStates, c.PollInterval, c.Timeout)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// return the corrent Content Type for the Update request depending on the Policy Type
// of the Policy being updated.
func (c *PolicyClient) getContentTypeForPolicyType(policyType string) string {
	return fmt.Sprintf("application/vnd.com.oracle.oracloud.lbaas.%s+json", policyType)
}

// WaitForPolicyState waits for the resource to be in one of a set of desired states
func (c *PolicyClient) WaitForPolicyState(lb LoadBalancerContext, name string, desiredStates, errorStates []LBaaSState, pollInterval, timeoutSeconds time.Duration) (*PolicyInfo, error) {

	var getErr error
	info := &PolicyInfo{}
	err := c.client.WaitFor("Policy status update", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetPolicy(lb, name)
		if getErr != nil {
			return false, getErr
		}

		return c.checkPolicyState(info, desiredStates, errorStates)
	})
	return info, err
}

// check the State, returns in desired state (true), not ready yet (false) or errored state (error)
func (c *PolicyClient) checkPolicyState(info *PolicyInfo, desiredStates, errorStates []LBaaSState) (bool, error) {

	c.client.DebugLogString(fmt.Sprintf("Policy %v state is %v", info.Name, info.State))

	state := LBaaSState(info.State)

	if isStateInLBaaSStates(state, desiredStates) {
		// we're good, return okay
		return true, nil
	}
	if isStateInLBaaSStates(state, errorStates) {
		// not good, return error
		return false, fmt.Errorf("Policy %v in errored state %v", info.Name, info.State)
	}
	// not ready lifecycleTimeout
	return false, nil
}

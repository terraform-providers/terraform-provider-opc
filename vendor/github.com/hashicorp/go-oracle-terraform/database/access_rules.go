// Manages Access Rules for a DBaaS Service Instance.
// The only fields that can be Updated for an Access Rule is the desired state
// of the access rule. From Enabled -> Disabled.
// Deleting an Access Rule also requires an Update call, instead of a Delete API request,
// but the Operation body parameter changes from `update` to `delete`.
// All other parameters for the resource, aside from Status should be ForceNew.
// The READ function for the AccessRule resource is tricky, as there is
// no exposed `GET` function on the AccessRule API.
// There is an API endpoint to view "all" rules, however, which will be used as a
// data source to match on a supplied AccessRule name.
// Timeout only supported for the CREATE method

package database

import "time"

// API URI Paths for Container and Root objects
const (
	DBAccessContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules"
	DBAccessRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules/%s"
)

// Default Timeout value for Create
const WaitForAccessRuleTimeout = time.Duration(10 * time.Second)

// AccessRules returns a UtilityClient for managing SSH Keys and Access Rules for a DBaaS Service Instance
func (c *DatabaseClient) AccessRules() *UtilityClient {
	return &UtilityClient{
		UtilityResourceClient: UtilityResourceClient{
			DatabaseClient:   c,
			ContainerPath:    DBAccessContainerPath,
			ResourceRootPath: DBAccessRootPath,
		},
	}
}

// Status Constants for an Access Rule
type AccessRuleStatus string

const (
	AccessRuleEnabled  AccessRuleStatus = "enabled"
	AccessRuleDisabled AccessRuleStatus = "disabled"
)

// Operational Constants for either Updating/Deleting an Access Rule
type AccessRuleOperation string

const (
	AccessRuleUpdate AccessRuleOperation = "update"
	AccessRuleDelete AccessRuleOperation = "delete"
)

// Default Destination for an Access Rule
type AccessRuleDestination string

const (
	AccessRuleDefaultDestination AccessRuleDestination = "DB"
)

// Used for the GET request, as there's no direct GET request for a single Access Rule
type AccessRules struct {
	Rules []AccessRuleInfo `json:"accessRules"`
}

type AccessRuleType string

const (
	AccessRuleTypeDefault AccessRuleType = "DEFAULT"
	AccessRuleTypeSystem  AccessRuleType = "SYSTEM"
	AccessRuleTypeUser    AccessRuleType = "USER"
)

// AccessRuleInfo holds all of the known information for a single AccessRule
type AccessRuleInfo struct {
	// The Description of the Access Rule
	Description string `json:"description"`
	// The destination of the Access Rule. Should always be "DB".
	Destination AccessRuleDestination `json:"destination"`
	// The ports for the rule.
	Ports string `json:"ports"`
	// The name of the Access Rule
	Name string `json:"ruleName"`
	// The Type of the rule. One of: "DEFAULT", "SYSTEM", or "USER".
	// Computed Value
	RuleType AccessRuleType `json:"ruleType"`
	// The IP Addresses and subnets from which traffic is allowed
	Source string `json:"source"`
	// The current status of the Access Rule
	Status AccessRuleStatus `json:"status"`
}

// CreateAccessRuleInput defines the input parameters needed to create an Access Rule for a DBaaS Service Instance.
type CreateAccessRuleInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Description of the Access Rule.
	// Required
	Description string `json:"description"`
	// Destination to which traffic is allowed. Specify the value "DB".
	// Required
	Destination AccessRuleDestination `json:"destination"`
	// The network port or ports to allow traffic on. Specified as a single port or a range.
	// Required
	Ports string `json:"ports"`
	// The name of the Access Rule
	// Required
	Name string `json:"ruleName"`
	// The IP addresses and subnets from which traffic is allowed.
	// Valid values are:
	//   - "DB" for any other cloud service instance in the service instances `ora_db` security list
	//   - "PUBLIC-INTERNET" for any host on the internet.
	//   - A single IP address or comma-separated list of subnets (in CIDR format) or IPv4 addresses.
	// Required
	Source string `json:"source"`
	// Desired Status of the rule. Either "disabled" or "enabled".
	// Required
	Status AccessRuleStatus `json:"status"`
	// Time to wait for an access rule to be ready
	Timeout time.Duration `json:"-"`
}

// Creates an AccessRule with the supplied input struct.
// The API call to Create returns a nil body object, and a 202 status code on success.
// Thus, the Create method will return the resulting object from an internal GET call
// during the WaitForReady timeout.
func (c *UtilityClient) CreateAccessRule(input *CreateAccessRuleInput) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRule AccessRuleInfo
	if err := c.createResource(input, &accessRule); err != nil {
		return nil, err
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForAccessRuleTimeout
	}

	getInput := &GetAccessRuleInput{
		Name: input.Name,
	}

	result, err := c.WaitForAccessRuleReady(getInput, timeout)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccessRuleInput defines the input parameters needed to retrieve information
// on an AccessRule for a DBaas Service Instance.
type GetAccessRuleInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule.
	// Because there is no native "GET" to return a single AccessRuleInfo object, we don't
	// need to marshal a request body for the GET request. Instead the request returns a slice
	// of AccessRuleInfo structs, which we iterate on to interpret the desired AccessRuleInfo struct
	// Required
	Name string `json:"-"`
}

// Get's a slice of every AccessRule, and iterates on the result until
// we find the correctly matching access rule. This is likely an expensive operation depending
// on how many access rules the customer has. However, since there's no direct GET API endpoint
// for a single Access Rule, it's not able to be optimized yet.
func (c *UtilityClient) GetAccessRule(input *GetAccessRuleInput) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRules
	if err := c.getResource("", &accessRules); err != nil {
		return nil, err
	}

	// This is likely not the most optimal path for this, however, the upper bound on
	// performance here is the actual API request, not the iteration.
	for _, rule := range accessRules.Rules {
		if rule.Name == input.Name {
			return &rule, nil
		}
	}

	// Iterated through entire slice, rule was not found.
	// No error occured though, return a nil struct, and allow the Provdier to handle
	// a Nil response case.
	return nil, nil
}

// UpdateAccessRuleInput defines the Update parameters needed to update an AccessRule
// for a DBaaS Service Instance.
type UpdateAccessRuleInput struct {
	// Name of the DBaaS Service Instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule. Used in the request's URI, not as a body parameter.
	// Required
	Name string `json:"-"`
	// Type of Operation being performed. This should never be set in the Provider,
	// as we're explicitly calling an Update function here, so the SDK uses the constant
	// defined for Updating an AccessRule
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// Desired Status of the Access Rule. This is the only attribute that can actually be
	// modified on an access rule.
	// Required
	Status AccessRuleStatus `json:"status"`
}

// Updates an AccessRule with the provided input struct. Returns a fully populated Info struct
// and any errors encountered
func (c *UtilityClient) UpdateAccessRule(input *UpdateAccessRuleInput,
) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleUpdate
	// Initialize the response struct
	var accessRule AccessRuleInfo
	if err := c.updateResource(input.Name, input, &accessRule); err != nil {
		return nil, err
	}
	return &accessRule, nil
}

// DeleteAccessRuleInput defines the Delete parameters needed to delete an AccessRule
// for a DBaaS Service Instance. There's no dedicated DELETE method on the API, so this
// mimics the same behavior of the Update method, but using the Delete operational constant.
// Instead of implementing, choosing to be verbose here for ease of use in the Provider, and clarity.
type DeleteAccessRuleInput struct {
	// Name of the DBaaS Service Instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule. Used in the request's URI, not as a body parameter.
	// Required
	Name string `json:"-"`
	// Type of Operation being performed. This should never be set in the Provider,
	// as we're explicitly calling an Delete function here, so the SDK uses the constant
	// defined for Deleting an AccessRule
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// Desired Status of the Access Rule. This is the only attribute that can actually be
	// modified on an access rule.
	// Required
	Status AccessRuleStatus `json:"status"`
}

// Deletes an AccessRule with the provided input struct. Returns any errors that occurred.
func (c *UtilityClient) DeleteAccessRule(input *DeleteAccessRuleInput) error {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleDelete
	// The Update API call with a `DELETE` operation actually returns the same access rule info
	// in a response body. As we are deleting the AccessRule, we don't actually need to parse that.
	// However, the Update API call requires a pointer to parse, or else we throw an error during the
	// json unmarshal
	var result AccessRuleInfo
	if err := c.updateResource(input.Name, input, &result); err != nil {
		return err
	}
	return nil
}

func (c *UtilityClient) WaitForAccessRuleReady(input *GetAccessRuleInput, timeout time.Duration) (*AccessRuleInfo, error) {
	var info *AccessRuleInfo
	var getErr error
	err := c.client.WaitFor("access rule to be ready", timeout, func() (bool, error) {
		info, getErr = c.GetAccessRule(input)
		if getErr != nil {
			return false, getErr
		}
		if info != nil {
			// Rule found, return. Desired case
			return true, nil
		}
		// Rule not found, wait
		return false, nil
	})
	return info, err
}

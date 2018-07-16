package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceLBaaSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// ForceNew when changing parent load_balancer
			if diff.HasChange("load_balancer") {
				diff.ForceNew("load_balancer")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"load_balancer": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateLoadBalancerID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateLoadBalancerPolicyName,
			},
			"application_cookie_stickiness_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"cloudgate_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudgate_application": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cloudgate_policy_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"identity_service_instance_guid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"virtual_hostname_for_policy_attribution": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"load_balancer_cookie_stickiness_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_expiration_period": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"load_balancing_mechanism_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancing_mechanism": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"round_robin",
								"least_conn",
								"ip_hash",
							}, true),
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"rate_limiting_request_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"burst_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"delay_excessive_requests": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"http_error_code": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      503,
							ValidateFunc: validation.IntBetween(405, 599),
						},
						"logging_level": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "warn",
							ValidateFunc: validation.StringInSlice([]string{
								"info",
								"notice",
								"warn",
								"error",
							}, true),
						},
						"rate_limiting_criteria": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Default:  "server",
							ValidateFunc: validation.StringInSlice([]string{
								"server",
								"remote_address",
								"host",
							}, true),
						},
						"requests_per_second": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"zone_memory_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
						},
						"zone": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateLoadBalancerPolicyName,
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"redirect_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"response_code": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(300, 399),
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"resource_access_control_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disposition": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"DENY_ALL",
								"ALLOW_ALL",
							}, true),
						},
						"denied_clients": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"permitted_clients": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"set_request_header_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action_when_header_exists": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"NOOP",
								"PREPEND",
								"APPEND",
								"OVERWRITE",
								"CLEAR",
							}, true),
						},
						"action_when_header_value_is": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"action_when_header_value_is_not": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"ssl_negotiation_policy",
					"trusted_certificate_policy",
				},
			},
			"ssl_negotiation_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"server_order_preference": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
						},
						"ssl_protocol": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ssl_ciphers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"trusted_certificate_policy",
				},
			},
			"trusted_certificate_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trusted_certificate": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{
					"application_cookie_stickiness_policy",
					"cloudgate_policy",
					"load_balancer_cookie_stickiness_policy",
					"load_balancing_mechanism_policy",
					"rate_limiting_request_policy",
					"redirect_policy",
					"resource_access_control_policy",
					"set_request_header_policy",
					"ssl_negotiation_policy",
				},
			},

			// Read only attributes
			"state": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	policyClient := lbaasClient.PolicyClient()
	lb := getLoadBalancerContextFromID(d.Get("load_balancer").(string))

	input := lbaas.CreatePolicyInput{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("application_cookie_stickiness_policy"); ok {
		input.Type = "AppCookieStickinessPolicy"
		input.ApplicationCookieStickinessPolicyInfo = expandApplicationCookieStickinessPolicy(d)
	}
	if _, ok := d.GetOk("cloudgate_policy"); ok {
		input.Type = "CloudGatePolicy"
		input.CloudGatePolicyInfo = expandCloudGatePolicy(d)
	}
	if _, ok := d.GetOk("load_balancer_cookie_stickiness_policy"); ok {
		input.Type = "LBCookieStickinessPolicy"
		input.LoadBalancerCookieStickinessPolicyInfo = expandLoadBalancerCookieStickinessPolicy(d)
	}
	if _, ok := d.GetOk("load_balancing_mechanism_policy"); ok {
		input.Type = "LoadBalancingMechanismPolicy"
		input.LoadBalancingMechanismPolicyInfo = expandLoadBalancingMechanismPolicy(d)
	}
	if _, ok := d.GetOk("rate_limiting_request_policy"); ok {
		input.Type = "RateLimitingRequestPolicy"
		input.RateLimitingRequestPolicyInfo = expandRateLimitingRequestPolicy(d)
	}
	if _, ok := d.GetOk("redirect_policy"); ok {
		input.Type = "RedirectPolicy"
		input.RedirectPolicyInfo = expandRedirectPolicy(d)
	}
	if _, ok := d.GetOk("resource_access_control_policy"); ok {
		input.Type = "ResourceAccessControlPolicy"
		input.ResourceAccessControlPolicyInfo = expandResourceAccessControlPolicy(d)
	}
	if _, ok := d.GetOk("set_request_header_policy"); ok {
		input.Type = "SetRequestHeaderPolicy"
		input.SetRequestHeaderPolicyInfo = expandSetRequestHeaderPolicy(d)
	}
	if _, ok := d.GetOk("ssl_negotiation_policy"); ok {
		input.Type = "SSLNegotiationPolicy"
		input.SSLNegotiationPolicyInfo = exapndSSLNegotiationPolicy(d)
	}
	if _, ok := d.GetOk("trusted_certificate_policy"); ok {
		input.Type = "TrustedCertPolicy"
		input.TrustedCertificatePolicyInfo = expandTrustedCertificatePolicy(d)
	}

	info, err := policyClient.CreatePolicy(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer Policy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	policyClient := lbaasClient.PolicyClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	result, err := policyClient.GetPolicy(lb, name)
	if err != nil {
		// Policy does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Load Balancer Policy %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("state", result.State)
	d.Set("type", result.Type)
	d.Set("uri", result.URI)
	d.Set("load_balancer", fmt.Sprintf("%s/%s", lb.Region, lb.Name))

	if result.Type == "AppCookieStickinessPolicy" {
		flattenApplicationCookieStickinessPolicy(d, result)
	}
	if result.Type == "CloudGatePolicy" {
		flattenCloudGatePolicy(d, result)
	}
	if result.Type == "LBCookieStickinessPolicy" {
		flattenLoadBalancerCookieStickinessPolicy(d, result)
	}
	if result.Type == "LoadBalancingMechanismPolicy" {
		flattenLoadBalancingMechanismPolicy(d, result)
	}
	if result.Type == "RateLimitingRequestPolicy" {
		flattenRateLimitingRequestPolicy(d, result)
	}
	if result.Type == "ResourceAccessControlPolicy" {
		flattenResourceAccessControlPolicy(d, result)
	}
	if result.Type == "SetRequestHeaderPolicy" {
		flattenSetRequestHeaderPolicy(d, result)
	}
	if result.Type == "SSLNegotiationPolicy" {
		flattenSSLNegotiationPolicy(d, result)
	}
	if result.Type == "TrustedCertPolicy" {
		flattenTrustedCertificatePolicy(d, result)
	}

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	policyClient := lbaasClient.PolicyClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	input := lbaas.UpdatePolicyInput{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("application_cookie_stickiness_policy"); ok {
		input.Type = "AppCookieStickinessPolicy"
		input.ApplicationCookieStickinessPolicyInfo = expandApplicationCookieStickinessPolicy(d)
	}
	if _, ok := d.GetOk("cloudgate_policy"); ok {
		input.Type = "CloudGatePolicy"
		input.CloudGatePolicyInfo = expandCloudGatePolicy(d)
	}
	if _, ok := d.GetOk("load_balancer_cookie_stickiness_policy"); ok {
		input.Type = "LBCookieStickinessPolicy"
		input.LoadBalancerCookieStickinessPolicyInfo = expandLoadBalancerCookieStickinessPolicy(d)
	}
	if _, ok := d.GetOk("load_balancing_mechanism_policy"); ok {
		input.Type = "LoadBalancingMechanismPolicy"
		input.LoadBalancingMechanismPolicyInfo = expandLoadBalancingMechanismPolicy(d)
	}
	if _, ok := d.GetOk("rate_limiting_request_policy"); ok {
		input.Type = "RateLimitingRequestPolicy"
		input.RateLimitingRequestPolicyInfo = expandRateLimitingRequestPolicy(d)
	}
	if _, ok := d.GetOk("redirect_policy"); ok {
		input.Type = "RedirectPolicy"
		input.RedirectPolicyInfo = expandRedirectPolicy(d)
	}
	if _, ok := d.GetOk("resource_access_control_policy"); ok {
		input.Type = "ResourceAccessControlPolicy"
		input.ResourceAccessControlPolicyInfo = expandResourceAccessControlPolicy(d)
	}
	if _, ok := d.GetOk("set_request_header_policy"); ok {
		input.Type = "SetRequestHeaderPolicy"
		input.SetRequestHeaderPolicyInfo = expandSetRequestHeaderPolicy(d)
	}
	if _, ok := d.GetOk("ssl_negotiation_policy"); ok {
		input.Type = "SSLNegotiationPolicy"
		input.SSLNegotiationPolicyInfo = exapndSSLNegotiationPolicy(d)
	}
	if _, ok := d.GetOk("trusted_certificate_policy"); ok {
		input.Type = "TrustedCertPolicy"
		input.TrustedCertificatePolicyInfo = expandTrustedCertificatePolicy(d)
	}

	result, err := policyClient.UpdatePolicy(lb, name, input.Type, &input)
	if err != nil {
		return fmt.Errorf("Error updating Policy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, result.Name))

	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	policyClient := lbaasClient.PolicyClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := policyClient.DeletePolicy(lb, name); err != nil {
		return fmt.Errorf("Error deleting Policy: %v", err)
	}
	return nil
}

// ApplicationCookieStickinessPolicy

func expandApplicationCookieStickinessPolicy(d *schema.ResourceData) lbaas.ApplicationCookieStickinessPolicyInfo {

	policy := d.Get("application_cookie_stickiness_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.ApplicationCookieStickinessPolicyInfo{
		AppCookieName: policy["cookie_name"].(string),
	}

	return info
}

func flattenApplicationCookieStickinessPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("application_cookie_stickiness_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["cookie_name"] = result.AppCookieName

	p = append(p, attrs)

	d.Set("application_cookie_stickiness_policy", p)
	return nil
}

// CloudGatePolicy

func expandCloudGatePolicy(d *schema.ResourceData) lbaas.CloudGatePolicyInfo {

	policy := d.Get("cloudgate_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.CloudGatePolicyInfo{
		CloudGateApplication:                policy["cloudgate_application"].(string),
		CloudGatePolicyName:                 policy["cloudgate_policy_name"].(string),
		IdentityServiceInstanceGuid:         policy["identity_service_instance_guid"].(string),
		VirtualHostnameForPolicyAttribution: policy["virtual_hostname_for_policy_attribution"].(string),
	}

	return info
}

func flattenCloudGatePolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("cloudgate_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["cloudgate_application"] = result.CloudGateApplication
	attrs["cloudgate_policy_name"] = result.CloudGatePolicyName
	attrs["identity_service_instance_guid"] = result.IdentityServiceInstanceGuid
	attrs["virtual_hostname_for_policy_attribution"] = result.VirtualHostnameForPolicyAttribution

	p = append(p, attrs)

	d.Set("cloudgate_policy", p)
	return nil
}

// LoadBalancerCookieStickinessPolicy

func expandLoadBalancerCookieStickinessPolicy(d *schema.ResourceData) lbaas.LoadBalancerCookieStickinessPolicyInfo {

	policy := d.Get("load_balancer_cookie_stickiness_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.LoadBalancerCookieStickinessPolicyInfo{
		CookieExpirationPeriod: policy["cookie_expiration_period"].(int),
	}

	return info
}

func flattenLoadBalancerCookieStickinessPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("load_balancer_cookie_stickiness_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["cookie_expiration_period"] = result.CookieExpirationPeriod

	p = append(p, attrs)

	d.Set("load_balancer_cookie_stickiness_policy", p)
	return nil
}

// LoadBalancingMechanismPolicy

func expandLoadBalancingMechanismPolicy(d *schema.ResourceData) lbaas.LoadBalancingMechanismPolicyInfo {

	policy := d.Get("load_balancing_mechanism_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.LoadBalancingMechanismPolicyInfo{
		LoadBalancingMechanism: policy["load_balancing_mechanism"].(string),
	}

	return info
}

func flattenLoadBalancingMechanismPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("load_balancing_mechanism_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["load_balancing_mechanism"] = result.LoadBalancingMechanism

	p = append(p, attrs)

	d.Set("load_balancing_mechanism_policy", p)
	return nil
}

// RateLimitingRequestPolicy

func expandRateLimitingRequestPolicy(d *schema.ResourceData) lbaas.RateLimitingRequestPolicyInfo {

	policy := d.Get("rate_limiting_request_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.RateLimitingRequestPolicyInfo{
		BurstSize:                   policy["burst_size"].(int),
		DoNotDelayExcessiveRequests: !policy["delay_excessive_requests"].(bool),
		HttpStatusErrorCode:         policy["http_error_code"].(int),
		LogLevel:                    policy["logging_level"].(string),
		RateLimitingCriteria:        policy["rate_limiting_criteria"].(string),
		RequestsPerSecond:           policy["requests_per_second"].(int),
		StorageSize:                 policy["zone_memory_size"].(int),
		Zone:                        policy["zone"].(string),
	}

	return info
}

func flattenRateLimitingRequestPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("rate_limiting_request_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["burst_size"] = result.BurstSize
	attrs["delay_excessive_requests"] = !result.DoNotDelayExcessiveRequests
	attrs["http_error_code"] = result.HttpStatusErrorCode
	attrs["log_level"] = result.LogLevel
	attrs["rate_limiting_criteria"] = result.RateLimitingCriteria
	attrs["requests_per_second"] = result.RequestsPerSecond
	attrs["storage_size"] = result.StorageSize
	attrs["zone"] = result.Zone

	p = append(p, attrs)

	d.Set("rate_limiting_request_policy", p)
	return nil
}

// RedirectPolicy

func expandRedirectPolicy(d *schema.ResourceData) lbaas.RedirectPolicyInfo {

	policy := d.Get("redirect_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.RedirectPolicyInfo{
		RedirectURI:  policy["redirect_uri"].(string),
		ResponseCode: policy["response_code"].(int),
	}

	return info
}

func flattenRedirectPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("redirect_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["redirect_uri"] = result.RedirectURI
	attrs["response_code"] = result.ResponseCode

	p = append(p, attrs)

	d.Set("redirect_policy", p)
	return nil
}

// ResourceAccessControlPolicy

func expandResourceAccessControlPolicy(d *schema.ResourceData) lbaas.ResourceAccessControlPolicyInfo {

	policy := d.Get("resource_access_control_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.ResourceAccessControlPolicyInfo{
		Disposition: policy["disposition"].(string),
	}

	deniedClients := getStringSet(d, "resource_access_control_policy.0.denied_clients")
	if len(deniedClients) != 0 {
		info.DeniedClients = deniedClients
	}
	permittedClients := getStringSet(d, "resource_access_control_policy.0.permitted_clients")
	if len(permittedClients) != 0 {
		info.PermittedClients = permittedClients
	}

	return info
}

func flattenResourceAccessControlPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("resource_access_control_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["disposition"] = result.Disposition
	attrs["denied_clients"] = result.DeniedClients
	attrs["permitted_clients"] = result.PermittedClients

	p = append(p, attrs)

	d.Set("resource_access_control_policy", p)
	return nil
}

// SetRequestHeaderPolicy

func expandSetRequestHeaderPolicy(d *schema.ResourceData) lbaas.SetRequestHeaderPolicyInfo {

	policy := d.Get("set_request_header_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.SetRequestHeaderPolicyInfo{
		HeaderName: policy["header_name"].(string),
	}

	if val, ok := policy["value"].(string); ok && val != "" {
		info.Value = val
	}
	if val, ok := policy["action_when_header_exists"].(string); ok && val != "" {
		info.ActionWhenHeaderExists = val
	}

	info.ActionWhenHeaderValueIs = getStringSet(d, "set_request_header_policy.0.action_when_header_value_is")
	info.ActionWhenHeaderValueIsNot = getStringSet(d, "set_request_header_policy.0.action_when_header_value_is_not")

	return info
}

func flattenSetRequestHeaderPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("set_request_header_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["header_name"] = result.HeaderName
	attrs["value"] = result.Value
	attrs["action_when_header_exists"] = result.ActionWhenHeaderExists
	attrs["action_when_header_value_is"] = result.ActionWhenHeaderValueIs
	attrs["action_when_header_value_is_not"] = result.ActionWhenHeaderValueIsNot

	p = append(p, attrs)

	return d.Set("set_request_header_policy", p)
}

// SSLNegotiationPolicy

func exapndSSLNegotiationPolicy(d *schema.ResourceData) lbaas.SSLNegotiationPolicyInfo {

	policy := d.Get("ssl_negotiation_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.SSLNegotiationPolicyInfo{
		Port: policy["port"].(int),
	}

	if val, ok := policy["server_order_preference"].(string); ok && val != "" {
		info.ServerOrderPreference = val
	}

	sslProtocol := getStringSet(d, "ssl_negotiation_policy.0.ssl_protocol")
	if len(sslProtocol) != 0 {
		info.SSLProtocol = sslProtocol
	}
	sslCiphers := getStringSet(d, "ssl_negotiation_policy.0.ssl_ciphers")
	if len(sslCiphers) != 0 {
		info.SSLCiphers = sslCiphers
	}

	return info
}

func flattenSSLNegotiationPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("ssl_negotiation_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["port"] = result.Port
	attrs["server_order_preference"] = result.Value
	attrs["action_when_header_exists"] = result.ActionWhenHeaderExists
	attrs["ssl_protocol"] = result.SSLProtocol
	attrs["ssl_ciphers"] = result.SSLCiphers

	p = append(p, attrs)

	d.Set("ssl_negotiation_policy", p)
	return nil
}

// TrustedCertificatePolicy

func expandTrustedCertificatePolicy(d *schema.ResourceData) lbaas.TrustedCertificatePolicyInfo {

	policy := d.Get("trusted_certificate_policy").([]interface{})[0].(map[string]interface{})

	info := lbaas.TrustedCertificatePolicyInfo{
		TrustedCertificate: policy["trusted_certificate"].(string),
	}

	return info
}

func flattenTrustedCertificatePolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	attrs := make(map[string]interface{})
	p := make([]map[string]interface{}, 0)
	val, ok := d.GetOk("trusted_certificate_policy")
	if ok {
		policyConfiguration := val.([]interface{})
		attrs = policyConfiguration[0].(map[string]interface{})
	}

	attrs["trusted_certificate"] = result.TrustedCertificate

	p = append(p, attrs)

	d.Set("trusted_certificate_policy", p)
	return nil
}

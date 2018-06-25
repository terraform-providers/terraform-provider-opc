package opc

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: map[string]*schema.Schema{
			"load_balancer": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"set_request_header_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": {
							Type:     schema.TypeString,
							Required: true,
							// TODO Force New?
							// TODO Add validation
						},
						"action_when_header_exists": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_when_header_value_is": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"action_when_header_value_is_not": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
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

	client := meta.(*Client).lbaasClient.PolicyClient()

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		s := strings.Split(loadBalancer.(string), "/")
		lb = lbaas.LoadBalancerContext{
			Region: s[0],
			Name:   s[1],
		}
	}

	input := lbaas.CreatePolicyInput{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("set_request_header_policy"); ok {
		input.Type = "SetRequestHeaderPolicy"
		input.SetRequestHeaderPolicyInfo = expandSetRequestHeaderPolicy(d)
	}

	info, err := client.CreatePolicy(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.PolicyClient()
	name := getLastNameInURIPath(d.Id())

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		s := strings.Split(loadBalancer.(string), "/")
		lb = lbaas.LoadBalancerContext{
			Region: s[0],
			Name:   s[1],
		}
	}

	result, err := lbaasClient.GetPolicy(lb, name)
	if err != nil {
		// Policy does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Server Pool %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("state", result.State)
	d.Set("type", result.Type)
	d.Set("uri", result.URI)

	if result.Type == "SetRequestHeaderPolicy" {
		flattenSetRequestHeaderPolicy(d, result)
	}

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.PolicyClient()
	name := getLastNameInURIPath(d.Id())

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		s := strings.Split(loadBalancer.(string), "/")
		lb = lbaas.LoadBalancerContext{
			Region: s[0],
			Name:   s[1],
		}
	}

	input := lbaas.UpdatePolicyInput{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("set_request_header_policy"); ok {
		input.Type = "SetRequestHeaderPolicy"
		input.SetRequestHeaderPolicyInfo = expandSetRequestHeaderPolicy(d)
	}

	result, err := lbaasClient.UpdatePolicy(lb, name, input.Type, &input)
	if err != nil {
		return fmt.Errorf("Error updating Policy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, result.Name))

	// TODO instead of re-read, process info from UpdatePolicy()
	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.PolicyClient()
	name := d.Id()

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		s := strings.Split(loadBalancer.(string), "/")
		lb = lbaas.LoadBalancerContext{
			Region: s[0],
			Name:   s[1],
		}
	}

	if _, err := lbaasClient.DeletePolicy(lb, name); err != nil {
		return fmt.Errorf("Error deleting Policy")
	}
	return nil
}

func expandSetRequestHeaderPolicy(d *schema.ResourceData) lbaas.SetRequestHeaderPolicyInfo {
	// there can be only one
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
	if val, ok := policy["action_when_header_value_is"].(string); ok && val != "" {
		info.ActionWhenHeaderValueIs = getStringList(d, "set_request_header_policy.0.action_when_header_value_is")
	}
	if val, ok := policy["action_when_header_value_is_not"].(string); ok && val != "" {
		info.ActionWhenHeaderValueIsNot = getStringList(d, "set_request_header_policy.0.action_when_header_value_is_not")
	}
	return info
}

func flattenSetRequestHeaderPolicy(d *schema.ResourceData, result *lbaas.PolicyInfo) error {
	val, _ := d.GetOk("set_request_header_policy")

	p := make([]map[string]interface{}, 0)
	policyConfiguration := val.([]interface{})
	attrs := policyConfiguration[0].(map[string]interface{})

	if len(policyConfiguration) != 1 {
		return fmt.Errorf("Invalid Policy Configuration info")
	}

	attrs["header_name"] = result.HeaderName
	attrs["vale"] = result.Value
	attrs["action_when_header_exists"] = result.ActionWhenHeaderExists
	attrs["action_when_header_value_is"] = setStringList(d, "action_when_header_value_is", result.ActionWhenHeaderValueIs)
	attrs["action_when_header_value_is_not"] = setStringList(d, "action_when_header_value_is_not", result.ActionWhenHeaderValueIsNot)

	p = append(p, attrs)

	d.Set("set_request_header_policy", p)
	return nil
}

package opc

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceLBaaSLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCLoadBalancerCreate,
		Read:   resourceOPCLoadBalancerRead,
		Update: resourceOPCLoadBalancerUpdate,
		Delete: resourceOPCLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateLoadBalancerResourceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ip_network": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateComputeResourceFQDN,
			},
			"parent_load_balancer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permitted_clients": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"permitted_methods": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scheme": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"INTERNET_FACING",
					"INTERNAL",
				}, true),
			},
			"server_pool": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOriginServerPoolURI,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Read only attributes
			"balancer_vips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"canonical_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudgate_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	lbClient := meta.(*Client).lbaasClient.LoadBalancerClient()
	input := lbaas.CreateLoadBalancerInput{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		input.Disabled = getDisabledStateKeyword(enabled.(bool))
	}

	if ipNetwork, ok := d.GetOk("ip_network"); ok {
		input.IPNetworkName = ipNetwork.(string)
	}

	if scheme, ok := d.GetOk("scheme"); ok {
		input.Scheme = lbaas.LoadBalancerScheme(scheme.(string))
	}

	if serverPool, ok := d.GetOk("server_pool"); ok {
		input.OriginServerPool = serverPool.(string)
	}

	if parent, ok := d.GetOk("parent_load_balancer"); ok {
		input.ParentLoadBalancer = parent.(string)
	}

	permittedClients := getStringSet(d, "permitted_clients")
	if len(permittedClients) != 0 {
		input.PermittedClients = permittedClients
	}

	permittedMethods := getStringSet(d, "permitted_methods")
	if len(permittedMethods) != 0 {
		input.PermittedMethods = permittedMethods
	}

	policies := getStringSet(d, "policies")
	if len(policies) != 0 {
		input.Policies = policies
	}

	tags := getStringSet(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	info, err := lbClient.CreateLoadBalancer(&input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", info.Region, info.Name))
	return resourceOPCLoadBalancerRead(d, meta)
}

func resourceOPCLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	lbClient := meta.(*Client).lbaasClient.LoadBalancerClient()

	s := strings.Split(d.Id(), "/")
	lb := lbaas.LoadBalancerContext{
		Region: s[0],
		Name:   s[1],
	}

	result, err := lbClient.GetLoadBalancer(lb)
	if err != nil {
		// LoadBalancer does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Load Balancer %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("canonical_host_name", result.CanonicalHostName)
	d.Set("cloudgate_capable", result.CloudgateCapable)
	d.Set("enabled", getEnabledState(result.Disabled))
	d.Set("description", result.Description)
	d.Set("ip_network", result.IPNetworkName)
	d.Set("name", result.Name)
	d.Set("parent_load_balancer", result.ParentLoadBalancer)
	d.Set("region", result.Region)
	d.Set("scheme", result.Scheme)
	d.Set("server_pool", result.OriginServerPool)
	d.Set("uri", result.URI)

	if err := setStringList(d, "balancer_vips", result.BalancerVIPs); err != nil {
		return err
	}
	if err := setStringList(d, "permitted_clients", result.PermittedClients); err != nil {
		return err
	}
	if err := setStringList(d, "permitted_methods", result.PermittedMethods); err != nil {
		return err
	}
	if err := setStringList(d, "policies", result.Policies); err != nil {
		return err
	}
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	lbClient := meta.(*Client).lbaasClient.LoadBalancerClient()
	lb := getLoadBalancerContextFromID(d.Id())

	input := lbaas.UpdateLoadBalancerInput{
		Name: d.Get("name").(string),
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		input.Disabled = getDisabledStateKeyword(enabled.(bool))
	}

	input.Description = updateOrRemoveStringAttribute(d, "description")
	input.IPNetworkName = updateOrRemoveStringAttribute(d, "ip_network")
	input.OriginServerPool = updateOrRemoveStringAttribute(d, "server_pool")
	input.ParentLoadBalancer = updateOrRemoveStringAttribute(d, "parent_load_balancer")

	input.PermittedClients = updateOrRemoveStringListAttribute(d, "permitted_clients")
	input.PermittedMethods = updateOrRemoveStringListAttribute(d, "permitted_methods")
	input.Policies = updateOrRemoveStringListAttribute(d, "policies")
	input.Tags = updateOrRemoveStringListAttribute(d, "tags")

	result, err := lbClient.UpdateLoadBalancer(lb, &input)
	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", lb.Region, result.Name))

	return resourceOPCLoadBalancerRead(d, meta)
}

func resourceOPCLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	lbClient := meta.(*Client).lbaasClient.LoadBalancerClient()
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := lbClient.DeleteLoadBalancer(lb); err != nil {
		return fmt.Errorf("Error deleting LoadBalancer")
	}
	return nil
}

// return the Disbaled State keyword for Load Balancer enabled state
func getDisabledStateKeyword(enabled bool) lbaas.LBaaSDisabled {
	if enabled {
		return lbaas.LBaaSDisabledFalse
	} else {
		return lbaas.LBaaSDisabledTrue
	}
}

// convert the DisabledState attribute to a boolean representing the enabled state
func getEnabledState(state lbaas.LBaaSDisabled) bool {
	if state == lbaas.LBaaSDisabledFalse {
		return true
	}
	return false
}

// return the changed value, empty string if the attribute has been removed or nil if unchanged
// and optional string trnasformation function to apply to the returned value
func updateOrRemoveStringAttribute(d *schema.ResourceData, attributeName string) *string {
	if d.HasChange(attributeName) {
		val := ""
		if attribute, ok := d.GetOk(attributeName); ok {
			val = attribute.(string)
		}
		return &val
	}
	return nil
}

// return the updated list, empty list if attribute has been removed, or nil if unchanged
func updateOrRemoveStringListAttribute(d *schema.ResourceData, attributeName string) *[]string {
	if d.HasChange(attributeName) {
		val := getStringSet(d, attributeName)
		if val == nil {
			// return an empty list
			val = []string{}
		}
		return &val
	}
	return nil
}

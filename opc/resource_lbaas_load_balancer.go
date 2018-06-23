package opc

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // TODO name can be changed
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				// TODO separate enabled flag (desired state) from DisabledState (current state)
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// TODO add validation for 3 part name
				// TODO add valication only supported for INTERNAL load balancer?
			},
			"premitted_methods": {
				Type:     schema.TypeList, // TODO TypeSet? API returns ordered list
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// TODO add validation
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
				// TODO add validation
			},
			"server_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList, // TODO TypeSet?
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

	client := meta.(*Client).lbaasClient.LoadBalancerClient()
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
		input.ServerPool = serverPool.(string)
	}

	permittedMethods := getStringList(d, "premitted_methods")
	if len(permittedMethods) != 0 {
		input.PermittedMethods = permittedMethods
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	info, err := client.CreateLoadBalancer(&input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer: %s", err)
	}

	d.SetId(info.Region + "/" + info.Name)
	return resourceOPCLoadBalancerRead(d, meta)
}

func resourceOPCLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.LoadBalancerClient()

	s := strings.Split(d.Id(), "/")
	lb := lbaas.LoadBalancerContext{
		Region: s[0],
		Name:   s[1],
	}

	result, err := lbaasClient.GetLoadBalancer(lb)
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
	d.Set("region", result.Region)
	d.Set("scheme", result.Scheme)
	d.Set("server_pool", result.ServerPool)
	d.Set("uri", result.URI)

	if err := setStringList(d, "balancer_vips", result.BalancerVIPs); err != nil {
		return err
	}
	if err := setStringList(d, "premitted_methods", result.PermittedMethods); err != nil {
		return err
	}
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.LoadBalancerClient()

	s := strings.Split(d.Id(), "/")
	lb := lbaas.LoadBalancerContext{
		Region: s[0],
		Name:   s[1],
	}

	input := lbaas.UpdateLoadBalancerInput{
		Name: d.Get("name").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		input.Disabled = getDisabledStateKeyword(enabled.(bool))
	}

	// TODO API complains
	// * opc_lbaas_load_balancer.lb1: Error updating LoadBalancer: 400: ip_network_name should not be specified while updating lb1.
	// if ipNetwork, ok := d.GetOk("ip_network"); ok {
	// 	input.IPNetworkName = ipNetwork.(string)
	// }

	if serverPool, ok := d.GetOk("server_pool"); ok {
		input.ServerPool = serverPool.(string)
	}

	permittedMethods := getStringList(d, "premitted_methods")
	if len(permittedMethods) != 0 {
		input.PermittedMethods = permittedMethods
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	result, err := lbaasClient.UpdateLoadBalancer(lb, &input)
	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer: %s", err)
	}

	d.SetId(result.Region + "/" + result.Name)

	// TODO instead of re-read, process info from UpdateLoadBalancer()
	return resourceOPCLoadBalancerRead(d, meta)
}

func resourceOPCLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.LoadBalancerClient()

	s := strings.Split(d.Id(), "/")
	lb := lbaas.LoadBalancerContext{
		Region: s[0],
		Name:   s[1],
	}

	if _, err := lbaasClient.DeleteLoadBalancer(lb); err != nil {
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

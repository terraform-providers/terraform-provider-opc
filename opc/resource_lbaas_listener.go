package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLBaaSListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerCreate,
		Read:   resourceListenerRead,
		Update: resourceListenerUpdate,
		Delete: resourceListenerDelete,
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
			"balancer_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"server_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssl_certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"virtual_hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Read only attributes
			"inline_policies": {
				// TODO not returned from API, Remove?
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"operation_details": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_listener": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
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

func resourceListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).lbaasClient.ListenerClient()

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		lb = getLoadBalancerContextFromID(loadBalancer.(string))
	}

	input := lbaas.CreateListenerInput{
		Name:                 d.Get("name").(string),
		BalancerProtocol:     lbaas.Protocol(d.Get("balancer_protocol").(string)),
		OriginServerProtocol: lbaas.Protocol(d.Get("server_protocol").(string)),
		Port:                 d.Get("port").(int),
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		input.Disabled = getDisabledStateKeyword(enabled.(bool))
	}

	if serverPool, ok := d.GetOk("server_pool"); ok {
		// Only the URI Path is need on Create
		input.OriginServerPool = getURIRequestPath(serverPool.(string))
	}

	pathPrefixes := getStringList(d, "path_prefixes")
	if len(pathPrefixes) != 0 {
		input.PathPrefixes = pathPrefixes
	}

	policies := getStringList(d, "policies")
	if len(policies) != 0 {
		input.Policies = policies
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	virtualHosts := getStringList(d, "virtual_hosts")
	if len(virtualHosts) != 0 {
		input.VirtualHosts = virtualHosts
	}

	info, err := client.CreateListener(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer Listener: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourceListenerRead(d, meta)
}

func resourceListenerRead(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.ListenerClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	result, err := lbaasClient.GetListener(lb, name)
	if err != nil {
		// Listener does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Load Balancer Listener %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("balancer_protocol", result.BalancerProtocol)
	d.Set("enabled", getEnabledState(result.Disabled))
	d.Set("name", result.Name)
	d.Set("operation_details", result.OperationDetails)
	d.Set("parent_listener", result.ParentListener)
	d.Set("port", result.Port)
	d.Set("server_pool", result.OriginServerPool)
	d.Set("server_protocol", result.OriginServerProtocol)
	d.Set("state", result.State)
	d.Set("uri", result.URI)
	d.Set("load_balancer", fmt.Sprintf("%s/%s", lb.Region, lb.Name))

	if err := setStringList(d, "inline_policies", result.InlinePolicies); err != nil {
		return err
	}

	if err := setStringList(d, "path_prefixes", result.PathPrefixes); err != nil {
		return err
	}

	if err := setStringList(d, "policies", result.Policies); err != nil {
		return err
	}

	if err := setStringList(d, "ssl_certificates", result.SSLCerts); err != nil {
		return err
	}

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	if err := setStringList(d, "virtual_hosts", result.VirtualHosts); err != nil {
		return err
	}

	return nil
}

func resourceListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.ListenerClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	input := lbaas.UpdateListenerInput{
		Name: d.Get("name").(string),
	}

	if balancerProtocol, ok := d.GetOk("balancer_protocol"); ok {
		input.BalancerProtocol = lbaas.Protocol(balancerProtocol.(string))
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		input.Disabled = getDisabledStateKeyword(enabled.(bool))
	}

	if port, ok := d.GetOk("port"); ok {
		input.Port = port.(int)
	}

	if serverPool, ok := d.GetOk("server_pool"); ok {
		// Only the URI Path is need on Update
		input.OriginServerPool = getURIRequestPath(serverPool.(string))
	}

	if serverProtocol, ok := d.GetOk("server_protocol"); ok {
		input.OriginServerProtocol = lbaas.Protocol(serverProtocol.(string))
	}

	pathPrefixes := getStringList(d, "path_prefixes")
	if len(pathPrefixes) != 0 {
		input.PathPrefixes = pathPrefixes
	}

	policies := getStringList(d, "policies")
	if len(policies) != 0 {
		input.Policies = policies
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	virtualHosts := getStringList(d, "virtual_hosts")
	if len(virtualHosts) != 0 {
		input.VirtualHosts = virtualHosts
	}

	result, err := lbaasClient.UpdateListener(lb, name, &input)
	if err != nil {
		return fmt.Errorf("Error updating Listener: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, result.Name))

	// TODO instead of re-read, process info from UpdateListener()
	return resourceListenerRead(d, meta)
}

func resourceListenerDelete(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.ListenerClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := lbaasClient.DeleteListener(lb, name); err != nil {
		return fmt.Errorf("Error deleting Listener")
	}
	return nil
}

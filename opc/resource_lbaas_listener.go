package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateLoadBalancerID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateLoadBalancerResourceName,
			},
			"balancer_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, true),
			},
			"server_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, true),
			},
			"port": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"server_pool": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOriginServerPoolURI,
			},
			"path_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"virtual_hosts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Read only attributes
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

	listenerClient := meta.(*Client).lbaasClient.ListenerClient()

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

	pathPrefixes := getStringSet(d, "path_prefixes")
	if len(pathPrefixes) != 0 {
		input.PathPrefixes = pathPrefixes
	}

	policies := getStringSet(d, "policies")
	if len(policies) != 0 {
		input.Policies = policies
	}

	sslCerts := getStringSet(d, "certificates")
	if len(sslCerts) != 0 {
		input.SSLCerts = sslCerts
	}

	tags := getStringSet(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	virtualHosts := getStringSet(d, "virtual_hosts")
	if len(virtualHosts) != 0 {
		input.VirtualHosts = virtualHosts
	}

	info, err := listenerClient.CreateListener(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer Listener: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourceListenerRead(d, meta)
}

func resourceListenerRead(d *schema.ResourceData, meta interface{}) error {
	listenerClient := meta.(*Client).lbaasClient.ListenerClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	result, err := listenerClient.GetListener(lb, name)
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

	if err := setStringList(d, "path_prefixes", result.PathPrefixes); err != nil {
		return err
	}

	if err := setStringList(d, "policies", result.Policies); err != nil {
		return err
	}

	if err := setStringList(d, "certificates", result.SSLCerts); err != nil {
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
	listenerClient := meta.(*Client).lbaasClient.ListenerClient()
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

	if serverProtocol, ok := d.GetOk("server_protocol"); ok {
		input.OriginServerProtocol = lbaas.Protocol(serverProtocol.(string))
	}

	if d.HasChange("server_pool") {
		if _, ok := d.GetOk("server_pool"); ok {
			// Only the URI Path is need on Update
			serverPool := updateOrRemoveStringAttribute(d, "server_pool")
			*serverPool = getURIRequestPath(*serverPool)
			input.OriginServerPool = serverPool
		} else {
			// server pool removed
			input.OriginServerPool = nil
		}
	}

	input.PathPrefixes = updateOrRemoveStringListAttribute(d, "path_prefixes")
	input.Policies = updateOrRemoveStringListAttribute(d, "policies")
	input.SSLCerts = updateOrRemoveStringListAttribute(d, "certificates")
	input.Tags = updateOrRemoveStringListAttribute(d, "tags")
	input.VirtualHosts = updateOrRemoveStringListAttribute(d, "virtual_hosts")

	result, err := listenerClient.UpdateListener(lb, name, &input)
	if err != nil {
		return fmt.Errorf("Error updating Listener: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, result.Name))

	return resourceListenerRead(d, meta)
}

func resourceListenerDelete(d *schema.ResourceData, meta interface{}) error {
	listenerClient := meta.(*Client).lbaasClient.ListenerClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := listenerClient.DeleteListener(lb, name); err != nil {
		return fmt.Errorf("Error deleting Listener")
	}
	return nil
}

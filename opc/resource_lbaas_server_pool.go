package opc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLBaaSOriginServerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceOriginServerPoolCreate,
		Read:   resourceOriginServerPoolRead,
		Update: resourceOriginServerPoolUpdate,
		Delete: resourceOriginServerPoolDelete,
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"servers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vnic_set": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateComputeResourceFQDN,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Read only attributes
			"consumers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_details": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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

func resourceOriginServerPoolCreate(d *schema.ResourceData, meta interface{}) error {

	serverPoolClient := meta.(*Client).lbaasClient.OriginServerPoolClient()

	var lb lbaas.LoadBalancerContext
	if loadBalancer, ok := d.GetOk("load_balancer"); ok {
		lb = getLoadBalancerContextFromID(loadBalancer.(string))
	}

	input := lbaas.CreateOriginServerPoolInput{
		Name: d.Get("name").(string),
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		if enabled.(bool) {
			input.Status = lbaas.LBaaSStatusEnabled
		} else {
			input.Status = lbaas.LBaaSStatusDisabled
		}
	}

	if vnicSet, ok := d.GetOk("vnic_set_name"); ok {
		input.VnicSetName = vnicSet.(string)
	}

	originServers := getStringList(d, "servers")
	if len(originServers) != 0 {
		servers, err := expandOriginServerConfig(originServers)
		if err != nil {
			return err
		}
		input.OriginServers = servers
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	info, err := serverPoolClient.CreateOriginServerPool(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer Server Pool: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourceOriginServerPoolRead(d, meta)
}

func resourceOriginServerPoolRead(d *schema.ResourceData, meta interface{}) error {
	serverPoolClient := meta.(*Client).lbaasClient.OriginServerPoolClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	result, err := serverPoolClient.GetOriginServerPool(lb, name)
	if err != nil {
		// OriginServerPool does not exist
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

	d.Set("consumers", result.Consumers)
	d.Set("enabled", result.Status == lbaas.LBaaSStatusEnabled)
	d.Set("name", result.Name)
	d.Set("operation_details", result.OperationDetails)
	d.Set("reason_for_disabling", result.ReasonForDisabling)
	d.Set("state", result.State)
	d.Set("uri", result.URI)
	d.Set("vnic_set_name", result.VnicSetName)
	d.Set("load_balancer", fmt.Sprintf("%s/%s", lb.Region, lb.Name))

	if err := setStringList(d, "servers", flattenOriginServerConfig(result.OriginServers)); err != nil {
		return err
	}

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	return nil
}

func resourceOriginServerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	serverPoolClient := meta.(*Client).lbaasClient.OriginServerPoolClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	input := lbaas.UpdateOriginServerPoolInput{
		Name: name,
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		if enabled.(bool) {
			input.Status = lbaas.LBaaSStatusEnabled
		} else {
			input.Status = lbaas.LBaaSStatusDisabled
		}
	}

	if d.HasChange("servers") {
		originServers := getStringList(d, "servers")
		servers := []lbaas.CreateOriginServerInput{}
		if len(originServers) > 0 {
			expanded, err := expandOriginServerConfig(originServers)
			if err != nil {
				return err
			}
			servers = expanded
		}
		input.OriginServers = &servers
	}

	input.Tags = updateOrRemoveStringListAttribute(d, "tags")

	result, err := serverPoolClient.UpdateOriginServerPool(lb, name, &input)
	if err != nil {
		return fmt.Errorf("Error updating OriginServerPool: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, result.Name))

	return resourceOriginServerPoolRead(d, meta)
}

func resourceOriginServerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	lbaasClient := meta.(*Client).lbaasClient.OriginServerPoolClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := lbaasClient.DeleteOriginServerPool(lb, name); err != nil {
		return fmt.Errorf("Error deleting OriginServerPool")
	}
	return nil
}

// convert the list of "server:port" strings to a list of CreateOriginServerInput
func expandOriginServerConfig(servers []string) ([]lbaas.CreateOriginServerInput, error) {
	config := []lbaas.CreateOriginServerInput{}
	for _, element := range servers {
		s := strings.Split(element, ":")
		port, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("Server Pool servers must be in the format \"host:port\"")
		}
		server := lbaas.CreateOriginServerInput{
			Hostname: s[0],
			Port:     port,
			Status:   lbaas.LBaaSStatusEnabled,
		}
		config = append(config, server)
	}
	return config, nil
}

// convert the OriginServerInfo reponse to a listing of "server:port" strings
func flattenOriginServerConfig(info []lbaas.OriginServerInfo) []string {
	servers := []string{}
	for _, config := range info {
		servers = append(servers, fmt.Sprintf("%s:%d", config.Hostname, config.Port))
	}
	return servers
}

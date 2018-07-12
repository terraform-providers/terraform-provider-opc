package opc

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_return_codes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							// Computed: true,
							Default:      5,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							// Computed: true,
							Default:      30,
							ValidateFunc: validation.IntBetween(5, 300),
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							// Computed: true,
							Default: "",
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							// Computed: true,
							Default:      20,
							ValidateFunc: validation.IntBetween(2, 60),
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "http",
							ValidateFunc: validation.StringInSlice([]string{
								"http",
							}, true),
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							// Computed: true,
							Default:      7,
							ValidateFunc: validation.IntBetween(2, 10),
						},
					},
				},
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
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	serverPoolClient := lbaasClient.OriginServerPoolClient()
	lb := getLoadBalancerContextFromID(d.Get("load_balancer").(string))

	status := lbaas.LBaaSStatusEnabled
	if !d.Get("enabled").(bool) {
		status = lbaas.LBaaSStatusDisabled
	}

	input := lbaas.CreateOriginServerPoolInput{
		Name:   d.Get("name").(string),
		Status: status,
	}

	if vnicSet, ok := d.GetOk("vnic_set_name"); ok {
		input.VnicSetName = vnicSet.(string)
	}

	originServers := getStringSet(d, "servers")
	if len(originServers) != 0 {
		servers, err := expandOriginServerConfig(originServers)
		if err != nil {
			return err
		}
		input.OriginServers = servers
	}

	tags := getStringSet(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if _, ok := d.GetOk("health_check"); ok {
		healthCheck := expandHealthCheckConfig(d)
		input.HealthCheck = &healthCheck
	}

	info, err := serverPoolClient.CreateOriginServerPool(lb, &input)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer Server Pool: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lb.Region, lb.Name, info.Name))
	return resourceOriginServerPoolRead(d, meta)
}

func resourceOriginServerPoolRead(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	serverPoolClient := lbaasClient.OriginServerPoolClient()
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

	d.Set("tags", result.Tags)

	if result.HealthCheck.Enabled != "" {
		if err := flattenHealthCheckConfig(d, result.HealthCheck); err != nil {
			return err
		}
	}

	return nil
}

func resourceOriginServerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	serverPoolClient := lbaasClient.OriginServerPoolClient()
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
		originServers := getStringSet(d, "servers")
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

	if d.HasChange("health_check") {
		if _, ok := d.GetOk("health_check"); ok {
			healthCheck := expandHealthCheckConfig(d)
			input.HealthCheck = &healthCheck
		} else {
			input.HealthCheck = &lbaas.HealthCheckInfo{}
		}
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
	lbaasClient, err := meta.(*Client).getLBaaSClient()
	if err != nil {
		return err
	}
	serverPoolClient := lbaasClient.OriginServerPoolClient()
	name := getLastNameInPath(d.Id())
	lb := getLoadBalancerContextFromID(d.Id())

	if _, err := serverPoolClient.DeleteOriginServerPool(lb, name); err != nil {
		return fmt.Errorf("Error deleting Server Pool: %v", err)
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

func expandHealthCheckConfig(d *schema.ResourceData) lbaas.HealthCheckInfo {

	if v, ok := d.GetOk("health_check"); ok {
		vL := v.([]interface{})
		config := vL[0].(map[string]interface{})

		enabled := "FALSE"
		if config["enabled"].(bool) {
			enabled = "TRUE"
		}

		info := lbaas.HealthCheckInfo{
			Type:               config["type"].(string),
			Path:               config["path"].(string),
			Enabled:            enabled,
			Interval:           config["interval"].(int),
			Timeout:            config["timeout"].(int),
			HealthyThreshold:   config["healthy_threshold"].(int),
			UnhealthyThreshold: config["unhealthy_threshold"].(int),
		}

		returnCodes := []string{}
		if codes, ok := config["accepted_return_codes"]; ok && v != nil {
			for _, code := range codes.([]interface{}) {
				returnCodes = append(returnCodes, code.(string))
			}
			if len(returnCodes) > 0 {
				sort.Strings(returnCodes)
				info.AcceptedReturnCodes = returnCodes
			}
		}

		return info
	}
	return lbaas.HealthCheckInfo{}
}

func flattenHealthCheckConfig(d *schema.ResourceData, info lbaas.HealthCheckInfo) error {

	config := make([]map[string]interface{}, 0)

	attrs := make(map[string]interface{})
	attrs["type"] = info.Type
	attrs["path"] = info.Path
	attrs["accepted_return_codes"] = info.AcceptedReturnCodes
	attrs["interval"] = info.Interval
	attrs["timeout"] = info.Timeout
	attrs["healthy_threshold"] = info.HealthyThreshold
	attrs["unhealthy_threshold"] = info.UnhealthyThreshold

	if info.Enabled == "TRUE" {
		attrs["enabled"] = true
	} else {
		attrs["enabled"] = false
	}

	config = append(config, attrs)

	return d.Set("health_check", config)
}

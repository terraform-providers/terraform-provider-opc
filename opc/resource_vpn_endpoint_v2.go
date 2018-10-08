package opc

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCVPNEndpointV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCVPNEndpointV2Create,
		Read:   resourceOPCVPNEndpointV2Read,
		Update: resourceOPCVPNEndpointV2Update,
		Delete: resourceOPCVPNEndpointV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"customer_vpn_gateway": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ike_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"require_perfect_forward_secrecy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"phase_one_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"hash": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"dh_group": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"phase_two_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"hash": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"pre_shared_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reachable_routes": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vnic_sets": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tunnel_status": {
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

func resourceOPCVPNEndpointV2Create(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating VPNEndpointV2")

	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.VPNEndpointV2s()
	input := compute.CreateVPNEndpointV2Input{
		Name:               d.Get("name").(string),
		Enabled:            d.Get("enabled").(bool),
		CustomerVPNGateway: d.Get("customer_vpn_gateway").(string),
		IPNetwork:          d.Get("ip_network").(string),
		PSK:                d.Get("pre_shared_key").(string),
		VNICSets:           getStringList(d, "vnic_sets"),
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	if ikeIdentifier, ok := d.GetOk("ike_identifier"); ok {
		input.IKEIdentifier = ikeIdentifier.(string)
	}

	if pfsFlag, ok := d.GetOk("require_perfect_forward_secrecy"); ok {
		input.PFSFlag = pfsFlag.(bool)
	}

	if _, ok := d.GetOk("phase_one_settings"); ok {
		input.Phase1Settings = expandVPNEndpoingV2PhaseOneSettings(d)
	}

	if _, ok := d.GetOk("phase_two_settings"); ok {
		input.Phase2Settings = expandVPNEndpoingV2PhaseTwoSettings(d)
	}

	info, err := resClient.CreateVPNEndpointV2(&input)
	if err != nil {
		return fmt.Errorf("Error creating VPNEndpointV2: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCVPNEndpointV2Read(d, meta)
}

func resourceOPCVPNEndpointV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.VPNEndpointV2s()

	log.Printf("[DEBUG] Reading state of VPNEndpointV2 %s", d.Id())
	getInput := compute.GetVPNEndpointV2Input{
		Name: d.Id(),
	}

	result, err := resClient.GetVPNEndpointV2(&getInput)
	if err != nil {
		// VPNEndpointV2 does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading VPNEndpointV2 %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of VPNEndpointV2 %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("customer_vpn_endpoint", result.CustomerVPNGateway)
	d.Set("enabled", result.Enabled)
	d.Set("ike_identifier", result.IKEIdentifier)
	d.Set("ip_network", result.IPNetwork)
	d.Set("require_perfect_forward_secrecy", result.PFSFlag)
	d.Set("pre_shared_key", result.PSK)
	d.Set("uri", result.URI)
	d.Set("tunnel_status", string(result.TunnelStatus))

	if err := setStringList(d, "reachable_routes", result.ReachableRoutes); err != nil {
		return err
	}
	if err := setStringList(d, "vnic_sets", result.VNICSets); err != nil {
		return err
	}

	if p1Settings := result.Phase1Settings; &p1Settings != nil {
		if err := d.Set("phase_one_settings", flattenVPNEndpointV2PhaseOneSettings(p1Settings)); err != nil {
			return err
		}
	}

	if p2Settings := result.Phase2Settings; &p2Settings != nil {
		if err := d.Set("phase_two_settings", flattenVPNEndpointV2PhaseTwoSettings(p2Settings)); err != nil {
			return err
		}
	}

	return nil
}

func resourceOPCVPNEndpointV2Update(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating VPNEndpointV2")

	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.VPNEndpointV2s()
	input := compute.UpdateVPNEndpointV2Input{
		Name:               d.Get("name").(string),
		Enabled:            d.Get("enabled").(bool),
		CustomerVPNGateway: d.Get("customer_vpn_gateway").(string),
		IPNetwork:          d.Get("ip_network").(string),
		PSK:                d.Get("pre_shared_key").(string),
		VNICSets:           getStringList(d, "vnic_sets"),
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	if ikeIdentifier, ok := d.GetOk("ike_identifier"); ok {
		input.IKEIdentifier = ikeIdentifier.(string)
	}

	if pfsFlag, ok := d.GetOk("require_perfect_forward_secrecy"); ok {
		input.PFSFlag = pfsFlag.(bool)
	}

	if _, ok := d.GetOk("phase_one_settings"); ok {
		input.Phase1Settings = expandVPNEndpoingV2PhaseOneSettings(d)
	}

	if _, ok := d.GetOk("phase_two_settings"); ok {
		input.Phase2Settings = expandVPNEndpoingV2PhaseTwoSettings(d)
	}

	info, err := resClient.UpdateVPNEndpointV2(&input)
	if err != nil {
		return fmt.Errorf("Error creating VPNEndpointV2: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCVPNEndpointV2Read(d, meta)
}

func resourceOPCVPNEndpointV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.VPNEndpointV2s()
	name := d.Id()

	log.Printf("[DEBUG] Deleting VPNEndpointV2: %v", name)

	input := compute.DeleteVPNEndpointV2Input{
		Name: name,
	}
	if err := resClient.DeleteVPNEndpointV2(&input); err != nil {
		return fmt.Errorf("Error deleting VPNEndpointV2")
	}
	return nil
}

func flattenVPNEndpointV2PhaseOneSettings(input compute.Phase1Settings) []interface{} {

	settings := make(map[string]interface{}, 0)

	settings["encryption"] = input.Encryption
	settings["hash"] = input.Hash
	settings["dh_group"] = input.DHGroup

	return []interface{}{settings}
}

func flattenVPNEndpointV2PhaseTwoSettings(input compute.Phase2Settings) []interface{} {

	settings := make(map[string]interface{}, 0)

	settings["encryption"] = input.Encryption
	settings["hash"] = input.Hash

	return []interface{}{settings}
}

func expandVPNEndpoingV2PhaseOneSettings(d *schema.ResourceData) *compute.Phase1Settings {
	phase1Settings := d.Get("phase_one_settings").([]interface{})

	attrs := phase1Settings[0].(map[string]interface{})

	result := &compute.Phase1Settings{
		Encryption: attrs["encryption"].(string),
		Hash:       attrs["hash"].(string),
		DHGroup:    attrs["dh_group"].(string),
	}

	return result
}

func expandVPNEndpoingV2PhaseTwoSettings(d *schema.ResourceData) *compute.Phase2Settings {
	phase1Settings := d.Get("phase_one_settings").([]interface{})

	attrs := phase1Settings[0].(map[string]interface{})

	result := &compute.Phase2Settings{
		Encryption: attrs["encryption"].(string),
		Hash:       attrs["hash"].(string),
	}

	return result
}

package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCVNICSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCVNICSetCreate,
		Read:   resourceOPCVNICSetRead,
		Update: resourceOPCVNICSetUpdate,
		Delete: resourceOPCVNICSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"applied_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"virtual_nics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOPCVNICSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).VirtNICSets()

	name := d.Get("name").(string)
	desc, descOk := d.GetOk("description")

	input := &compute.CreateVirtualNICSetInput{
		Name: name,
	}

	if descOk {
		input.Description = desc.(string)
	}

	acls := getStringList(d, "applied_acls")
	if len(acls) != 0 {
		input.AppliedACLs = acls
	}

	vnics := getStringList(d, "virtual_nics")
	if len(vnics) != 0 {
		input.VirtualNICs = vnics
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	vnicSet, err := client.CreateVirtualNICSet(input)
	if err != nil {
		return fmt.Errorf("Error creating Virtual NIC Set: %s", err)
	}

	d.SetId(vnicSet.Name)

	return resourceOPCVNICSetRead(d, meta)
}

func resourceOPCVNICSetRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*compute.ComputeClient).VirtNICSets()

	name := d.Id()
	input := &compute.GetVirtualNICSetInput{
		Name: name,
	}

	result, err := computeClient.GetVirtualNICSet(input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Virtual NIC Set '%s': %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
	if err := setStringList(d, "applied_acls", result.AppliedACLs); err != nil {
		return err
	}
	if err := setStringList(d, "virtual_nics", result.VirtualNICs); err != nil {
		return err
	}
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCVNICSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).VirtNICSets()

	name := d.Id()
	desc, descOk := d.GetOk("description")

	input := &compute.UpdateVirtualNICSetInput{
		Name: name,
	}

	if descOk {
		input.Description = desc.(string)
	}

	acls := getStringList(d, "applied_acls")
	if len(acls) != 0 {
		input.AppliedACLs = acls
	}

	vnics := getStringList(d, "virtual_nics")
	if len(vnics) != 0 {
		input.VirtualNICs = vnics
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	info, err := client.UpdateVirtualNICSet(input)
	if err != nil {
		return fmt.Errorf("Error updating Virtual NIC Set: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCVNICSetRead(d, meta)
}

func resourceOPCVNICSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).VirtNICSets()

	name := d.Id()
	input := &compute.DeleteVirtualNICSetInput{
		Name: name,
	}

	if err := client.DeleteVirtualNICSet(input); err != nil {
		return fmt.Errorf("Error deleting Virtual NIC Set '%s': %s", name, err)
	}
	return nil
}

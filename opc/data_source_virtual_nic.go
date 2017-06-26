package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVNIC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVNICRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsComputedSchema(),

			"transit_flag": {
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

func dataSourceVNICRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*compute.ComputeClient).VirtNICs()

	name := d.Get("name").(string)

	input := &compute.GetVirtualNICInput{
		Name: name,
	}

	vnic, err := computeClient.GetVirtualNIC(input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading vnic %s: %s", name, err)
	}

	if vnic == nil {
		d.SetId("")
		return nil
	}

	d.SetId(name)
	d.Set("description", vnic.Description)
	d.Set("mac_address", vnic.MACAddress)
	d.Set("transit_flag", vnic.TransitFlag)
	d.Set("uri", vnic.Uri)
	if err := setStringList(d, "tags", vnic.Tags); err != nil {
		return err
	}
	return nil
}

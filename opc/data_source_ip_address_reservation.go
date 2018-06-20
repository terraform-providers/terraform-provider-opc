package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIPAddressReservation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPAddressReservationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ip_address_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsComputedSchema(),
			"ip_address": {
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

func dataSourceIPAddressReservationRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*Client).computeClient.IPAddressReservations()
	name := d.Get("name").(string)

	input := compute.GetIPAddressReservationInput{
		Name: name,
	}

	result, err := computeClient.GetIPAddressReservation(&input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading ip address reservation %s: %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.SetId(name)
	d.Set("description", result.Description)
	d.Set("ip_address_pool", result.IPAddressPool)
	d.Set("ip_address", result.IPAddress)
	d.Set("uri", result.URI)

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	return nil
}

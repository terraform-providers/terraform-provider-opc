package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*OPCClient).computeClient.SSHKeys()
	name := d.Get("name").(string)

	input := compute.GetSSHKeyInput{
		Name: name,
	}

	result, err := computeClient.GetSSHKey(&input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading vnic %s: %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.SetId(name)
	d.Set("key", result.Key)
	d.Set("enabled", result.Enabled)

	return nil
}

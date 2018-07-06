package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
)

func dataSourceMachineImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMachineImageRead,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Required: true,
			},

			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"error_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hypervisor": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"image_format": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"no_upload": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sizes": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"state": {
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

func dataSourceMachineImageRead(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.MachineImages()

	// Get required attributes
	account := d.Get("account").(string)
	name := d.Get("name").(string)

	input := compute.GetMachineImageInput{
		Account: account,
		Name:    name,
	}

	result, err := resClient.GetMachineImage(&input)
	if err != nil {
		return err
	}

	// Not found, don't error
	if result == nil {
		d.SetId("")
		return nil
	}

	// Flatten JSON attributes
	attributes, err := structure.FlattenJsonToString(result.Attributes)
	if err != nil {
		return err
	}

	// Populate schema attributes
	d.SetId(fmt.Sprintf("%s", result.Name))
	d.Set("account", result.Account)
	d.Set("audited", result.Audited)
	d.Set("description", result.Description)
	d.Set("error_reason", result.ErrorReason)
	d.Set("hypervisor", result.Hypervisor)
	d.Set("image_format", result.ImageFormat)
	d.Set("file", result.File)
	d.Set("name", result.Name)
	d.Set("no_upload", result.NoUpload)
	d.Set("platform", result.Platform)
	d.Set("sizes", result.Sizes)
	d.Set("state", result.State)
	d.Set("uri", result.URI)

	if err := d.Set("attributes", attributes); err != nil {
		return err
	}

	return nil
}

package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCMachineImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCMachineImageCreate,
		Read:   resourceOPCMachineImageRead,
		Delete: resourceOPCMachineImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// ValidateFunc: validateAccount, TODO
			},

			// "attributes": {
			// 	Type:     schema.TypeMap,
			// 	Optional: true,
			// 	ForceNew: true,
			// },

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"no_upload": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// "sizes": {
			// 	Type:     schema.TypeMap,
			// 	Computed: true,
			// },

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

func resourceOPCMachineImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).computeClient.MachineImages()

	// Get required attributes
	name := d.Get("name").(string)
	file := d.Get("file").(string)

	input := &compute.CreateMachineImageInput{
		Name: name,
		File: file,
	}

	if account, ok := d.GetOk("account"); ok && account != nil {
		input.Account = account.(string)
	} else {
		input.Account = fmt.Sprintf("/Compute-%s/cloud_storage", d.Get("identity_domain").(string)) // XXX TODO BROKEN
	}

	// Get Optional attributes
	if desc, ok := d.GetOk("description"); ok && desc != nil {
		input.Description = desc.(string)
	}

	info, err := client.CreateMachineImage(input)
	if err != nil {
		return fmt.Errorf("Error creating Machine Image '%s': %v", name, err)
	}

	d.SetId(info.Name)

	return resourceOPCMachineImageRead(d, meta)
}

func resourceOPCMachineImageRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*OPCClient).computeClient.MachineImages()

	name := d.Id()
	input := &compute.GetMachineImageInput{
		Name: name,
	}

	result, err := computeClient.GetMachineImage(input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Machine Image '%s': %v", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("account", result.Account)
	// d.Set("attributes", result.Attributes)
	d.Set("audited", result.Audited)
	d.Set("description", result.Description)
	d.Set("error_reason", result.ErrorReason)
	d.Set("hypervisor", result.Hypervisor)
	d.Set("image_format", result.ImageFormat)
	d.Set("file", result.File)
	d.Set("name", result.Name)
	d.Set("no_upload", result.NoUpload)
	d.Set("platform", result.Platform)
	// d.Set("sizes", result.Sizes)
	d.Set("state", result.State)
	d.Set("uri", result.URI)
	return nil
}

func resourceOPCMachineImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).computeClient.MachineImages()

	name := d.Id()
	input := &compute.DeleteMachineImageInput{
		Name: name,
	}

	if err := client.DeleteMachineImage(input); err != nil {
		return fmt.Errorf("Error deleting Machine Image '%s': %v", name, err)
	}
	return nil
}

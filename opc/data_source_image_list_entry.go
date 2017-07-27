package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
)

func dataSourceImageListEntry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageListEntryRead,

		Schema: map[string]*schema.Schema{
			"entry": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"image_list": {
				Type:     schema.TypeString,
				Required: true,
			},

			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Computed Attributes
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"machine_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceImageListEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).computeClient.ImageListEntries()

	// Get required attributes
	image_list := d.Get("image_list").(string)
	version := d.Get("version").(int)

	// Get the image list (we'll parse out the entry later on)
	input := compute.GetImageListEntryInput{
		Name:    image_list,
		Version: version,
	}

	result, err := client.GetImageListEntry(&input)
	if err != nil {
		return err
	}

	// Not found, don't error
	if result == nil {
		d.SetId("")
		return nil
	}

	// If entry was specified, only return the single entry, otherwise the entire list
	var images []string
	var entry int
	if v, ok := d.GetOk("entry"); ok {
		entry = v.(int)
		// Protect against index panic
		if len(result.MachineImages) <= entry-1 {
			return fmt.Errorf("Invalid entry specified. Image list only has %d entries, got: %d",
				len(result.MachineImages), entry)
		}
		images = append(images, result.MachineImages[entry-1])
	} else {
		images = result.MachineImages
	}

	// Flatten JSON attributes
	attrs, err := structure.FlattenJsonToString(result.Attributes)
	if err != nil {
		return err
	}

	// Populate schema attributes
	d.SetId(fmt.Sprintf("%s|%s:%d", image_list, version, entry))
	d.Set("uri", result.Uri)
	if err := d.Set("attributes", attrs); err != nil {
		return err
	}
	if err := d.Set("machine_images", images); err != nil {
		return err
	}

	return nil
}

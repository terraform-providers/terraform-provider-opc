package opc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCImageListEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCImageListEntryCreate,
		Read:   resourceOPCImageListEntryRead,
		Delete: resourceOPCImageListEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_images": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"attributes": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCImageListEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageListEntries()

	name := d.Get("name").(string)
	machineImages := expandOPCImageListEntryMachineImages(d)
	version := d.Get("version").(int)

	createInput := &compute.CreateImageListEntryInput{
		Name:          name,
		MachineImages: machineImages,
		Version:       version,
	}

	if v, ok := d.GetOk("attributes"); ok {
		attributesString := v.(string)
		attributes, err := structure.ExpandJsonFromString(attributesString)
		if err != nil {
			return err
		}

		createInput.Attributes = attributes
	}

	_, err := client.CreateImageListEntry(createInput)
	if err != nil {
		return err
	}

	id := generateOPCImageListEntryID(name, version)
	d.SetId(id)
	return resourceOPCImageListEntryRead(d, meta)
}

func resourceOPCImageListEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageListEntries()

	name, version, err := parseOPCImageListEntryID(d.Id())
	if err != nil {
		return err
	}

	input := compute.GetImageListEntryInput{
		Name:    *name,
		Version: *version,
	}

	result, err := client.GetImageListEntry(&input)
	if err != nil {
		return err
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	attrs, err := structure.FlattenJsonToString(result.Attributes)
	if err != nil {
		return err
	}

	d.Set("name", name)
	d.Set("machine_images", result.MachineImages)
	d.Set("version", result.Version)
	d.Set("attributes", attrs)
	d.Set("uri", result.Uri)

	return nil
}

func resourceOPCImageListEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageListEntries()

	name, version, err := parseOPCImageListEntryID(d.Id())
	if err != nil {
		return err
	}

	deleteInput := &compute.DeleteImageListEntryInput{
		Name:    *name,
		Version: *version,
	}
	err = client.DeleteImageListEntry(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func parseOPCImageListEntryID(id string) (*string, *int, error) {
	s := strings.Split(id, "|")
	name, versionString := s[0], s[1]
	version, err := strconv.Atoi(versionString)
	if err != nil {
		return nil, nil, err
	}

	return &name, &version, nil
}

func expandOPCImageListEntryMachineImages(d *schema.ResourceData) []string {
	machineImages := []string{}
	for _, i := range d.Get("machine_images").([]interface{}) {
		machineImages = append(machineImages, i.(string))
	}
	return machineImages
}

func generateOPCImageListEntryID(name string, version int) string {
	return fmt.Sprintf("%s|%d", name, version)
}

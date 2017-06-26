package opc

import (
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCImageList() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCImageListCreate,
		Read:   resourceOPCImageListRead,
		Update: resourceOPCImageListUpdate,
		Delete: resourceOPCImageListDelete,
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
				Required: true,
			},
			"default": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceOPCImageListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageList()

	name := d.Get("name").(string)

	createInput := &compute.CreateImageListInput{
		Name:        name,
		Description: d.Get("description").(string),
		Default:     d.Get("default").(int),
	}

	createResult, err := client.CreateImageList(createInput)
	if err != nil {
		return err
	}

	d.SetId(createResult.Name)

	return resourceOPCImageListRead(d, meta)
}

func resourceOPCImageListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageList()

	name := d.Id()

	updateInput := &compute.UpdateImageListInput{
		Name:        name,
		Description: d.Get("description").(string),
		Default:     d.Get("default").(int),
	}

	_, err := client.UpdateImageList(updateInput)
	if err != nil {
		return err
	}

	return resourceOPCImageListRead(d, meta)
}

func resourceOPCImageListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageList()

	input := &compute.GetImageListInput{
		Name: d.Id(),
	}

	result, err := client.GetImageList(input)
	if err != nil {
		return err
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("default", result.Default)

	return nil
}

func resourceOPCImageListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*compute.ComputeClient).ImageList()

	deleteInput := &compute.DeleteImageListInput{
		Name: d.Id(),
	}
	err := client.DeleteImageList(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

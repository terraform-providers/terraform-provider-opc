package opc

import (
	"fmt"
	"log"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
  "github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCJavaServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCJavaServiceInstanceCreate,
		Read:   resourceOPCJavaServiceInstanceRead,
		Delete: resourceOPCJavaServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
      "edition": {
        Type: schema.TypeString,
        Required: true,
        ForceNew: true,
        ValidateFunc: validation.StringInSlice([]string{
 				  string(java.ServiceInstanceEditionSE),
 					string(java.ServiceInstanceEditionEE),
          string(java.ServiceInstanceEditionSuite),
 				}, true),
      },
      "level": {
        Type: schema.TypeString,
        Optional: true,
        ForceNew: true,
        Default: "PAAS",
        ValidateFunc: validation.StringInSlice([]string{
 				  string(java.ServiceInstanceLevelPAAS),
 					string(java.ServiceInstanceLevelBasic),
 				}, true),
      },
		},
	}
}

func resourceOPCJavaServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating JavaServiceInstance")

	client := meta.(*OPCClient).computeClient.JavaServiceInstances()
	input := compute.CreateJavaServiceInstanceInput{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	info, err := client.CreateJavaServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating JavaServiceInstance: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCJavaServiceInstanceRead(d, meta)
}

func resourceOPCJavaServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).computeClient.JavaServiceInstances()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := compute.GetJavaServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetJavaServiceInstance(&getInput)
	if err != nil {
		// Java Service Instance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading JavaServiceInstance %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of JavaServiceInstance %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("enabled", result.Enabled)
	d.Set("description", result.Description)
	d.Set("uri", result.URI)
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCJavaServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).computeClient.JavaServiceInstances()
	name := d.Id()

	log.Printf("[DEBUG] Deleting JavaServiceInstance: %v", name)

	input := compute.DeleteJavaServiceInstanceInput{
		Name: name,
	}
	if err := client.DeleteJavaServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting JavaServiceInstance")
	}
	return nil
}

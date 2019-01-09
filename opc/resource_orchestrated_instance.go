package opc

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCOrchestratedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCOrchestratedInstanceCreate,
		Read:   resourceOPCOrchestratedInstanceRead,
		Delete: resourceOPCOrchestratedInstanceDelete,
		Update: resourceOPCOrchestratedInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
			"desired_state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"inactive",
					"suspend",
				}, true),
			},
			"tags": tagsOptionalSchema(),

			"instance": orchestrationInstanceSchema(),

			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceOPCOrchestratedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating Orchestration")

	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Orchestrations()
	input := compute.CreateOrchestrationInput{
		Name:         d.Get("name").(string),
		DesiredState: compute.OrchestrationDesiredState(d.Get("desired_state").(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	instances, err := expandOrchestrationInstances(d)
	if err != nil {
		return err
	}
	input.Objects = instances

	info, err := resClient.CreateOrchestration(&input)
	if err != nil {
		return fmt.Errorf("Error creating Orchestration: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCOrchestratedInstanceRead(d, meta)
}

func resourceOPCOrchestratedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Orchestrations()

	log.Printf("[DEBUG] Reading state of orchestrated instance %s", d.Id())
	getInput := compute.GetOrchestrationInput{
		Name: d.Id(),
	}

	result, err := resClient.GetOrchestration(&getInput)
	if err != nil {
		// Orchestration does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Orchestration %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of Orchestration %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("version", result.Version)
	d.Set("description", result.Description)
	d.Set("desired_state", result.DesiredState)

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	if result.DesiredState == "active" {
		instances, err := flattenOrchestratedInstances(d, meta, result.Objects)
		if err != nil {
			return err
		}
		if err := d.Set("instance", instances); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Instances error: %#v", err)
		}
	}

	return nil
}

func resourceOPCOrchestratedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Updating Orchestration")

	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Orchestrations()

	// Obtain orchestration so we can grab the instance information
	getInput := compute.GetOrchestrationInput{
		Name: d.Id(),
	}

	result, err := resClient.GetOrchestration(&getInput)
	if err != nil {
		// Orchestration does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Orchestration %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	input := compute.UpdateOrchestrationInput{
		Name:         d.Get("name").(string),
		DesiredState: compute.OrchestrationDesiredState(d.Get("desired_state").(string)),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Version:      d.Get("version").(int),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	input.Objects = result.Objects

	info, err := resClient.UpdateOrchestration(&input)
	if err != nil {
		return fmt.Errorf("Error updating Orchestration: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCOrchestratedInstanceRead(d, meta)
}

func resourceOPCOrchestratedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Orchestrations()

	name := d.Id()

	input := &compute.DeleteOrchestrationInput{
		Name:    name,
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	log.Printf("[DEBUG] Deleting orchestration %s", name)

	if err := resClient.DeleteOrchestration(input); err != nil {
		return fmt.Errorf("Error deleting orchestration %s for instance %s: %s", name, d.Id(), err)
	}

	return nil
}

func expandOrchestrationInstances(d *schema.ResourceData) ([]compute.Object, error) {
	instances_info := d.Get("instance").([]interface{})
	instances := make([]compute.Object, 0, len(instances_info))
	for i := range instances_info {
		// The value for orchestration is the name of the orchestration
		orchestrationName := d.Get("name").(string)
		objectLabel := d.Get(fmt.Sprintf("instance.%d.name", i)).(string)
		persistent := d.Get(fmt.Sprintf("instance.%d.persistent", i)).(bool)

		instanceCreateInput, instanceErr := expandCreateInstanceInput(fmt.Sprintf("instance.%d", i), d)
		if instanceErr != nil {
			return nil, instanceErr
		}

		instance := compute.Object{
			Label:         objectLabel,
			Orchestration: orchestrationName,
			Type:          compute.OrchestrationTypeInstance,
			Template:      instanceCreateInput,
			Persistent:    persistent,
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

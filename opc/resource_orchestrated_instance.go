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
		// TODO Change this to resourceOPCOrchestratedInstanceUpdate
		Update: resourceOPCOrchestratedInstanceCreate,
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

	client := meta.(*OPCClient).computeClient.Orchestrations()
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

	info, err := client.CreateOrchestration(&input)
	if err != nil {
		return fmt.Errorf("Error creating Orchestration: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCOrchestratedInstanceRead(d, meta)
}

func resourceOPCOrchestratedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient := meta.(*OPCClient).computeClient.Orchestrations()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := compute.GetOrchestrationInput{
		Name: d.Id(),
	}

	result, err := computeClient.GetOrchestration(&getInput)
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
	return nil
}

func resourceOPCOrchestratedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).computeClient.Orchestrations()

	name := d.Id()

	input := &compute.DeleteOrchestrationInput{
		Name:    name,
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	log.Printf("[DEBUG] Deleting orchestration %s", name)

	if err := client.DeleteOrchestration(input); err != nil {
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

		instanceCreateInput, instanceErr := getCreateInstanceInput(fmt.Sprintf("instance.%d", i), d)
		if instanceErr != nil {
			return nil, instanceErr
		}

		instance := compute.Object{
			Label:         orchestrationName,
			Orchestration: orchestrationName,
			Type:          compute.OrchestrationTypeInstance,
			Template:      instanceCreateInput,
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

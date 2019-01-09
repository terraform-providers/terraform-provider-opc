package opc

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCSnapshotCreate,
		Read:   resourceOPCSnapshotRead,
		Delete: resourceOPCSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_image": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Snapshots()

	instance := d.Get("instance").(string)

	input := compute.CreateSnapshotInput{
		Instance: instance,
		Timeout:  d.Timeout(schema.TimeoutCreate),
	}

	if account, ok := d.GetOk("description"); ok {
		input.Account = account.(string)
	}

	if machineImage, ok := d.GetOk("machine_image"); ok {
		input.MachineImage = machineImage.(string)
	}

	info, err := resClient.CreateSnapshot(&input)
	if err != nil {
		return fmt.Errorf("Error creating snapshot %s: %s", instance, err)
	}

	d.SetId(info.Name)

	return resourceOPCSnapshotRead(d, meta)
}

func resourceOPCSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.Snapshots()

	name := d.Id()

	input := compute.GetSnapshotInput{
		Name: name,
	}
	result, err := resClient.GetSnapshot(&input)
	if err != nil {
		// Sec Rule does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading snapshot %s: %s", name, err)
	}

	d.Set("name", result.Name)
	d.Set("account", result.Account)
	d.Set("creation_time", result.CreationTime)
	d.Set("machine_image", result.MachineImage)
	d.Set("instance", result.Instance)
	d.Set("uri", result.URI)

	return nil
}

func resourceOPCSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	snapshotClient := computeClient.Snapshots()
	machineImageClient := computeClient.MachineImages()
	name := d.Id()

	getInput := compute.GetSnapshotInput{
		Name: name,
	}
	result, err := snapshotClient.GetSnapshot(&getInput)
	if err != nil {
		// Snapshot does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading snapshot %s: %s", name, err)
	}

	input := compute.DeleteSnapshotInput{
		Snapshot:     name,
		MachineImage: result.MachineImage,
		Timeout:      d.Timeout(schema.TimeoutDelete),
	}
	if err := snapshotClient.DeleteSnapshot(machineImageClient, &input); err != nil {
		return fmt.Errorf("Error deleting snapshot %s: %s", name, err)
	}

	return nil
}

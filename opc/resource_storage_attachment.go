package opc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCStorageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCStorageAttachmentCreate,
		Read:   resourceOPCStorageAttachmentRead,
		Delete: resourceOPCStorageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"index": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"storage_volume": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOPCStorageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating storage_attachment")

	volumeName := d.Get("storage_volume").(string)
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.StorageVolumes()
	getVolumeInput := compute.GetStorageVolumeInput{
		Name: volumeName,
	}
	storageVolume, err := resClient.GetStorageVolume(&getVolumeInput)
	if err != nil {
		if client.WasNotFoundError(err) {
			return fmt.Errorf("Unable to find storage volume: %s", volumeName)
		}
		return fmt.Errorf("Error reading storage volume %s: %s", volumeName, err)
	}
	if storageVolume == nil {
		// Volume doesn't exist
		return fmt.Errorf("Unable to find storage volume: %s", volumeName)
	}

	instanceName := d.Get("instance").(string)
	instanceClient := meta.(*Client).computeClient.Instances()
	getInstanceInput := &compute.GetInstanceIDInput{
		Name: instanceName,
	}
	instance, err := instanceClient.GetInstanceFromName(getInstanceInput)
	if err != nil {
		return fmt.Errorf("Unable to find Instance: %s", instanceName)
	}
	if instance == nil {
		return fmt.Errorf("Unable to find Instance: %s", instanceName)
	}

	volumeIndex := d.Get("index").(int)
	if !checkForEmptyIndex(instance.Storage, volumeIndex) {
		return fmt.Errorf("Storage index %d is already in use on instance %s", volumeIndex, instanceName)
	}

	storageAttachmentClient := meta.(*Client).computeClient.StorageAttachments()
	input := compute.CreateStorageAttachmentInput{
		StorageVolumeName: storageVolume.Name,
		InstanceName:      fmt.Sprintf("%s/%s", instance.Name, instance.ID),
		Index:             volumeIndex,
		Timeout:           d.Timeout(schema.TimeoutCreate),
	}

	info, err := storageAttachmentClient.CreateStorageAttachment(&input)
	if err != nil {
		return fmt.Errorf("Error creating StorageAttachment: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCStorageAttachmentRead(d, meta)
}

// Need to confirm that the index specified is not already in use.
func checkForEmptyIndex(attachments []compute.StorageAttachment, index int) bool {
	for _, attachment := range attachments {
		if attachment.Index == index {
			return false
		}
	}

	return true
}

func resourceOPCStorageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.StorageAttachments()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := compute.GetStorageAttachmentInput{
		Name: d.Id(),
	}

	result, err := resClient.GetStorageAttachment(&getInput)
	if err != nil {
		// StorageAttachment does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading storage_attachment %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of storage_attachment %s: %#v", d.Id(), result)
	d.Set("index", result.Index)
	d.Set("instance", strings.Split(result.InstanceName, "/")[0])
	d.Set("storage_volume", result.StorageVolumeName)
	return nil
}

func resourceOPCStorageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.StorageAttachments()
	name := d.Id()

	log.Printf("[DEBUG] Deleting StorageAttachment: %v", name)

	input := compute.DeleteStorageAttachmentInput{
		Name:    name,
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if err := resClient.DeleteStorageAttachment(&input); err != nil {
		return fmt.Errorf("Error deleting StorageAttachment")
	}
	return nil
}

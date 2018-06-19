package opc

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceStorageVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageVolumeSnapshotRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"collocated": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"machine_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parent_volume_bootable": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"property": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"start_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsComputedSchema(),

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*Client).computeClient.StorageVolumeSnapshots()

	name := d.Get("name").(string)
	input := &compute.GetStorageVolumeSnapshotInput{
		Name: name,
	}

	result, err := computeClient.GetStorageVolumeSnapshot(input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading storage volume snapshot '%s': %v", name, err)
	}

	if result == nil {
		// No Storage volume snapshot was found
		d.SetId("")
		return nil
	}

	d.SetId(result.Name)
	d.Set("volume_name", result.Volume)
	d.Set("description", result.Description)
	d.Set("name", result.Name)
	d.Set("property", result.Property)
	d.Set("platform", result.Platform)
	d.Set("account", result.Account)
	d.Set("machine_image_name", result.MachineImageName)
	d.Set("size", result.Size)
	d.Set("snapshot_timestamp", result.SnapshotTimestamp)
	d.Set("snapshot_id", result.SnapshotID)
	d.Set("start_timestamp", result.StartTimestamp)
	d.Set("status", result.Status)
	d.Set("status_detail", result.StatusDetail)
	d.Set("status_timestamp", result.StatusTimestamp)
	d.Set("uri", result.URI)

	bootable, err := strconv.ParseBool(result.ParentVolumeBootable)
	if err != nil {
		return fmt.Errorf("Error converting parent volume to boolean: %v", err)
	}
	d.Set("parent_volume_bootable", bootable)

	if result.Property != compute.SnapshotPropertyCollocated {
		d.Set("collocated", false)
	} else {
		d.Set("collocated", true)
	}

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	return nil
}

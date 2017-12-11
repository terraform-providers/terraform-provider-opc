package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCStorageVolume_Basic(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccStorageVolumeBasic, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_Complete(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccStorageVolumeComplete, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_MaxSize(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccStorageVolumeBasicMaxSize, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_Update(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccStorageVolumeComplete, ri)
	updatedConfig := fmt.Sprintf(testAccStorageVolumeUpdated, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_Bootable(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := testAccStorageVolumeBootable(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_BootableEndToEnd(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.restored"
	ri := acctest.RandInt()
	config := testAccStorageVolumeBootableEndToEnd(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_FromBootableSnapshot(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := testAccStorageVolumeFromBootableSnapshot(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
					resource.TestCheckResourceAttr(volumeResourceName, "size", "300"),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_ImageListEntry(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	ri := acctest.RandInt()
	config := testAccStorageVolumeImageListEntry(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_LowLatency(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageVolumeLowLatency(rInt),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
					resource.TestCheckResourceAttr(volumeResourceName, "storage_type", "/oracle/public/storage/latency"),
				),
			},
		},
	})
}

func TestAccOPCStorageVolume_FromSnapshot(t *testing.T) {
	volumeResourceName := "opc_compute_storage_volume.test"
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageVolumeFromSnapshot(rInt),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(volumeResourceName, testAccCheckStorageVolumeExists),
					resource.TestCheckResourceAttr(volumeResourceName, "name", fmt.Sprintf("test-acc-stor-vol-final-%d", rInt)),
					resource.TestCheckResourceAttrSet(volumeResourceName, "snapshot"),
					resource.TestCheckResourceAttrSet(volumeResourceName, "snapshot_id"),
					resource.TestCheckResourceAttr(volumeResourceName, "size", "5"),
				),
			},
		},
	})
}

func testAccCheckStorageVolumeExists(state *OPCResourceState) error {
	sv := state.ComputeClient.StorageVolumes()
	volumeName := state.Attributes["name"]

	input := compute.GetStorageVolumeInput{
		Name: volumeName,
	}
	info, err := sv.GetStorageVolume(&input)
	if err != nil {
		return fmt.Errorf("Error retrieving state of volume %s: %s", volumeName, err)
	}

	if info == nil {
		return fmt.Errorf("No info found for volume %s", volumeName)
	}

	return nil
}

func testAccCheckStorageVolumeDestroyed(state *OPCResourceState) error {
	sv := state.ComputeClient.StorageVolumes()

	volumeName := state.Attributes["name"]

	input := compute.GetStorageVolumeInput{
		Name: volumeName,
	}
	info, err := sv.GetStorageVolume(&input)
	if err != nil {
		return fmt.Errorf("Error retrieving state of volume %s: %s", volumeName, err)
	}

	if info != nil {
		return fmt.Errorf("Volume %s still exists", volumeName)
	}

	return nil
}

const testAccStorageVolumeBasic = `
resource "opc_compute_storage_volume" "test" {
  name = "test-acc-stor-vol-%d"
  size = 1
}
`

const testAccStorageVolumeComplete = `
resource "opc_compute_storage_volume" "test" {
  name        = "test-acc-stor-vol-%d"
  description = "Provider Acceptance Tests Storage Volume Initial"
  size        = 2
  tags        = ["foo"]
}
`

const testAccStorageVolumeUpdated = `
resource "opc_compute_storage_volume" "test" {
  name        = "test-acc-stor-vol-%d"
  description = "Provider Acceptance Tests Storage Volume Updated"
  size        = 2
  tags        = ["bar", "foo"]
}
`

func testAccStorageVolumeBootable(rInt int) string {
	return fmt.Sprintf(`
	resource "opc_compute_image_list" "test" {
	  name        = "test-acc-stor-vol-bootable-image-list-%d"
	  description = "Provider Acceptance Tests Storage Volume Bootable"
	}

	resource "opc_compute_image_list_entry" "test" {
	  name           = "${opc_compute_image_list.test.name}"
	  machine_images = [ "/oracle/public/oel_6.7_apaas_16.4.5_1610211300" ]
	  version        = 1
	}

	resource "opc_compute_storage_volume" "test" {
	  name             = "test-acc-stor-vol-bootable-%d"
	  description      = "Provider Acceptance Tests Storage Volume Bootable"
	  size             = 20
	  tags             = ["bar", "foo"]
	  bootable         = true
	  image_list       = "${opc_compute_image_list.test.name}"
	  image_list_entry = "${opc_compute_image_list_entry.test.version}"
	}
	`, rInt, rInt)
}

func testAccStorageVolumeBootableEndToEnd(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_image_list" "original" {
  name        = "test-acc-imagelist-%d"
  description = "Provider Acceptance Tests Original Image List"
}

resource "opc_compute_image_list_entry" "original" {
  name           = "${opc_compute_image_list.original.name}"
  machine_images = ["/oracle/public/oel_6.7_apaas_16.4.5_1610211300"]
  version        = 1
}

resource "opc_compute_storage_volume" "original" {
  name             = "test-acc-original-sv-%d"
  description      = "Provider Acceptance Tests Original Storage Volume"
  size             = 100
  bootable         = true
  image_list       = "${opc_compute_image_list.original.name}"
  image_list_entry = "${opc_compute_image_list_entry.original.version}"
}

resource "opc_compute_instance" "original" {
  name       = "test-acc-original-instance-%d"
  label      = "Original Instance"
  shape      = "oc3"
  image_list = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"

  storage {
    volume = "${opc_compute_storage_volume.original.name}"
    index  = 1
  }

  networking_info {
    index          = 0
    nat            = ["ippool:/oracle/public/ippool"]
    shared_network = true
  }
}

# Take a snapshot
resource "opc_compute_storage_volume_snapshot" "test" {
  name                   = "test-acc-snapshot-%d"
  description            = "Provider Acceptance Tests Snapshot"
  tags                   = ["example"]
  collocated             = true
  volume_name            = "${opc_compute_storage_volume.original.name}"
  parent_volume_bootable = true
}

# Restore the image onto the Volume
resource "opc_compute_storage_volume" "restored" {
  name        = "test-acc-volume-%d"
  description = "storage volume from snapshot"
  size        = 100
  bootable    = true
  snapshot_id = "${opc_compute_storage_volume_snapshot.test.snapshot_id}"
}

# Ensure it's bootable
resource "opc_compute_instance" "restored" {
  name       = "test-acc-restored-%d"
  label      = "Provider Acceptance Test Restored Instance"
  shape      = "oc3"
  image_list = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"

  storage {
    volume = "${opc_compute_storage_volume.restored.name}"
    index  = 1
  }

  networking_info {
    index          = 0
    nat            = ["ippool:/oracle/public/ippool"]
    shared_network = true
  }
}`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccStorageVolumeImageListEntry(rInt int) string {
	return fmt.Sprintf(`
	resource "opc_compute_image_list" "test" {
	  name        = "test-acc-stor-vol-bootable-image-list-%d"
	  description = "Provider Acceptance Tests Storage Volume Image List Entry"
	}

	resource "opc_compute_image_list_entry" "test" {
	  name           = "${opc_compute_image_list.test.name}"
	  machine_images = [ "/oracle/public/oel_6.7_apaas_16.4.5_1610211300" ]
	  version        = 1
	}

	resource "opc_compute_storage_volume" "test" {
	  name             = "test-acc-stor-vol-image-list-entry-%d"
	  description      = "Provider Acceptance Tests Storage Volume Image List Entry"
	  size             = 20
	  tags             = ["bar", "foo"]
	  image_list_entry = "${opc_compute_image_list_entry.test.version}"
	}
	`, rInt, rInt)
}

const testAccStorageVolumeBasicMaxSize = `
resource "opc_compute_storage_volume" "test" {
  name        = "test-acc-stor-vol-%d"
  description = "Provider Acceptance Tests Storage Volume Max Size"
  size        = 2048
}
`

func testAccStorageVolumeFromSnapshot(rInt int) string {
	return fmt.Sprintf(`
  // Initial Storage Volume to create snapshot with
  resource "opc_compute_storage_volume" "foo" {
    name        = "test-acc-stor-vol-%d"
    description = "Acc Test intermediary storage volume for snapshot"
    size        = 5
  }

  resource "opc_compute_storage_volume_snapshot" "foo" {
    description = "testing-acc"
    name        = "test-acc-stor-snapshot-%d"
    collocated  = true
    volume_name = "${opc_compute_storage_volume.foo.name}"
  }

  // Create storage volume from snapshot
  resource "opc_compute_storage_volume" "test" {
    name        = "test-acc-stor-vol-final-%d"
    description = "storage volume from snapshot"
    size        = 5
    snapshot_id = "${opc_compute_storage_volume_snapshot.foo.snapshot_id}"
  }`, rInt, rInt, rInt)
}

func testAccStorageVolumeLowLatency(rInt int) string {
	return fmt.Sprintf(`
  resource "opc_compute_storage_volume" "test" {
    name         = "test-acc-stor-vol-ll-%d"
    description  = "Acc Test Storage Volume Low Latency"
    storage_type = "/oracle/public/storage/latency"
    size         = 5
  }`, rInt)
}

func testAccStorageVolumeFromBootableSnapshot(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_image_list" "test" {
  name        = "test-acc-image-list-%d"
  default     = 21
  description = "description"
}

resource "opc_compute_image_list_entry" "test" {
  name           = "${opc_compute_image_list.test.name}"
  machine_images = [ "/oracle/public/oel_6.7_apaas_16.4.5_1610211300" ]
  version        = 1
}
resource "opc_compute_storage_volume" "foo" {
  name = "test-acc-stor-vol-%d"
  description = "testAccStorageVolumeSnapshot_basic"
  size = 300
	bootable = true
  image_list = "${opc_compute_image_list.test.name}"
  image_list_entry = "${opc_compute_image_list_entry.test.version}"
}

resource "opc_compute_storage_volume_snapshot" "test" {
  name = "test-acc-stor-vol-snapshot-%d"
  description = "storage volume snapshot"
  volume_name = "${opc_compute_storage_volume.foo.name}"
	parent_volume_bootable = true
}

resource "opc_compute_storage_volume" "test" {
  name = "test-acc-stor-vol-from-snapshot-%d"
  snapshot_id = "${opc_compute_storage_volume_snapshot.test.snapshot_id}"
  size = 300
}
`, rInt, rInt, rInt, rInt)
}

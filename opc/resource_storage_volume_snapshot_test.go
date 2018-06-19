package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCStorageVolumeSnapshot_basic(t *testing.T) {
	snapshotName := "opc_compute_storage_volume_snapshot.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(snapshotName, testAccCheckStorageVolumeSnapshotDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageVolumeSnapshot_basic(rInt),
				Check: resource.ComposeTestCheckFunc(opcResourceCheck(snapshotName, testAccCheckStorageVolumeSnapshotExists),
					resource.TestCheckResourceAttr(snapshotName, "name", fmt.Sprintf("test-acc-stor-vol-%d", rInt)),
					resource.TestCheckResourceAttr(snapshotName, "parent_volume_bootable", "false"),
					resource.TestCheckResourceAttr(snapshotName, "collocated", "true"),
					resource.TestCheckResourceAttr(snapshotName, "size", "5"),
				),
			},
		},
	})
}

func TestAccOPCStorageVolumeSnapshot_bootableSnapshot(t *testing.T) {
	snapshotName := "opc_compute_storage_volume_snapshot.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(snapshotName, testAccCheckStorageVolumeSnapshotDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageVolumeSnapshot_bootableSnapshot(rInt),
				Check: resource.ComposeTestCheckFunc(opcResourceCheck(snapshotName, testAccCheckStorageVolumeSnapshotExists),
					resource.TestCheckResourceAttr(snapshotName, "name", fmt.Sprintf("test-acc-stor-vol-snapshot-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckStorageVolumeSnapshotExists(state *OPCResourceState) error {
	client := state.Client.StorageVolumeSnapshots()
	snapshotName := state.Attributes["name"]

	input := &compute.GetStorageVolumeSnapshotInput{
		Name: snapshotName,
	}

	info, err := client.GetStorageVolumeSnapshot(input)
	if err != nil {
		return fmt.Errorf("Error retrieving state of snapshot '%s': %v", snapshotName, err)
	}

	if info == nil {
		return fmt.Errorf("No info found for snapshot '%s'", snapshotName)
	}

	return nil
}

func testAccCheckStorageVolumeSnapshotDestroyed(state *OPCResourceState) error {
	client := state.Client.StorageVolumeSnapshots()
	snapshotName := state.Attributes["name"]

	input := &compute.GetStorageVolumeSnapshotInput{
		Name: snapshotName,
	}
	info, err := client.GetStorageVolumeSnapshot(input)
	if err != nil {
		return fmt.Errorf("Error retrieving state of snapshot '%s': %v", snapshotName, err)
	}

	if info != nil {
		return fmt.Errorf("Snapshot '%s' still exists", snapshotName)
	}

	return nil
}

func testAccStorageVolumeSnapshot_basic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_storage_volume" "foo" {
  name = "test-acc-stor-vol-%d"
  description = "testAccStorageVolumeSnapshot_basic"
  size = 5
}

resource "opc_compute_storage_volume_snapshot" "test" {
  name = "test-acc-stor-vol-%d"
  description = "storage volume snapshot"
  collocated = true
  volume_name = "${opc_compute_storage_volume.foo.name}"
}
`, rInt, rInt)
}

func testAccStorageVolumeSnapshot_bootableSnapshot(rInt int) string {
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
  size = 50
	tags        = ["bar", "foo"]
	bootable = true
  image_list = "${opc_compute_image_list.test.name}"
  image_list_entry = "${opc_compute_image_list_entry.test.version}"
}

resource "opc_compute_storage_volume_snapshot" "test" {
  name = "test-acc-stor-vol-snapshot-%d"
  description = "storage volume snapshot"
  collocated = true
  volume_name = "${opc_compute_storage_volume.foo.name}"
	parent_volume_bootable = true
}
`, rInt, rInt, rInt)
}

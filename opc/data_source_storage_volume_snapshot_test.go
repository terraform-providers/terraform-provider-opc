package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDataSourceStorageVolumeSnapshot_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "data.opc_compute_storage_volume_snapshot.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageVolumeSnapshotBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "collocated", "true"),
					resource.TestCheckResourceAttr(resName, "description", "storage volume snapshot"),
					resource.TestCheckResourceAttr(resName, "volume_name", fmt.Sprintf("test-acc-stor-vol-%d", rInt)),
					resource.TestCheckResourceAttr(resName, "size", "5"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "parent_volume_bootable"),
					resource.TestCheckResourceAttrSet(resName, "property"),
					resource.TestCheckResourceAttrSet(resName, "snapshot_timestamp"),
					resource.TestCheckResourceAttrSet(resName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resName, "start_timestamp"),
					resource.TestCheckResourceAttrSet(resName, "uri"),
					resource.TestCheckResourceAttrSet(resName, "volume_name"),
				),
			},
		},
	})
}

func testAccDataSourceStorageVolumeSnapshotBasic(rInt int) string {
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
  tags = ["foo"]
}

data "opc_compute_storage_volume_snapshot" "test" {
  name = "${opc_compute_storage_volume_snapshot.test.name}"
}`, rInt, rInt)
}

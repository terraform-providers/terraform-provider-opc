package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const _TestAccSnapshotImage = "/oracle/public/OL_5.11_UEKR2_x86_64"

func TestAccOPCSnapshot_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSnapshotBasic(rInt),
				Check:  testAccCheckSnapshotExists,
			},
		},
	})
}

func TestAccOPCSnapshot_MachineImage(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSnapshotMachineImage(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists,
					resource.TestCheckResourceAttr(
						"opc_compute_snapshot.test", "machine_image", fmt.Sprintf("acc-test-snapshot-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckSnapshotExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).computeClient.Snapshots()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_snapshot" {
			continue
		}

		input := compute.GetSnapshotInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetSnapshot(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Snapshot %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).computeClient.Snapshots()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_snapshot" {
			continue
		}

		input := compute.GetSnapshotInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetSnapshot(&input); err == nil {
			return fmt.Errorf("Snapshot %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccOPCSnapshotBasic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_instance" "test" {
	name = "acc-test-snapshot-%d"
  label = "TestAccOPCSnapshot_basic"
	shape = "oc3"
	image_list = "%s"
}

resource "opc_compute_snapshot" "test" {
  instance = "${opc_compute_instance.test.name}/${opc_compute_instance.test.id}"
}`, rInt, _TestAccSnapshotImage)
}

func testAccOPCSnapshotMachineImage(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_instance" "test" {
	name = "acc-test-snapshot-%d"
  label = "TestAccOPCSnapshot_basic"
	shape = "oc3"
	image_list = "%s"
}

resource "opc_compute_snapshot" "test" {
  instance = "${opc_compute_instance.test.name}/${opc_compute_instance.test.id}"
  machine_image = "acc-test-snapshot-%d"
}`, rInt, _TestAccSnapshotImage, rInt)
}

package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCSnapshot_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccOPCSnapshotBasic, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testAccCheckSnapshotExists,
			},
		},
	})
}

func TestAccOPCSnapshot_MachineImage(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccOPCSnapshotMachineImage, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists,
					resource.TestCheckResourceAttr(
						"opc_compute_snapshot.test", "machine_image", fmt.Sprintf("acc-test-snapshot-%d", ri)),
				),
			},
		},
	})
}

func testAccCheckSnapshotExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*compute.ComputeClient).Snapshots()
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
	client := testAccProvider.Meta().(*compute.ComputeClient).Snapshots()

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

var testAccOPCSnapshotBasic = `
resource "opc_compute_instance" "test" {
	name = "acc-test-snapshot-%d"
  label = "TestAccOPCSnapshot_basic"
	shape = "oc3"
	image_list = "/oracle/public/JEOS_OL_6.6_10GB_RD-1.2.217-20151201-194209"
}

resource "opc_compute_snapshot" "test" {
  instance = "${opc_compute_instance.test.name}/${opc_compute_instance.test.id}"
}
`

var testAccOPCSnapshotMachineImage = `
resource "opc_compute_instance" "test" {
	name = "acc-test-snapshot-%d"
  label = "TestAccOPCSnapshot_basic"
	shape = "oc3"
	image_list = "/oracle/public/JEOS_OL_6.6_10GB_RD-1.2.217-20151201-194209"
}

resource "opc_compute_snapshot" "test" {
  instance = "${opc_compute_instance.test.name}/${opc_compute_instance.test.id}"
  machine_image = "acc-test-snapshot-%d"
}
`

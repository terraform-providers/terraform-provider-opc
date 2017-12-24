package opc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCMachineImage_Basic(t *testing.T) {

	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMachineImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMachineImage_basic(ri),
				Check:  testAccCheckMachineImageExists,
			},
		},
	})
}

func testAccCheckMachineImageExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.MachineImages()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_machine_image" {
			continue
		}

		input := compute.GetMachineImageInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetMachineImage(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Machine Image %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckMachineImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.MachineImages()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_machine_image" {
			continue
		}

		input := compute.GetMachineImageInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetMachineImage(&input); err == nil {
			return fmt.Errorf("Machine Image %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccMachineImage_basic(rInt int) string {

	identity_domain := os.Getenv("OPC_IDENTITY_DOMAIN")

	testAccMachineImageBasic := `
  resource "opc_storage_object" "acc-test-machine-image" {
		name         = "acc-test-machine-image.tar.gz"
		container    = "compute_images"
		file         = "test-fixtures/dummy.tar.gz"
		content_type = "application/tar+gzip;charset=UTF-8"
	}

	resource "opc_compute_machine_image" "test" {
		account = "/Compute-%s/cloud_storage"
	  name    = "acc-test-machine-image-%d"
	  file    = "${opc_storage_object.acc-test-machine-image.name}"
	}`

	return fmt.Sprintf(testAccMachineImageBasic, identity_domain, rInt)
}

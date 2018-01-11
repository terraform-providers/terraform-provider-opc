package opc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDataSourceMachineImage_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "data.opc_compute_machine_image.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMachineImageBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "file", "acc-test-machine-image.tar.gz"),
					resource.TestCheckResourceAttr(resName, "description", "machine image description"),
					resource.TestCheckResourceAttr(resName, "attributes", "{\"role\":\"master\"}"),
					resource.TestCheckResourceAttr(resName, "state", "available"),
				),
			},
		},
	})
}

func testAccDataSourceMachineImageBasic(rInt int) string {
	identity_domain := os.Getenv("OPC_IDENTITY_DOMAIN")

	testAccMachineImageBasic := `
  resource "opc_storage_object" "acc-test-machine-image" {
		name         = "acc-test-machine-image.tar.gz"
		container    = "compute_images"
		file         = "test-fixtures/dummy.tar.gz"
		content_type = "application/tar+gzip;charset=UTF-8"
	}

	resource "opc_compute_machine_image" "test" {
		account     = "/Compute-%s/cloud_storage"
	  name        = "acc-test-machine-image-%d"
	  file        = "${opc_storage_object.acc-test-machine-image.name}"
		description = "machine image description"
		attributes  = "{\"role\":\"master\"}"

	}

	data "opc_compute_machine_image" "test" {
		account = "${opc_compute_machine_image.test.account}"
		name    = "${opc_compute_machine_image.test.name}"
	}`

	return fmt.Sprintf(testAccMachineImageBasic, identity_domain, rInt)
}

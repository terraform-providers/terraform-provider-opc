package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDataSourceImageListEntry_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "data.opc_compute_image_list_entry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImageListEntryBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "machine_images.#", "3"),
					resource.TestCheckResourceAttr(resName, "machine_images.1",
						"/oracle/public/oel_6.7_apaas_16.4.5_1610211300"),
					resource.TestCheckResourceAttr(resName, "machine_images.2",
						"/oracle/public/OL_5.11_UEKR2_i386-17.2.2-20170405-205607"),
					resource.TestCheckResourceAttr(resName, "version", "1"),
				),
			},
		},
	})
}

func TestAccOPCDataSourceImageListEntry_entry(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "data.opc_compute_image_list_entry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImageListEntry_entry(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "machine_images.#", "1"),
					resource.TestCheckResourceAttr(resName, "machine_images.0",
						"/oracle/public/OL_5.11_UEKR2_i386-17.2.2-20170405-205607"),
					resource.TestCheckResourceAttr(resName, "version", "1"),
					resource.TestCheckResourceAttr(resName, "entry", "3"),
				),
			},
		},
	})
}

func testAccDataSourceImageListEntryBasic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_image_list" "test" {
 name        = "test-acc-image-list-entry-basic-%d"
 description = "Acceptance Test TestAccOPCImageListEntry_Basic"
 default     = 1
}

resource "opc_compute_image_list_entry" "test" {
  name           = "${opc_compute_image_list.test.name}"
  machine_images = [
    "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
    "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
    "/oracle/public/OL_5.11_UEKR2_i386-17.2.2-20170405-205607"
  ]
  version        = 1
}

data "opc_compute_image_list_entry" "test" {
  image_list = "${opc_compute_image_list_entry.test.name}"
  version    = 1
}`, rInt)
}

func testAccDataSourceImageListEntry_entry(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_image_list" "test" {
 name        = "test-acc-image-list-entry-basic-%d"
 description = "Acceptance Test TestAccOPCImageListEntry_Basic"
 default     = 1
}

resource "opc_compute_image_list_entry" "test" {
  name           = "${opc_compute_image_list.test.name}"
  machine_images = [
    "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
    "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
    "/oracle/public/OL_5.11_UEKR2_i386-17.2.2-20170405-205607"
  ]
  version        = 1
}

data "opc_compute_image_list_entry" "test" {
  image_list = "${opc_compute_image_list_entry.test.name}"
  version    = 1
  entry      = 3
}`, rInt)
}

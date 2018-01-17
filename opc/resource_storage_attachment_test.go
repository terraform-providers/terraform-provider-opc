package opc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCStorageAttachment_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccStorageAttachmentBasic(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageAttachmentExists,
				),
			},
		},
	})
}

func TestAccOPCStorageAttachment_InvalidIndex(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccStorageAttachmentInvalidIndex(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Storage index 1 is already in use"),
			},
		},
	})
}

func TestAccOPCStorageAttachment_BootAndAttached(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccStorageAttachmentBootAndAttached(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageAttachmentExists,
				),
			},
		},
	})
}

func testAccCheckStorageAttachmentExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.StorageAttachments()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_storage_attachment" {
			continue
		}

		input := compute.GetStorageAttachmentInput{
			Name: rs.Primary.Attributes["id"],
		}
		if _, err := client.GetStorageAttachment(&input); err != nil {
			return fmt.Errorf("Error retrieving state of StorageAttachment %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckStorageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.StorageAttachments()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_storage_attachment" {
			continue
		}

		input := compute.GetStorageAttachmentInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetStorageAttachment(&input); err == nil {
			return fmt.Errorf("StorageAttachment %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccStorageAttachmentBasic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_storage_volume" "foo" {
  name = "acc-test-storage-attachment-%d"
  size = 1
}

resource "opc_compute_instance" "test" {
	name = "acc-test-storage-attachment-%d"
	label = "TestAccOPCInstance_basic"
	shape = "oc3"
	image_list = "%s"
}

resource "opc_compute_storage_attachment" "test" {
  instance = "${opc_compute_instance.test.name}"
  storage_volume = "${opc_compute_storage_volume.foo.name}"
  index = 1
}
`, rInt, rInt, TEST_IMAGE_LIST)
}

func testAccStorageAttachmentInvalidIndex(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_storage_volume" "foo" {
  name = "acc-test-storage-attachment-%d"
  size = 1
}

resource "opc_compute_storage_volume" "bar" {
  name = "acc-test-storage-attachment-attached-%d"
  size = 1
}

resource "opc_compute_instance" "test" {
	name = "acc-test-storage-attachment-%d"
	label = "TestAccOPCInstance_basic"
	shape = "oc3"
	image_list = "%s"
  storage {
		volume = "${opc_compute_storage_volume.bar.name}"
		index = 1
	}
}

resource "opc_compute_storage_attachment" "test" {
  instance = "${opc_compute_instance.test.name}"
  storage_volume = "${opc_compute_storage_volume.foo.name}"
  index = 1
}
`, rInt, rInt, rInt, TEST_IMAGE_LIST)
}

func testAccStorageAttachmentBootAndAttached(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_storage_volume" "foo" {
  name = "acc-test-storage-attachment-%d"
  size = 1
}

resource "opc_compute_image_list" "test" {
  name = "acc-test-instance-%d"
  description = "testing instance start-stop"
}

resource "opc_compute_image_list_entry" "test" {
  name = "${opc_compute_image_list.test.name}"
	machine_images = [ "/oracle/public/oel_6.7_apaas_16.4.5_1610211300" ]
  version = 1
}

resource "opc_compute_storage_volume" "boot" {
  name = "acc-test-instance-%d"
  size = "20"
  image_list = "${opc_compute_image_list.test.name}"
  image_list_entry = "${opc_compute_image_list_entry.test.version}"
  bootable = true
}

resource "opc_compute_instance" "test" {
	name = "acc-test-storage-attachment-%d"
	label = "TestAccOPCInstance_basic"
	shape = "oc3"
	image_list = "%s"
  storage {
		volume = "${opc_compute_storage_volume.boot.name}"
		index = 1
	}
}

resource "opc_compute_storage_attachment" "test" {
  instance = "${opc_compute_instance.test.name}"
  storage_volume = "${opc_compute_storage_volume.foo.name}"
  index = 2
}
`, rInt, rInt, rInt, rInt, TEST_IMAGE_LIST)
}

package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCStorageContainer_Basic(t *testing.T) {
	containerResourceName := "opc_storage_container.test"
	ri := acctest.RandInt()
	config := testAccOPCStorageContainerBasic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageContainerExists,
					resource.TestCheckResourceAttr(containerResourceName, "max_age", "50"),
					resource.TestCheckResourceAttr(containerResourceName, "primary_key", "test-key"),
				),
			},
		},
	})
}

func TestAccOPCStorageContainer_Updated(t *testing.T) {
	containerResourceName := "opc_storage_container.test"
	ri := acctest.RandInt()
	config := testAccOPCStorageContainerBasic(ri)
	config2 := testAccOPCStorageContainerUpdated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageContainerExists,
					resource.TestCheckResourceAttr(containerResourceName, "max_age", "50"),
					resource.TestCheckResourceAttr(containerResourceName, "primary_key", "test-key"),
					resource.TestCheckResourceAttr(containerResourceName, "allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(containerResourceName, "allowed_origins.0", "origin-1"),
					resource.TestCheckResourceAttr(containerResourceName, "exposed_headers.#", "1"),
					resource.TestCheckResourceAttr(containerResourceName, "exposed_headers.0", "exposed-header-1"),
					resource.TestCheckResourceAttr(containerResourceName, "quota_bytes", "1000000000"),
					resource.TestCheckResourceAttr(containerResourceName, "quota_count", "1000"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.Foo", "bar"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.Abc-Def", "xyz"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageContainerExists,
					resource.TestCheckResourceAttr(containerResourceName, "max_age", "60"),
					resource.TestCheckResourceAttr(containerResourceName, "primary_key", "test-key-updated"),
					resource.TestCheckResourceAttr(containerResourceName, "secondary_key", "test-key"),
					resource.TestCheckResourceAttr(containerResourceName, "allowed_origins.#", "2"),
					resource.TestCheckResourceAttr(containerResourceName, "allowed_origins.1", "origin-2"),
					resource.TestCheckResourceAttr(containerResourceName, "exposed_headers.#", "2"),
					resource.TestCheckResourceAttr(containerResourceName, "exposed_headers.1", "exposed-header-2"),
					resource.TestCheckResourceAttr(containerResourceName, "quota_bytes", "2000000000"),
					resource.TestCheckResourceAttr(containerResourceName, "quota_count", "2000"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.Bar", "foo"),
					resource.TestCheckResourceAttr(containerResourceName, "metadata.Abc-Def", "xyz"),
				),
			},
		},
	})
}

func testAccCheckStorageContainerExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).storageClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_storage_container" {
			continue
		}

		input := storage.GetContainerInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetContainer(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Storage Container %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckStorageContainerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).storageClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_storage_container" {
			continue
		}

		input := storage.GetContainerInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetContainer(&input); err == nil {
			return fmt.Errorf("Storage Container %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccOPCStorageContainerBasic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_storage_container" "test" {
  name = "acc-storage-container-%d"
  max_age = 50
	quota_bytes = 1000000000
	quota_count = 1000
  primary_key = "test-key"
  allowed_origins = ["origin-1"]
	exposed_headers = ["exposed-header-1"]
	metadata {
		"Foo" = "bar",
		"Abc-Def" = "xyz"
	}
}
`, rInt)
}

func testAccOPCStorageContainerUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "opc_storage_container" "test" {
  name = "acc-storage-container-%d"
  max_age = 60
	quota_bytes = 2000000000
	quota_count = 2000
  primary_key = "test-key-updated"
  secondary_key = "test-key"
  allowed_origins = ["origin-1", "origin-2"]
	exposed_headers = ["exposed-header-1", "exposed-header-2"]
	metadata {
		"Bar" = "foo",
		"Abc-Def" = "xyz"
	}
}
`, rInt)
}

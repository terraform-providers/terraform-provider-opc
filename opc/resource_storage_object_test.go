package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const _TestStorageObjectPath = "test-fixtures"

func TestAccOPCStorageObject_contentSource(t *testing.T) {
	resName := "opc_storage_object.test"
	rInt := acctest.RandInt()

	body := _SourceInput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCStorageObject_contentSource(rInt, body),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageObjectExists,
					resource.TestCheckResourceAttr(resName, "name", fmt.Sprintf("test-acc-%d", rInt)),
					resource.TestCheckResourceAttr(resName, "container", fmt.Sprintf("acc-test-%d", rInt)),
					resource.TestCheckResourceAttrSet(resName, "content_length"),
					resource.TestCheckResourceAttr(resName, "content_type", "text/plain;charset=UTF-8"),
					resource.TestCheckResourceAttr(resName, "delete_at", "0"),
					resource.TestCheckResourceAttrSet(resName, "last_modified"),
					resource.TestCheckResourceAttrSet(resName, "timestamp"),
					resource.TestCheckResourceAttrSet(resName, "transaction_id"),
				),
			},
		},
	})
}

func TestAccOPCStorageObject_contentSource_nilContentType(t *testing.T) {
	resName := "opc_storage_object.test"
	rInt := acctest.RandInt()

	body := _SourceInput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCStorageObject_contentSource_nilContentType(rInt, body),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageObjectExists,
					resource.TestCheckResourceAttr(resName, "name", fmt.Sprintf("test-acc-%d", rInt)),
					resource.TestCheckResourceAttr(resName, "container", fmt.Sprintf("acc-test-%d", rInt)),
					resource.TestCheckResourceAttrSet(resName, "content_length"),
					resource.TestCheckResourceAttrSet(resName, "content_type"),
					resource.TestCheckResourceAttr(resName, "delete_at", "0"),
					resource.TestCheckResourceAttrSet(resName, "last_modified"),
					resource.TestCheckResourceAttrSet(resName, "timestamp"),
					resource.TestCheckResourceAttrSet(resName, "transaction_id"),
				),
			},
		},
	})
}

func TestAccOPCStorageObject_fileSource(t *testing.T) {
	resName := "opc_storage_object.test"
	rInt := acctest.RandInt()

	path := _TestStorageObjectPath + "/fileSource.txt"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCStorageObject_fileSource(rInt, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageObjectExists,
					resource.TestCheckResourceAttr(resName, "name", fmt.Sprintf("test-acc-%d", rInt)),
					resource.TestCheckResourceAttr(resName, "container", fmt.Sprintf("acc-test-%d", rInt)),
					resource.TestCheckResourceAttrSet(resName, "content_length"),
					resource.TestCheckResourceAttr(resName, "content_type", "text/plain;charset=UTF-8"),
					resource.TestCheckResourceAttr(resName, "delete_at", "0"),
					resource.TestCheckResourceAttrSet(resName, "last_modified"),
					resource.TestCheckResourceAttrSet(resName, "timestamp"),
					resource.TestCheckResourceAttrSet(resName, "transaction_id"),
				),
			},
		},
	})
}

func testAccCheckStorageObjectExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).storageClient.Objects()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_storage_object" {
			continue
		}

		input := &storage.GetObjectInput{
			ID: rs.Primary.Attributes["id"],
		}
		if _, err := client.GetObject(input); err != nil {
			return fmt.Errorf("Error retrieving state of Storage Object (%s): %s", input.ID, err)
		}
	}
	return nil
}

func testAccCheckStorageObjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).storageClient.Objects()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_storage_object" {
			continue
		}

		input := &storage.GetObjectInput{
			ID: rs.Primary.Attributes["id"],
		}

		if info, err := client.GetObject(input); err == nil {
			return fmt.Errorf("Storage Object (%s) still exists: %#v", input.ID, info)
		}
	}
	return nil
}

func testAccOPCStorageObject_testContainer(rInt int) string {
	return fmt.Sprintf(`
resource "opc_storage_container" "foo" {
  name = "acc-test-%d"
  max_age = 50
  primary_key = "test-key"
  allowed_origins = ["origin-1"]
}`, rInt)
}

func testAccOPCStorageObject_contentSource(rInt int, body string) string {
	return fmt.Sprintf(`
%s

resource "opc_storage_object" "test" {
  name = "test-acc-%d"
  container = "${opc_storage_container.foo.name}"
  content_type = "text/plain;charset=UTF-8"
  content = <<EOF
%s
EOF
}`,
		testAccOPCStorageObject_testContainer(rInt),
		rInt,
		body)
}

func testAccOPCStorageObject_contentSource_nilContentType(rInt int, body string) string {
	return fmt.Sprintf(`
%s

resource "opc_storage_object" "test" {
  name = "test-acc-%d"
  container = "${opc_storage_container.foo.name}"
  content = <<EOF
%s
EOF
}`,
		testAccOPCStorageObject_testContainer(rInt),
		rInt,
		body)
}

func testAccOPCStorageObject_fileSource(rInt int, path string) string {
	return fmt.Sprintf(`
%s

resource "opc_storage_object" "test" {
  name = "test-acc-%d"
  container = "${opc_storage_container.foo.name}"
  content_type = "text/plain;charset=UTF-8"
  file = "%s"
}`,
		testAccOPCStorageObject_testContainer(rInt),
		rInt,
		path)
}

const _SourceInput = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi auctor nisi id sem gravida, quis sollicitudin dolor
maximus. Sed est lectus, mollis sit amet neque eu, pulvinar aliquet turpis. Aenean in euismod erat. Proin pulvinar
ex vel lorem malesuada, sed tincidunt urna posuere. Sed fringilla, elit et faucibus maximus, dui orci blandit lectus,
ullamcorper fringilla felis nisl at nisl. Ut leo elit, semper non dui sit amet, sagittis commodo nulla. Nulla pulvinar
purus a nunc pellentesque scelerisque at id elit. Etiam quis bibendum eros. Etiam erat elit, feugiat non ante tempus,
mattis consectetur purus. Cras nunc nibh, fringilla in imperdiet a, tempus porta nisl. Curabitur nec justo nec leo
luctus scelerisque quis sit amet risus. Curabitur finibus fringilla lacus eu vestibulum. Nunc pellentesque aliquam
semper. Proin nec ligula urna. Donec lobortis aliquam nunc vitae feugiat. Integer blandit risus in gravida facilisis.
Pellentesque vitae lectus sed est pretium finibus. Morbi sed lacus purus. Duis nec condimentum urna. Donec vel velit
purus. Ut a velit risus. Vivamus ac euismod magna, eget convallis quam. Sed tincidunt, nisl nec rhoncus facilisis,
orci mauris commodo leo, ut eleifend nisi nisi sit amet mauris. Ut lacinia viverra rhoncus. Phasellus lacinia eleifend
turpis eu rutrum. Donec sed gravida eros, eget molestie ipsum.In hac habitasse platea dictumst. Duis a libero ante.
Quisque euismod placerat risus sit amet maximus. Praesent malesuada velit nec dui tincidunt rutrum. Proin commodo ex
non consectetur cursus. Pellentesque egestas pharetra mauris, et condimentum nibh rhoncus nec. Morbi hendrerit vel
ligula vel varius. Vestibulum in faucibus metus, eget euismod justo. Cras dolor sem, dictum eget scelerisque at,
scelerisque eu enim. Aliquam vulputate rutrum orci, vitae convallis mauris sollicitudin ut.Quisque eu accumsan massa.
Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Ut nisl magna, vulputate eget
eleifend id, tincidunt eget dolor. Nam pulvinar, dui non pellentesque dignissim, nulla neque iaculis sapien, id commodo
nisi nunc vel turpis. Vivamus eget dapibus lacus. Mauris convallis mi sit amet faucibus placerat. Mauris gravida neque
tortor, vel placerat sem elementum venenatis. Integer eu placerat est. Sed sem massa, volutpat eget augue eget, aliquam
semper sem.`

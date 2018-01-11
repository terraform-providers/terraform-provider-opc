package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCVNICSet_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := fmt.Sprintf("testing-acc-%d", rInt)
	rDesc := fmt.Sprintf("acctesting vnic set %d", rInt)
	resourceName := "opc_compute_vnic_set.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOPCCheckVNICSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVnicSetBasic(rName, rDesc, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckVNICSetExists,
					resource.TestCheckResourceAttr(
						resourceName, "name", rName),
					resource.TestCheckResourceAttr(
						resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(
						resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "virtual_nics.#", "2"),
				),
			},
			{
				Config: testAccVnicSetBasic_Update(rName, rDesc, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckVNICSetExists,
					resource.TestCheckResourceAttr(
						resourceName, "name", rName),
					resource.TestCheckResourceAttr(
						resourceName, "description", fmt.Sprintf("%s-updated", rDesc)),
					resource.TestCheckResourceAttr(
						resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "virtual_nics.#", "2"),
				),
			},
		},
	})
}

// Setting this takes two applies. This is... not "optimal"
// However, fixing this would require some core changes that
// allow for multiple level dependencies, ie: refresh after apply
// for a resource dependency loop.
func TestAccOPCVNICSet_UpdateFromInstance(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opc_compute_vnic_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOPCCheckVNICSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVnicSetInstance(rInt),
			},
			{
				Config: testAccVnicSetInstance(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckVNICSetExists,
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("vnicset-acctest-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "virtual_nics.#", "1"),
				),
			},
		},
	})
}

func testAccOPCCheckVNICSetExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.VirtNICSets()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_vnic_set" {
			continue
		}

		input := compute.GetVirtualNICSetInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetVirtualNICSet(&input); err != nil {
			return fmt.Errorf("Error retrieving state of VNIC Set %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccOPCCheckVNICSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.VirtNICSets()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_vnic_set" {
			continue
		}

		input := compute.GetVirtualNICSetInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetVirtualNICSet(&input); err == nil {
			return fmt.Errorf("VNIC Set %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccVnicSetBasic(rName, rDesc string, rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_ip_network" "foo" {
  name = "testing-vnic-set-%d"
  description = "testing-vnic-set"
  ip_address_prefix = "10.1.14.0/24"
}

resource "opc_compute_ip_network" "bar" {
  name = "testing-vnic-set2-%d"
  description = "testing-vnic-set2"
  ip_address_prefix = "10.1.15.0/24"
}

resource "opc_compute_instance" "foo" {
  name = "test-vnic-set-%d"
  label = "testing"
  shape = "oc3"
  image_list = "%s"
  networking_info {
    index = 0
    ip_network = "${opc_compute_ip_network.foo.id}"
    vnic = "test-vnic-set-%d"
    shared_network = false
  }
  networking_info {
    index = 1
    ip_network = "${opc_compute_ip_network.bar.id}"
    vnic = "test-vnic-set2-%d"
    shared_network = false
  }
}

data "opc_compute_network_interface" "foo" {
  instance_name = "${opc_compute_instance.foo.name}"
  instance_id = "${opc_compute_instance.foo.id}"
  interface = "eth0"
}

data "opc_compute_network_interface" "bar" {
  instance_name = "${opc_compute_instance.foo.name}"
  instance_id = "${opc_compute_instance.foo.id}"
  interface = "eth1"
}

resource "opc_compute_vnic_set" "test" {
  name = "%s"
  description = "%s"
  tags = ["tag1", "tag2"]
  virtual_nics = [
    "${data.opc_compute_network_interface.foo.vnic}",
    "${data.opc_compute_network_interface.bar.vnic}",
  ]
}`, rInt, rInt, rInt, TEST_IMAGE_LIST, rInt, rInt, rName, rDesc)
}

func testAccVnicSetBasic_Update(rName, rDesc string, rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_ip_network" "foo" {
  name = "testing-vnic-set-%d"
  description = "testing-vnic-set"
  ip_address_prefix = "10.1.14.0/24"
}

resource "opc_compute_ip_network" "bar" {
  name = "testing-vnic-set2-%d"
  description = "testing-vnic-set2"
  ip_address_prefix = "10.1.15.0/24"
}

resource "opc_compute_instance" "foo" {
  name = "test-vnic-set-%d"
  label = "testing"
  shape = "oc3"
  image_list = "%s"
  networking_info {
    index = 0
    ip_network = "${opc_compute_ip_network.foo.id}"
    vnic = "test-vnic-set-%d"
    shared_network = false
  }
  networking_info {
    index = 1
    ip_network = "${opc_compute_ip_network.bar.id}"
    vnic = "test-vnic-set2-%d"
    shared_network = false
  }
}

data "opc_compute_network_interface" "foo" {
  instance_name = "${opc_compute_instance.foo.name}"
  instance_id = "${opc_compute_instance.foo.id}"
  interface = "eth0"
}

data "opc_compute_network_interface" "bar" {
  instance_name = "${opc_compute_instance.foo.name}"
  instance_id = "${opc_compute_instance.foo.id}"
  interface = "eth1"
}

resource "opc_compute_vnic_set" "test" {
  name = "%s"
  description = "%s-updated"
  tags = ["tag1"]
  virtual_nics = [
    "${data.opc_compute_network_interface.foo.vnic}",
    "${data.opc_compute_network_interface.bar.vnic}",
  ]
}`, rInt, rInt, rInt, TEST_IMAGE_LIST, rInt, rInt, rName, rDesc)
}

func testAccVnicSetInstance(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_ip_network" "foo" {
  name                = "acctest-ip-network-%d"
  ip_address_prefix   = "192.168.1.0/24"
}

resource "opc_compute_acl" "foo" {
  name        = "acctest-acl-%d"
}

resource "opc_compute_vnic_set" "test" {
  name         = "vnicset-acctest-%d"
  applied_acls = ["${opc_compute_acl.foo.name}"]
}

resource "opc_compute_instance" "foo" {
  name = "acctest-%d"
  hostname = "my-instance"
  label = "my-instance"
  shape = "oc3"
  image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
  networking_info {
    index = 0
    ip_network = "${opc_compute_ip_network.foo.name}"
    ip_address = "192.168.1.100"
    vnic = "my-instance_eth0"
    vnic_sets = [ "${opc_compute_vnic_set.test.name}"]
  }
}`, rInt, rInt, rInt, rInt)
}

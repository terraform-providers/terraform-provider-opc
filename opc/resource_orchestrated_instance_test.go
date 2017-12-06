package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCOrchestratedInstance_Basic(t *testing.T) {
	resName := "opc_compute_orchestrated_instance.test"
	ri := acctest.RandInt()
	config := testAccOrchestrationBasic(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_BasicTwoInstance(t *testing.T) {
	resName := "opc_compute_orchestrated_instance.test"
	ri := acctest.RandInt()
	config := testAccOrchestrationBasicTwoInstance(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
					resource.TestCheckResourceAttrSet(resName, "instance.1.id"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_sharedNetworking(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_compute_orchestrated_instance.test"
	instancePath := "instance.0"
	dataName := "data.opc_compute_network_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestratedInstanceSharedNetworking(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckInstanceExists,
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "id")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "availability_domain")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "domain")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "hostname")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "ip_address")),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "name"), fmt.Sprintf("acc-test-instance-%d", rInt)),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "networking_info.#"), "1"),
					// Default Placement Reqs
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "placement_requirements.#"), "2"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "placement_requirements.0"), "/system/compute/allow_instances"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "placement_requirements.1"), "/system/compute/placement/default"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "platform"), "linux"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "priority"), "/oracle/public/default"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "reverse_dns"), "true"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "state"), "running"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "tags.#"), "2"),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "vcable")),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "virtio"), "false"),

					// Check Data Source to validate networking attributes
					resource.TestCheckResourceAttr(dataName, "shared_network", "true"),
					resource.TestCheckResourceAttr(dataName, "nat.#", "1"),
					resource.TestCheckResourceAttr(dataName, "sec_lists.#", "1"),
					resource.TestCheckResourceAttr(dataName, "name_servers.#", "0"),
					resource.TestCheckResourceAttr(dataName, "vnic_sets.#", "0"),
				),
			},
		},
	})
}

func testAccCheckOrchestrationExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.Orchestrations()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_orchestration" {
			continue
		}

		input := compute.GetOrchestrationInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetOrchestration(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Orchestration %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckOrchestrationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.Orchestrations()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_orchestrated_instance" {
			continue
		}

		input := compute.GetOrchestrationInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetOrchestration(&input); err == nil {
			return fmt.Errorf("Orchestration %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccGetInstanceAttr(instance, attr string) string {
	return fmt.Sprintf("%s.%s", instance, attr)
}

func testAccOrchestrationBasic(rInt int) string {
	return fmt.Sprintf(`
  resource "opc_compute_orchestrated_instance" "test" {
    name        = "test_orchestration-%d"
    desired_state = "active"
		instance {
			name = "acc-test-instance-%d"
			label = "TestAccOPCInstance_basic"
			shape = "oc3"
			image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		}
  }
  `, rInt, rInt)
}

func testAccOrchestrationBasicTwoInstance(rInt int) string {
	return fmt.Sprintf(`
  resource "opc_compute_orchestrated_instance" "test" {
    name        = "test_orchestration-%d"
    desired_state = "active"
		instance {
			name = "acc-test-instance-%d"
			label = "TestAccOPCInstance_basic"
			shape = "oc3"
			image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		}
		instance {
			name = "acc-test-instance-two-%d"
			label = "TestAccOPCInstance_basicTwo"
			shape = "oc3"
			image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		}
  }
  `, rInt, rInt, rInt)
}

func testAccOrchestratedInstanceSharedNetworking(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_orchestrated_instance" "test" {
	name        = "test_orchestration-%d"
	desired_state = "active"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_sharedNetworking"
		shape = "oc3"
		image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		tags = ["tag1", "tag2"]
		networking_info {
			index = 0
			nat = ["ippool:/oracle/public/ippool"]
			shared_network = true
		}
	}
}

data "opc_compute_network_interface" "test" {
  instance_name = "${opc_compute_orchestrated_instance.test.instance.0.name}"
  instance_id = "${opc_compute_orchestrated_instance.test.instance.0.id}"
  interface = "eth0"
}
`, rInt, rInt)
}

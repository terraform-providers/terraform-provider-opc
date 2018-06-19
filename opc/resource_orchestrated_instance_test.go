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

func TestAccOPCOrchestratedInstance_Persistent(t *testing.T) {
	resName := "opc_compute_orchestrated_instance.test"
	ri := acctest.RandInt()
	config := testAccOrchestrationPersistent(ri)
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
					resource.TestCheckResourceAttr(resName, "instance.0.persistent", "true"),
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
					testAccCheckOrchestrationExists,
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

func TestAccOPCOrchestratedInstance_ipNetwork(t *testing.T) {
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
				Config: testAccOrchestratedInstanceIPNetworking(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "id")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "availability_domain")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "domain")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "ip_address")),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "name"), fmt.Sprintf("acc-test-instance-%d", rInt)),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "networking_info.#"), "1"),
					// Default Placement Reqs
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "placement_requirements.#"), "2"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "placement_requirements.0"), "/system/compute/allow_instances"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "platform"), "linux"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "priority"), "/oracle/public/default"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "reverse_dns"), "true"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "state"), "running"),
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "virtio"), "false"),

					// Check Data Source to validate networking attributes
					resource.TestCheckResourceAttr(dataName, "ip_network", fmt.Sprintf("testing-ip-network-%d", rInt)),
					resource.TestCheckResourceAttr(dataName, "vnic", fmt.Sprintf("ip-network-test-%d", rInt)),
					resource.TestCheckResourceAttr(dataName, "shared_network", "false"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_ipNetworkIsDefaultGateway(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_compute_orchestrated_instance.test"
	instancePath := "instance.0"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestratedInstanceIPNetworkingDefaultGateway(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "id")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "availability_domain")),
					resource.TestCheckResourceAttrSet(resName, testAccGetInstanceAttr(instancePath, "domain")),
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
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "virtio"), "false"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_storage(t *testing.T) {
	resName := "opc_compute_orchestrated_instance.test"
	instancePath := "instance.0"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestratedInstanceStorage(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttr(resName, testAccGetInstanceAttr(instancePath, "storage.#"), "2"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_noBoot(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccOrchestrationBasic_noBoot(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("One of boot_order or image_list must be set for instance to be created"),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_inactive(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccOrchestrationInactive(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_activeToInactive(t *testing.T) {
	ri := acctest.RandInt()
	resName := "opc_compute_orchestrated_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationBasic(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
				),
			},
			{
				Config: testAccOrchestrationInactive(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_inactiveToActive(t *testing.T) {
	ri := acctest.RandInt()
	resName := "opc_compute_orchestrated_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationInactive(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
				),
			},
			{
				Config: testAccOrchestrationBasic(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_activeToSuspend(t *testing.T) {
	ri := acctest.RandInt()
	resName := "opc_compute_orchestrated_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationBasic(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
					resource.TestCheckResourceAttrSet(resName, "instance.0.id"),
				),
			},
			{
				Config: testAccOrchestrationSuspend(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationExists,
				),
			},
		},
	})
}

func TestAccOPCOrchestratedInstance_105(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrchestration_105(ri),
				ExpectError: regexp.MustCompile("Error creating Orchestration"),
			},
		},
	})
}

func testAccCheckOrchestrationExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).computeClient.Orchestrations()

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
	client := testAccProvider.Meta().(*Client).computeClient.Orchestrations()

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

func testAccOrchestration_105(rInt int) string {
	return fmt.Sprintf(`
		resource "opc_compute_ip_network" "orchestration-test" {
		  name              = "test_orchestration_%d"
		  ip_address_prefix = "192.168.1.0/24"
		}

		resource "opc_compute_orchestrated_instance" "orchestrated-instance-test" {
		  name          = "test_orchestration_%d"
		  description   = "test_orchestration_%d"
		  desired_state = "active"

		  instance {
		    name       = "orchestration-test-%d"
		    hostname   = "orchestration-test-%d"
		    label      = "orchestration-test-%d"
		    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		    shape      = "oc3"

		    networking_info {
		      index          = 0
		      ip_network     = "${opc_compute_ip_network.orchestration-test.name}"
		      vnic           = "orchestrationtest_eth0"
		      ip_address      = "192.168.0.101"
		    }
		  }
		}
  `, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccOrchestrationSuspend(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_orchestrated_instance" "test" {
  name        = "test_orchestration-%d"
  desired_state = "suspend"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_basic"
		shape = "oc3"
		image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
	}
}
  `, rInt, rInt)
}

func testAccOrchestrationPersistent(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_orchestrated_instance" "test" {
  name        = "test_orchestration-%d"
  desired_state = "active"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_basic"
		shape = "oc3"
		image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		persistent = true
	}
}
  `, rInt, rInt)
}

func testAccOrchestrationInactive(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_orchestrated_instance" "test" {
  name        = "test_orchestration-%d"
  desired_state = "inactive"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_basic"
		shape = "oc3"
		image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
	}
}
  `, rInt, rInt)
}

func testAccOrchestrationBasic_noBoot(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_orchestrated_instance" "test" {
  name        = "test_orchestration-%d"
  desired_state = "active"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_basic"
		shape = "oc3"
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

func testAccOrchestratedInstanceIPNetworking(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_ip_network" "foo" {
  name = "testing-ip-network-%d"
  description = "testing-ip-network-instance"
  ip_address_prefix = "10.1.12.0/24"
}

resource "opc_compute_orchestrated_instance" "test" {
	name        = "test_orchestration-%d"
	desired_state = "active"
	instance {
	  name = "acc-test-instance-%d"
	  label = "TestAccOPCInstance_ipNetwork"
	  shape = "oc3"
	  image_list = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
	  networking_info {
	    index = 0
	    ip_network = "${opc_compute_ip_network.foo.id}"
	    vnic = "ip-network-test-%d"
	    shared_network = false
	  }
	}
}

data "opc_compute_network_interface" "test" {
	instance_name = "${opc_compute_orchestrated_instance.test.instance.0.name}"
	instance_id = "${opc_compute_orchestrated_instance.test.instance.0.id}"
  interface = "eth0"
}
`, rInt, rInt, rInt, rInt)
}

func testAccOrchestratedInstanceIPNetworkingDefaultGateway(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_ip_network" "foo" {
  name = "testing-ip-network-%d"
  description = "testing-ip-network-instance"
  ip_address_prefix = "10.1.12.0/24"
}

resource "opc_compute_orchestrated_instance" "test" {
	name        = "test_orchestration-%d"
	desired_state = "active"
	instance {
	  name = "acc-test-instance-%d"
	  label = "TestAccOPCInstance_ipNetwork"
	  shape = "oc3"
	  image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
	  networking_info {
	    index = 0
	    ip_network = "${opc_compute_ip_network.foo.id}"
	    vnic = "ip-network-test-%d"
	    shared_network = false
			is_default_gateway = true
	  }
	}
}
`, rInt, rInt, rInt, rInt)
}

func testAccOrchestratedInstanceStorage(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_storage_volume" "foo" {
  name = "acc-test-orchestration-%d"
  size = 1
}

resource "opc_compute_storage_volume" "bar" {
  name = "acc-test-orchestration-2-%d"
  size = 1
}

resource "opc_compute_orchestrated_instance" "test" {
	name        = "test_orchestration-%d"
	desired_state = "active"
	instance {
		name = "acc-test-instance-%d"
		label = "TestAccOPCInstance_basic"
		shape = "oc3"
		image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
		storage {
			volume = "${opc_compute_storage_volume.foo.name}"
			index = 1
		}
		storage {
		  volume = "${opc_compute_storage_volume.bar.name}"
		  index = 2
		}
	}
}`, rInt, rInt, rInt, rInt)
}

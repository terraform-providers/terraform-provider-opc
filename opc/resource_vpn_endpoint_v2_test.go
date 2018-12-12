package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCVPNEndpointV2_Basic(t *testing.T) {
	resourceName := "opc_compute_vpn_endpoint_v2.test"
	ri := acctest.RandInt()
	config := testAccVPNEndpointV2Basic(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNEndpointV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEndpointV2Exists,
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pre_shared_key", "customer_vpn_gateway"},
			},
		},
	})
}

func TestAccOPCVPNEndpointV2_Update(t *testing.T) {
	resourceName := "opc_compute_vpn_endpoint_v2.test"
	ri := acctest.RandInt()
	config := testAccVPNEndpointV2Basic(ri)
	config2 := testAccVPNEndpointV2Update(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNEndpointV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEndpointV2Exists,
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEndpointV2Exists,
					resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "fdsafdsa"),
				),
			},
		},
	})
}

func testAccCheckVPNEndpointV2Exists(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).computeClient.VPNEndpointV2s()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_vpn_endpoint_v2" {
			continue
		}

		input := compute.GetVPNEndpointV2Input{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetVPNEndpointV2(&input); err != nil {
			return fmt.Errorf("Error retrieving state of VPNEndpointV2 %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckVPNEndpointV2Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).computeClient.VPNEndpointV2s()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_vpn_endpoint_v2" {
			continue
		}

		input := compute.GetVPNEndpointV2Input{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetVPNEndpointV2(&input); err == nil {
			return fmt.Errorf("VPNEndpointV2 %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccVPNEndpointV2Basic(rInt int) string {
	return fmt.Sprintf(`
	resource "opc_compute_ip_network" "test" {
		name = "testing-ip-network-%d"
		ip_address_prefix = "10.0.12.0/24"
	}

	resource "opc_compute_vpn_endpoint_v2" "test" {
	  name        = "test_vpn_endpoint_v2-%d"
	  customer_vpn_gateway = "127.0.0.1"
	  ip_network = "${opc_compute_ip_network.test.name}"
	  pre_shared_key = "asdfasdf"
	  reachable_routes = ["127.0.0.1/24"]
	  vnic_sets = ["default"]
	}
	`, rInt, rInt)
}

func testAccVPNEndpointV2Update(rInt int) string {
	return fmt.Sprintf(`
	resource "opc_compute_ip_network" "test" {
		name = "testing-ip-network-%d"
		ip_address_prefix = "10.0.12.0/24"
	}

	resource "opc_compute_vpn_endpoint_v2" "test" {
	  name        = "test_vpn_endpoint_v2-%d"
	  customer_vpn_gateway = "127.0.0.1"
	  ip_network = "${opc_compute_ip_network.test.name}"
	  pre_shared_key = "fdsafdsa"
	  reachable_routes = ["127.0.0.1/24"]
	  vnic_sets = ["default"]
	}
	`, rInt, rInt)
}

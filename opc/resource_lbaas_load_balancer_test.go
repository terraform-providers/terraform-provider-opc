package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/lbaas"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCLoadBalancer_Basic(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccLoadBalancerBasic, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists,
				),
			},
		},
	})
}

func testAccCheckLoadBalancerExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).lbaasClient.LoadBalancerClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_lbaas_load_balancer" {
			continue
		}

		lb := lbaas.LoadBalancerContext{
			Region: rs.Primary.Attributes["region"],
			Name:   rs.Primary.Attributes["name"],
		}

		if _, err := client.GetLoadBalancer(lb); err != nil {
			return fmt.Errorf("Error retrieving state of Load Balancer %s: %s", lb.Name, err)
		}
	}

	return nil
}

func testAccCheckLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).lbaasClient.LoadBalancerClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_lbaas_load_balancer" {
			continue
		}

		lb := lbaas.LoadBalancerContext{
			Region: rs.Primary.Attributes["region"],
			Name:   rs.Primary.Attributes["name"],
		}

		if info, err := client.GetLoadBalancer(lb); err == nil {
			return fmt.Errorf("Load Balancer %s still exists: %#v", lb.Name, info)
		}
	}

	return nil
}

var testAccLoadBalancerBasic = `
resource "opc_lbaas_load_balancer" "test" {
	region      = "uscom-central-1"
  name        = "acctest%d"
	scheme      = "INTERNET_FACING"
}
`

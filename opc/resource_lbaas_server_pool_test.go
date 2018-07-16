package opc

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBaaSServerPool_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_server_pool.test"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckServerPoolDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSServerPoolConfig_Basic(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckServerPoolExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "vnic_set", ""),
					resource.TestCheckResourceAttr(resName, "servers.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "health_check.0.type", "http"),
					resource.TestCheckResourceAttr(resName, "health_check.0.interval", "30"),
					resource.TestCheckResourceAttr(resName, "health_check.0.timeout", "20"),
					resource.TestCheckResourceAttr(resName, "health_check.0.healthy_threshold", "10"),
					resource.TestCheckResourceAttr(resName, "health_check.0.unhealthy_threshold", "5"),
					resource.TestCheckResourceAttr(resName, "health_check.0.accepted_return_codes.#", "2"),
				),
			},
		},
	})
}

func testAccLBaaSServerPoolConfig_Basic(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_server_pool" "test" {
  load_balancer = "%s"

  name    = "acctest-%d"
  servers = ["192.168.1.100:8080","192.168.1.101:8080"]
  tags = ["TESTING", "Terraform"]

  health_check {
    type = "http"
    interval = 30
    timeout = 20
    healthy_threshold = 10
    unhealthy_threshold = 5
    accepted_return_codes = [ "2xx", "3xx" ]
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSCheckServerPoolExists(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.OriginServerPoolClient()

	if _, err := client.GetOriginServerPool(lb, name); err != nil {
		return fmt.Errorf("Error retrieving state of Server Pool '%s': %v", name, err)
	}

	return nil
}

func testAccLBaaSCheckServerPoolDestroyed(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.OriginServerPoolClient()

	if info, _ := client.GetOriginServerPool(lb, name); info != nil {
		return fmt.Errorf("Server Pool '%s' still exists: %+v", name, info)
	}
	return nil
}

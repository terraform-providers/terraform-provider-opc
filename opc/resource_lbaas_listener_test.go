package opc

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBaaSListener_Basic(t *testing.T) {
	if checkSkipLBTests() {
		t.Skip(fmt.Printf("`OPC_LBAAS_ENDPOINT` not set, skipping test"))
	}

	rInt := acctest.RandInt()
	resName := "opc_lbaas_listener.test"
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
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckListenerDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSListenerConfig_Basic(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckListenerExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "balancer_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "virtual_hosts.#", "1"),
					resource.TestCheckResourceAttr(resName, "path_prefixes.#", "3"),
				),
			},
		},
	})
}

func TestAccLBaaSListener_BasicUpdate(t *testing.T) {
	if checkSkipLBTests() {
		t.Skip(fmt.Printf("`OPC_LBAAS_ENDPOINT` not set, skipping test"))
	}

	rInt := acctest.RandInt()
	resName := "opc_lbaas_listener.test"
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
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckListenerDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSListenerConfig_Basic(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckListenerExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "balancer_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "virtual_hosts.#", "1"),
					resource.TestCheckResourceAttr(resName, "path_prefixes.#", "3"),
				),
			},
			{
				Config: testAccLBaaSListenerConfig_BasicUpdate(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckListenerExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "balancer_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resName, "virtual_hosts.#", "2"),
					resource.TestCheckResourceAttr(resName, "path_prefixes.#", "2"),
				),
			},
		},
	})
}

func testAccLBaaSListenerConfig_Basic(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_listener" "test" {
	load_balancer = "%s"
  name          = "acctest-%d"

  port              = 8080
  balancer_protocol = "HTTP"
  server_protocol   = "HTTP"
  virtual_hosts     = ["mywebapp.example.com"]
	path_prefixes     = ["one","two","three"]

  tags = [ "TESTING", "Terraform"]
}

%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSListenerConfig_BasicUpdate(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_listener" "test" {
	load_balancer = "%s"
  name          = "acctest-%d"

  port              = 8080
  balancer_protocol = "HTTP"
  server_protocol   = "HTTP"
  virtual_hosts     = ["mywebapp.example.com","mywebapp2.example.com"]
	path_prefixes     = ["one","two"]

  tags = [ "TESTING", "Terraform", "extra"]
}

%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSCheckListenerExists(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.ListenerClient()

	if _, err := client.GetListener(lb, name); err != nil {
		return fmt.Errorf("Error retrieving state of Listener '%s': %v", name, err)
	}

	return nil
}

func testAccLBaaSCheckListenerDestroyed(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.ListenerClient()

	if info, _ := client.GetListener(lb, name); info != nil {
		return fmt.Errorf("Listener '%s' still exists: %+v", name, info)
	}
	return nil
}

func testAccParentLoadBalancerConfig(lbCount, rInt int) string {
	if lbCount == 0 {
		return ""
	}

	return fmt.Sprintf(`
resource "opc_lbaas_load_balancer" "test" {
	count = %d
	name = "acctest-%d"
	region = "uscom-central-1"
	scheme = "INTERNET_FACING"
}
`, lbCount, rInt)
}

package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCSecurityList_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := "opc_compute_security_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSecurityListBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityListExists,
					resource.TestCheckResourceAttr(rName, "policy", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "outbound_cidr_policy", "DENY"),
				),
			},
		},
	})
}

func TestAccOPCSecurityList_complete(t *testing.T) {
	rInt := acctest.RandInt()
	rName := "opc_compute_security_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSecurityListComplete(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityListExists,
					resource.TestCheckResourceAttr(rName, "policy", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "outbound_cidr_policy", "DENY"),
					resource.TestCheckResourceAttr(rName, "description", "Acceptance Test Security List Complete"),
				),
			},
		},
	})
}

func TestAccOPCSecurityList_lowercasePolicies(t *testing.T) {
	rInt := acctest.RandInt()
	rName := "opc_compute_security_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSecurityListLowercasePolicies(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityListExists,
					resource.TestCheckResourceAttr(rName, "policy", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "outbound_cidr_policy", "DENY"),
					resource.TestCheckResourceAttr(rName, "description", "Acceptance Test Security List Lowercase"),
				),
			},
		},
	})
}

func testAccCheckSecurityListExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.SecurityLists()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_security_list" {
			continue
		}

		input := compute.GetSecurityListInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetSecurityList(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Security List %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccCheckSecurityListDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).computeClient.SecurityLists()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_security_list" {
			continue
		}

		input := compute.GetSecurityListInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetSecurityList(&input); err == nil {
			return fmt.Errorf("Security List %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccOPCSecurityListBasic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_security_list" "test" {
	name                 = "acc-test-sec-list-%d"
	policy               = "PERMIT"
	outbound_cidr_policy = "DENY"
}`, rInt)
}

func testAccOPCSecurityListComplete(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_security_list" "test" {
 name                 = "acc-test-sec-list-%d"
 description          = "Acceptance Test Security List Complete"
 policy               = "PERMIT"
 outbound_cidr_policy = "DENY"
}`, rInt)
}

func testAccOPCSecurityListLowercasePolicies(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_security_list" "test" {
 name                 = "acc-test-sec-list-%d"
 description          = "Acceptance Test Security List Lowercase"
 policy               = "permit"
 outbound_cidr_policy = "deny"
}`, rInt)
}

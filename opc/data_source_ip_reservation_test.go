package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDataSourceIPReservation_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIPReservationBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opc_compute_ip_reservation.bar", "parent_pool", "/oracle/public/ippool"),
				),
			},
		},
	})
}

func testAccDataSourceIPReservationBasic(rInt int) string {

	config := `
  resource "opc_compute_ip_reservation" "foo" {
		name        = "acc-test-ip-reservation-%d"
	  parent_pool = "/oracle/public/ippool"
	  permanent   = true
	}

	data "opc_compute_ip_reservation" "bar" {
		name = "${opc_compute_ip_reservation.foo.name}"
	}
	`

	return fmt.Sprintf(config, rInt)
}

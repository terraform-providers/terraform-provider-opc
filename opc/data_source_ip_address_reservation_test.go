package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDataSourceIPAddressReservation_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIPAddressReservationBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opc_compute_ip_address_reservation.bar", "ip_address_pool", "public-ippool"),
				),
			},
		},
	})
}

func testAccDataSourceIPAddressReservationBasic(rInt int) string {

	config := `
  resource "opc_compute_ip_address_reservation" "foo" {
		name        = "acc-test-ip-address-reservation-%d"
		description = "terraform acceptance test"
	  ip_address_pool = "public-ippool"
	}

	data "opc_compute_ip_address_reservation" "bar" {
		name = "${opc_compute_ip_address_reservation.foo.name}"
	}
	`

	return fmt.Sprintf(config, rInt)
}

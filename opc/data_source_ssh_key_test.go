package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

const test_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ=="

func TestAccOPCDataSourceSSHKey_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSSHKeyBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opc_compute_ssh_key.bar", "key", test_ssh_key),
				),
			},
		},
	})
}

func testAccDataSourceSSHKeyBasic(rInt int) string {

	testAccMachineImageBasic := `
        resource "opc_compute_ssh_key" "foo" {
		name    = "acc-test-ssh%d"
		key     = "%s"
		enabled = true
	}

	data "opc_compute_ssh_key" "bar" {
		name = "${opc_compute_ssh_key.foo.name}"
	}`

	return fmt.Sprintf(testAccMachineImageBasic, rInt, test_ssh_key)
}

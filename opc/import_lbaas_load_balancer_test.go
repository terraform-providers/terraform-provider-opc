package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCLoadBalancer_importBasic(t *testing.T) {
	if checkSkipLBTests() {
		t.Skip(fmt.Printf("`OPC_LBAAS_ENDPOINT` not set, skipping test"))
	}

	resourceName := "opc_lbaas_load_balancer.test"

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccLoadBalancerBasic, ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

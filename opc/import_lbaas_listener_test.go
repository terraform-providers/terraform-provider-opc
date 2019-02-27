package opc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBaaSListener_importBasic(t *testing.T) {
	if checkSkipLBTests() {
		t.Skip(fmt.Printf("`OPC_LBAAS_ENDPOINT` not set, skipping test"))
	}

	rInt := acctest.RandInt()
	resName := "opc_lbaas_listener.test"

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	config := testAccLBaaSListenerConfig_Basic(lbID, rInt, lbCount)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckListenerDestroyed),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

package opc

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCDatabaseServiceInstance_importBasic(t *testing.T) {
	t.Skip("Skipping test until we release this resource")

	resourceName := "opc_database_service_instance.test"

	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceBasic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter", "vm_public_key"},
			},
		},
	})
}

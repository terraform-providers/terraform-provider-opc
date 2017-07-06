package opc

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCSecurityList_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	rName := "opc_compute_security_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSecurityListBasic(rInt),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOPCSecurityList_importComplete(t *testing.T) {
	rInt := acctest.RandInt()
	rName := "opc_compute_security_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCSecurityListComplete(rInt),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

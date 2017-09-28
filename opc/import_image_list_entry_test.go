package opc

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCImageListEntry_importBasic(t *testing.T) {
	rInt := acctest.RandInt()

	rName := "opc_compute_image_list_entry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImageListEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImageListEntry_basic(rInt),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

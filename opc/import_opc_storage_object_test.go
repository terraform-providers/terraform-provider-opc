package opc

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPCStorageObject_importBasic(t *testing.T) {
	rName := "opc_storage_object.test"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCStorageObject_contentSource(rInt, _SourceInput),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"transaction_id", "content"},
			},
		},
	})
}

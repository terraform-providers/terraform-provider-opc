package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCJavaServiceInstance_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccJavaServiceInstanceBasic(ri)
	resourceName := "opc_java_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJavaServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "shape", "oc1m"),
					resource.TestCheckResourceAttr(
						resourceName, "level", "PAAS"),
				),
			},
		},
	})
}

func testAccCheckJavaServiceInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).javaClient.ServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_java_service_instance" {
			continue
		}

		input := java.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetServiceInstance(&input); err != nil {
			return fmt.Errorf("Error retrieving state of JavaServiceInstance %s: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckJavaServiceInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).javaClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Java Client is not initialized. Make sure to use `java_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_java_service_instance" {
			continue
		}

		input := java.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetServiceInstance(&input); err == nil {
			return fmt.Errorf("JavaServiceInstance %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccJavaServiceInstanceBasic(rInt int) string {
	return fmt.Sprintf(`resource "opc_database_service_instance" "test" {
    name        = "test-service-instance-%d"
    description = "test service instance"
    edition = "EE_EP"
    level = "PAAS"
    shape = "oc1m"
    subscription_type = "HOURLY"
    version = "12.2.0.1"
    vm_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    parameter {
      admin_password = "Test_String7"
      backup_destination = "OSS"
      failover_database = false
      sid = "ORCL"
      usable_storage = 15
    }
    cloud_storage {
      container = "Storage-canonical/matthew-test-%d"
      create_if_missing = true
    }
	}

  resource "opc_java_service_instance" "test" {
    name = "test-java-service-instance-%d"
    type = "weblogic"
    shape = "oc1m"
    version = "12.2.1"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    subscription_type = "HOURLY"
    cloud_storage {
      container = "Storage-canonical/test-terraform-java-instance"
      create_if_missing = true
    }
    database {
      name = "${opc_database_service_instance.test.name}"
      username = "sys"
      password = "Test_String7"
    }
    admin {
      username = "terraform-user"
      password = "Test_String7"
    }
  }

  `, rInt, rInt, rInt)
}

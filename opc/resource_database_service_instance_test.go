package opc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPCDatabaseServiceInstance_Basic(t *testing.T) {
	t.Skip("Skipping test until we release this resource")

	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceBasic(ri)
	resourceName := "opc_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test service instance"),
					resource.TestCheckResourceAttr(
						resourceName, "edition", "EE_EP"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "12.2.0.1"),
					resource.TestCheckResourceAttr(
						resourceName, "parameter.#", "1"),
				),
			},
		},
	})
}

func TestAccOPCDatabaseServiceInstance_CloudStorage(t *testing.T) {
	t.Skip("Skipping test until we release this resource")

	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceCloudStorage(ri)
	resourceName := "opc_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "cloud_storage_container", fmt.Sprintf("Storage-%s/acctest-%d", os.Getenv("OPC_IDENTITY_DOMAIN"), ri)),
				),
			},
		},
	})
}

func testAccCheckDatabaseServiceInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).databaseClient.ServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_database_service_instance" {
			continue
		}

		input := database.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetServiceInstance(&input); err != nil {
			return fmt.Errorf("Error retrieving state of DatabaseServiceInstance %s: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckDatabaseServiceInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPCClient).databaseClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_database_service_instance" {
			continue
		}

		input := database.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetServiceInstance(&input); err == nil {
			return fmt.Errorf("DatabaseServiceInstance %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccDatabaseServiceInstanceBasic(rInt int) string {
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
      backup_destination = "NONE"
      sid = "ORCL"
      usable_storage = 15
    }
  }`, rInt)
}

func testAccDatabaseServiceInstanceCloudStorage(rInt int) string {
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
      container = "Storage-%s/acctest-%d"
      create_if_missing = true
    }
	}`, rInt, os.Getenv("OPC_IDENTITY_DOMAIN"), rInt)
}

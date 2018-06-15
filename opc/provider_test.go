package opc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"opc": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Error creating Provider: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	required := []string{"OPC_USERNAME", "OPC_PASSWORD", "OPC_IDENTITY_DOMAIN", "OPC_ENDPOINT", "OPC_STORAGE_ENDPOINT"}
	for _, prop := range required {
		if os.Getenv(prop) == "" {
			t.Fatalf("%s must be set for acceptance test", prop)
		}
	}
	config := Config{
		User:             os.Getenv("OPC_USERNAME"),
		Password:         os.Getenv("OPC_PASSWORD"),
		IdentityDomain:   os.Getenv("OPC_IDENTITY_DOMAIN"),
		Endpoint:         os.Getenv("OPC_ENDPOINT"),
		MaxRetries:       1,
		Insecure:         false,
		StorageEndpoint:  os.Getenv("OPC_STORAGE_ENDPOINT"),
		StorageServiceId: os.Getenv("OPC_STORAGE_SERVICE_ID"),
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}
	if client.storageClient == nil {
		t.Fatalf("Storage Client is nil. Make sure your Oracle Cloud Account has access to the Storage Cloud")
	}
}

type OPCResourceState struct {
	*compute.ComputeClient
	*terraform.InstanceState
}

func opcResourceCheck(resourceName string, f func(checker *OPCResourceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		state := &OPCResourceState{
			ComputeClient: testAccProvider.Meta().(*OPCClient).computeClient,
			InstanceState: rs.Primary,
		}

		return f(state)
	}
}

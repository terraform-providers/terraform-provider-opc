package lbaas

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/terraform/helper/acctest"
)

// GetTestClient obtains a client for testing purposes
func GetTestClient(c *opc.Config) (*Client, error) {
	// Build up config with default values if omitted

	if c.Username == nil {
		username := os.Getenv("OPC_USERNAME")
		c.Username = &username
	}

	if c.Password == nil {
		password := os.Getenv("OPC_PASSWORD")
		c.Password = &password
	}

	if c.APIEndpoint == nil {
		apiEndpoint, err := url.Parse(os.Getenv("OPC_LBAAS_ENDPOINT"))
		if err != nil {
			return nil, err
		}
		c.APIEndpoint = apiEndpoint
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				TLSHandshakeTimeout: 120 * time.Second},
		}
	}

	return NewClient(c)
}

func getLoadBalancerClient() (*LoadBalancerClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &LoadBalancerClient{}, err
	}
	return client.LoadBalancerClient(), nil
}

func destroyLoadBalancer(t *testing.T, client *LoadBalancerClient, lb LoadBalancerContext) {
	if _, err := client.DeleteLoadBalancer(lb); err != nil {
		t.Fatal(err)
	}
}

// utility function to create a load balancer instance needed for testing child resources
func createParentLoadBalancer(t *testing.T) (LoadBalancerContext, *LoadBalancerClient) {

	// if environment variable `OPC_TEST_USE_EXISTING_LB` is set an existing LB instance
	// can be using instead of waiting for a new one to be created.
	if existing := os.Getenv("OPC_TEST_USE_EXISTING_LB"); existing != "" {
		// expecting LB instance id in the format `region/name`
		s := strings.Split(existing, "/")
		lb := LoadBalancerContext{
			Region: s[0],
			Name:   s[1],
		}
		return lb, nil
	}

	// create a new Load Balancer instance

	rInt := acctest.RandInt()
	name := fmt.Sprintf("acctestlb-%d", rInt)

	var region string
	if region = os.Getenv("OPC_TEST_LBAAS_REGION"); region == "" {
		region = "uscom-central-1"
	}

	lbClient, err := getLoadBalancerClient()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        name,
		Region:      region,
		Description: "Terraform Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    LBaaSDisabledTrue,
	}

	_, err = lbClient.CreateLoadBalancer(createLoadBalancerInput)
	if err != nil {
		t.Fatal(err)
	}

	lb := LoadBalancerContext{
		Region: createLoadBalancerInput.Region,
		Name:   createLoadBalancerInput.Name,
	}

	return lb, lbClient
}

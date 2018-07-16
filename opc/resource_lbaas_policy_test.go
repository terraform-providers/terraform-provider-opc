package opc

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBaaSPolicy_ApplicationCookieStickinessPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.application_cookie_stickiness_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_ApplicationCookieStickinessPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.0.cookie_name", "MY_APP_COOKIE"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_CloudGatePolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.cloudgate_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_CloudGatePolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "cloudgate_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "cloudgate_policy.0.cloudgate_application", "example-cloudgate-app"),
					resource.TestCheckResourceAttr(resName, "cloudgate_policy.0.cloudgate_policy_name", "example-cloudgate-policy"),
					resource.TestCheckResourceAttr(resName, "cloudgate_policy.0.identity_service_instance_guid", "blahblahblah"),
					resource.TestCheckResourceAttr(resName, "cloudgate_policy.0.virtual_hostname_for_policy_attribution", "host1.example.com"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_LoadBalancerCookieStickinessPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.load_balancer_cookie_stickiness_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_LoadBalancerCookieStickinessPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "load_balancer_cookie_stickiness_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "load_balancer_cookie_stickiness_policy.0.cookie_expiration_period", "60"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_LoadBalancingMechanismPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.load_balancing_mechanism_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_LoadBalancingMechanismPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "load_balancing_mechanism_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "load_balancing_mechanism_policy.0.load_balancing_mechanism", "round_robin"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_RateLimitingRequestPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.rate_limiting_request_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_RateRequestLimitingPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.burst_size", "10"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.delay_excessive_requests", "true"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.http_error_code", "503"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.logging_level", "notice"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.rate_limiting_criteria", "server"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.requests_per_second", "1"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.zone_memory_size", "10"),
					resource.TestCheckResourceAttr(resName, "rate_limiting_request_policy.0.zone", "examplezone"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_RedirectPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.redirect_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_RedirectPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "redirect_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "redirect_policy.0.redirect_uri", "https://redirect.example.com"),
					resource.TestCheckResourceAttr(resName, "redirect_policy.0.response_code", "306"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_ResourceAccessControlPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.resource_access_control_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_ResourceAccessControlPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "resource_access_control_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "resource_access_control_policy.0.disposition", "DENY_ALL"),
					resource.TestCheckResourceAttr(resName, "resource_access_control_policy.0.denied_clients.#", "2"),
					resource.TestCheckResourceAttr(resName, "resource_access_control_policy.0.permitted_clients.#", "1"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_SetRequestHeaderPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.set_request_header_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_SetRequestHeaderPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.0.header_name", "X-Custom-Header"),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.0.value", "foo-bar"),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.0.action_when_header_exists", "OVERWRITE"),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.0.action_when_header_value_is.#", "3"),
					resource.TestCheckResourceAttr(resName, "set_request_header_policy.0.action_when_header_value_is_not.#", "3"),
				),
			},
		},
	})
}

func TestAccLBaaSPolicy_SSLNegotiationPolicy(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.ssl_negotiation_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_SSLNegotiationPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "ssl_negotiation_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "ssl_negotiation_policy.0.port", "8022"),
					resource.TestCheckResourceAttr(resName, "ssl_negotiation_policy.0.server_order_preference", "ENABLED"),
					resource.TestCheckResourceAttr(resName, "ssl_negotiation_policy.0.ssl_protocol.#", "3"),
					resource.TestCheckResourceAttr(resName, "ssl_negotiation_policy.0.ssl_ciphers.#", "1"),
				),
			},
		},
	})
}

// TODO TrustedCertificatePolicy Test disabled due to issue deleting Trusted certificates.
// func TestAccLBaaSPolicy_TrustedCertificatePolicy(t *testing.T) {
// 	rInt := acctest.RandInt()
// 	resName := "opc_lbaas_policy.trusted_certificate_policy"
// 	testName := fmt.Sprintf("acctest-%d", rInt)
//
// 	// use existing LB instance from environment if set
// 	lbCount := 0
// 	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
// 	if lbID == "" {
// 		lbCount = 1
// 		lbID = "${opc_lbaas_load_balancer.test.id}"
// 	}
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccLBaaSPolicyConfig_TrustedCertificatePolicy(lbID, rInt, lbCount),
// 				Check: resource.ComposeTestCheckFunc(
// 					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
// 					resource.TestCheckResourceAttr(resName, "name", testName),
// 					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
// 					resource.TestCheckResourceAttr(resName, "trusted_certificate_policy.#", "1"),
// 					resource.TestCheckResourceAttrSet(resName, "trusted_certificate_policy.0.trusted_certificate"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccLBaaSPolicy_ApplicationCookieStickinessPolicy_Update(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.application_cookie_stickiness_policy"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_ApplicationCookieStickinessPolicy(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.0.cookie_name", "MY_APP_COOKIE"),
				),
			},
			{
				Config: testAccLBaaSPolicyConfig_ApplicationCookieStickinessPolicyUpdate(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.0.cookie_name", "MY_APP_COOKIE_UPDATED"),
				),
			},
		},
	})
}

// Test to check that policy type can be completely changed successfully
func TestAccLBaaSPolicy_Update_ChangePolicyType(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_lbaas_policy.update_policy_type"
	testName := fmt.Sprintf("acctest-%d", rInt)

	// use existing LB instance from environment if set
	lbCount := 0
	lbID := os.Getenv("OPC_TEST_USE_EXISTING_LB")
	if lbID == "" {
		lbCount = 1
		lbID = "${opc_lbaas_load_balancer.test.id}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: opcResourceCheck(resName, testAccLBaaSCheckPolicyDestroyed),
		Steps: []resource.TestStep{
			{
				Config: testAccLBaaSPolicyConfig_UpdatePolicyType_Create(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "application_cookie_stickiness_policy.0.cookie_name", "MY_APP_COOKIE"),
				),
			},
			{
				Config: testAccLBaaSPolicyConfig_UpdatePolicyType_Update(lbID, rInt, lbCount),
				Check: resource.ComposeTestCheckFunc(
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					opcResourceCheck(resName, testAccLBaaSCheckPolicyExists),
					resource.TestCheckResourceAttr(resName, "name", testName),
					resource.TestMatchResourceAttr(resName, "uri", regexp.MustCompile(testName)),
					resource.TestCheckResourceAttr(resName, "redirect_policy.#", "1"),
					resource.TestCheckResourceAttr(resName, "redirect_policy.0.redirect_uri", "https://redirect.example.com"),
					resource.TestCheckResourceAttr(resName, "redirect_policy.0.response_code", "306"),
				),
			},
		},
	})
}

func testAccLBaaSPolicyConfig_ApplicationCookieStickinessPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "application_cookie_stickiness_policy" {
  load_balancer = "%s"
  name          = "acctest-%d"

  application_cookie_stickiness_policy {
    cookie_name = "MY_APP_COOKIE"
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_ApplicationCookieStickinessPolicyUpdate(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "application_cookie_stickiness_policy" {
  load_balancer = "%s"
  name          = "acctest-%d"

  application_cookie_stickiness_policy {
    cookie_name = "MY_APP_COOKIE_UPDATED"
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_CloudGatePolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "cloudgate_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  cloudgate_policy {
    cloudgate_application = "example-cloudgate-app"
    cloudgate_policy_name = "example-cloudgate-policy"
    identity_service_instance_guid = "blahblahblah"
    virtual_hostname_for_policy_attribution = "host1.example.com"
  }
}
`, lbID, rInt)
}

func testAccLBaaSPolicyConfig_LoadBalancerCookieStickinessPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "load_balancer_cookie_stickiness_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  load_balancer_cookie_stickiness_policy {
    cookie_expiration_period = 60
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_LoadBalancingMechanismPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "load_balancing_mechanism_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  load_balancing_mechanism_policy {
    load_balancing_mechanism = "round_robin"
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_RateRequestLimitingPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "rate_limiting_request_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  rate_limiting_request_policy {
    burst_size = 10
    delay_excessive_requests = true
    http_error_code = 503
    logging_level = "notice"
    rate_limiting_criteria = "server"
    requests_per_second = 1
    zone_memory_size = 10
    zone = "examplezone"
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_RedirectPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "redirect_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  redirect_policy {
    redirect_uri = "https://redirect.example.com"
    response_code = 306
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_ResourceAccessControlPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "resource_access_control_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  resource_access_control_policy {
    disposition = "DENY_ALL"
    denied_clients = ["192.168.0.1/24", "192.168.0.2/24"]
    permitted_clients = ["10.0.0.0/16"]
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_SetRequestHeaderPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "set_request_header_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  set_request_header_policy {
    header_name                     = "X-Custom-Header"
    value                           = "foo-bar"
    action_when_header_exists       = "OVERWRITE"
    action_when_header_value_is     = ["bar", "foo", "adc"]
    action_when_header_value_is_not = ["ABC", "ZZZ", "XYZ"]
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_SSLNegotiationPolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "ssl_negotiation_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  ssl_negotiation_policy {
    port = 8022
    server_order_preference = "ENABLED"
    ssl_protocol = ["SSLv3", "TLSv1.1", "TLSv1.2"]
    ssl_ciphers = ["AES256-SHA"]
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_TrustedCertificatePolicy(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "trusted_certificate_policy" {
	load_balancer = "%s"
  name          = "acctest-%d"

  trusted_certificate_policy {
    trusted_certificate = "${opc_lbaas_certificate.trusted-cert.uri}"
  }
}

resource "opc_lbaas_certificate" "trusted-cert" {
  name = "acctest-%d"
  type = "TRUSTED"
  certificate_body = "-----BEGIN CERTIFICATE-----\nMIIFjzCCA3egAwIBAgIRAPHidcMOfXqzFXEVlGi/148wDQYJKoZIhvcNAQELBQAw\nPDEVMBMGA1UEChMMc2Nyb3Nzb3JhY2xlMSMwIQYDVQQDExpteXdlYmFwcC5zY3Jv\nc3NvcmFjbGUuc2l0ZTAeFw0xODA2MzAxMTM5MTVaFw0xOTA3MDQxNTM5MTVaMFsx\nCzAJBgNVBAYTAkNBMRAwDgYDVQQIEwdPbnRhcmlvMRUwEwYDVQQKEwxzY3Jvc3Nv\ncmFjbGUxIzAhBgNVBAMTGm15d2ViYXBwLnNjcm9zc29yYWNsZS5zaXRlMIICIjAN\nBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtl8cFy2b92+hop3p2DU1gp3OnQf6\nc7njsrRYfKRaCHd7L93aiVKGkccFNNHVzhBc/bScfEhYtqFCstMPBNVgTqqzulF0\nifoNDRxg2Dwbk8SsO0aWzmz38vEew9xLdLE4KAze0DVfgPB/nrr75osBCoLnxC1T\nRBMMzp26Ff/0S1W9k/D8zUGezOPKH26x8hUdQaVxS+wZQQrWfvsWkBgFdDvysZXi\nPdpo8v09neeh616FNrrT6Z9LP2VDvc3YToqFaDgV6Yy0Qm7zwz4IOpc6NQ1sP/Mx\ns9f36Hb+k1KXvAFt/dB5ROeCYDvQIMbEN+tjTENMRukjyfgc5MOb8qViSANPl6pe\nZtloP92l3bYOD8LzXFfqAGtaY0XBEn1Aw0TLWmlVPjyrsv7uXSIlCm53YOH0FtQX\nq6oikQCj0/S/Z+D0LpOM+JElvRb4JjiZIqGSVpIgxcZs23xt3GyQMAEm64CF5/y6\nIsB/E9W1Ue10u1/41CcPpc9pzS8QbxwJyIQPuZHljl/R7RpZMf9m8VU2/lwYfZfI\n5iGbdTA9oc/t62KFHtk/uxh9fG1V9T30AtlpvP6sHKdSW9AKjPlcrV+qxvdQQznW\nNI8HuaBO+ynducTnfCCVBC05okFQge2NlFPE0p47w8wu4s2Yv/OICxiXZsm6gERu\nDJqxBFWnOvj2HFECAwEAAaNtMGswEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0T\nAQH/BAIwADAfBgNVHSMEGDAWgBSJNATOWcuXm+Jv0H7UCkIId2M0WTAlBgNVHREE\nHjAcghpteXdlYmFwcC5zY3Jvc3NvcmFjbGUuc2l0ZTANBgkqhkiG9w0BAQsFAAOC\nAgEAqcTBLaW4D5PcEYSwMhNYqdACCuV6mc1o18PzIDHn/VDqF8pVzvaSTdEuMTte\noz5W8JwBpG6jH8E3YKEMMC4f/CI09PdDM8nBr4yDOHlaTIt1jRWFjG7gBGfe6rZw\ntBkrz9fteU80ST8LBcEFnoov7Txss54amS0L+vXU1ddwx6e6k9Ta3eAWMNn/JkVg\n63uCgiYueLT62AJUJZvBwPJdBnYASpJxh/AN8biIkWqWnoERVoofGfngiGoQ9DLo\niq6So6Ix4D95eOmIRpf/MC2yTTOeIxiQXi5LMk/NZ9oRUbc0JOinMVFfLStynnlP\n6xu0RjBKtCO2EjiRWl8sdQIVEgY3MicDoCoBt5HwJdIBkR955l4or5aFjxuLeV2E\nn2q2RacaUqV/xp46RCjm1hbJYkwcWXHnQzHkx6Jk6Y5kDcYp85CkUnBfzdfsC6/I\nbKByd3Sfp2wWvB9f+D1rI2ZsfOc18N/S+9AM/SVn0WHIr1DZBm2yaABe1m5Qo4jq\nXlRYAfjVg+tdhEwy6X6V8APW2ZLKXoVHTBl6XEc0Kqgths3r551nIHci9S1skMaN\ngO09dPxaJG+TcDKkO3YaxppszSY/IJa+h4nFv6j/mr+3tYRJ+Qs70tu+pqm2qcXI\noqW2x+i+3LFQD4vclv0aqrhVfeerpiSXpcGGyg65p9KtBCg=\n-----END CERTIFICATE-----\n"
  certificate_chain = "-----BEGIN CERTIFICATE-----\nMIIFazCCA1OgAwIBAgIQIPffd3HSy1azL08m3+003TANBgkqhkiG9w0BAQsFADA8\nMRUwEwYDVQQKEwxzY3Jvc3NvcmFjbGUxIzAhBgNVBAMTGm15d2ViYXBwLnNjcm9z\nc29yYWNsZS5zaXRlMB4XDTE4MDYzMDExMzkxNVoXDTE5MDcwNDE1MzkxNVowPDEV\nMBMGA1UEChMMc2Nyb3Nzb3JhY2xlMSMwIQYDVQQDExpteXdlYmFwcC5zY3Jvc3Nv\ncmFjbGUuc2l0ZTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBALZfHBct\nm/dvoaKd6dg1NYKdzp0H+nO547K0WHykWgh3ey/d2olShpHHBTTR1c4QXP20nHxI\nWLahQrLTDwTVYE6qs7pRdIn6DQ0cYNg8G5PErDtGls5s9/LxHsPcS3SxOCgM3tA1\nX4Dwf566++aLAQqC58QtU0QTDM6duhX/9EtVvZPw/M1Bnszjyh9usfIVHUGlcUvs\nGUEK1n77FpAYBXQ78rGV4j3aaPL9PZ3noetehTa60+mfSz9lQ73N2E6KhWg4FemM\ntEJu88M+CDqXOjUNbD/zMbPX9+h2/pNSl7wBbf3QeUTngmA70CDGxDfrY0xDTEbp\nI8n4HOTDm/KlYkgDT5eqXmbZaD/dpd22Dg/C81xX6gBrWmNFwRJ9QMNEy1ppVT48\nq7L+7l0iJQpud2Dh9BbUF6uqIpEAo9P0v2fg9C6TjPiRJb0W+CY4mSKhklaSIMXG\nbNt8bdxskDABJuuAhef8uiLAfxPVtVHtdLtf+NQnD6XPac0vEG8cCciED7mR5Y5f\n0e0aWTH/ZvFVNv5cGH2XyOYhm3UwPaHP7etihR7ZP7sYfXxtVfU99ALZabz+rByn\nUlvQCoz5XK1fqsb3UEM51jSPB7mgTvsp3bnE53wglQQtOaJBUIHtjZRTxNKeO8PM\nLuLNmL/ziAsYl2bJuoBEbgyasQRVpzr49hxRAgMBAAGjaTBnMA4GA1UdDwEB/wQE\nAwICBDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSJNATOWcuXm+Jv0H7UCkII\nd2M0WTAlBgNVHREEHjAcghpteXdlYmFwcC5zY3Jvc3NvcmFjbGUuc2l0ZTANBgkq\nhkiG9w0BAQsFAAOCAgEAVREHaF6qPPe27c/IMThVLDDzihaT83DWjzcel/OI4MPP\nsd1HV01tWtt56jP4IjzH1SpdHQMncTbSAbEOwsbclmrLe4E/Hd58Dzjo6apkKx2M\nieX+XVBi0KJ5pKh+OJHug8CGpnFu7IWla5zUiRY2Mm4Y3EdZNn4NH0smd8Expqck\nqehI0xsuN4blj3KFRtmgN1Zm48qSavah9PfGpicCPs1ZvoJJ8v17DFE4uFbkGqZl\nRRFpCybOmW7KeU5v8lDhmkcP6bu72xw+J7VGT0TfHotXTLXSPNjRlD13m1idvk0o\nXiosdLoQWvMpq71mstG3b11fwCA/EXuJkgANTxTkpjo7S5fWvgDUPqaVTt9nSbt0\nPHID1OvxfKZeuNIB0hM0oA+C5ZbSULuWTEaHPIwM3xgM+I7gCoJItJpzruyrtSjE\nUNJlMlo9zoJptx/a6ZguIvyu95MQbDnTJfq8sZjK1r0mxMBvx9tE8qTHXgAkuIC3\nFpDuFtfIDUgiWweSk5js19/deiP+tQ2abd/Z8MCR++e0bHNMdyyXS9CahOcSWCCJ\nHomAUmji594MTlP37MfkufA9NGegIwACf0VqE6FWrriO6VThvNpjnkNBewttlymu\nRshjxzWs/8bVI3HyOIz3CVh2gD2477D+kDsJJICBxkmz2eizt8EUcZwVsXO0KRs=\n-----END CERTIFICATE-----\n"
}
%s
`, lbID, rInt, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_UpdatePolicyType_Create(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "update_policy_type" {
  load_balancer = "%s"
  name          = "acctest-%d"

  application_cookie_stickiness_policy {
    cookie_name = "MY_APP_COOKIE"
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSPolicyConfig_UpdatePolicyType_Update(lbID string, rInt, lbCount int) string {
	return fmt.Sprintf(`
resource "opc_lbaas_policy" "update_policy_type" {
  load_balancer = "%s"
  name          = "acctest-%d"

	redirect_policy {
    redirect_uri = "https://redirect.example.com"
    response_code = 306
  }
}
%s
`, lbID, rInt, testAccParentLoadBalancerConfig(lbCount, rInt))
}

func testAccLBaaSCheckPolicyExists(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.PolicyClient()

	if _, err := client.GetPolicy(lb, name); err != nil {
		return fmt.Errorf("Error retrieving state of Policy '%s': %v", name, err)
	}

	return nil
}

func testAccLBaaSCheckPolicyDestroyed(state *OPCResourceState) error {
	lb := getLoadBalancerContextFromID(state.Attributes["load_balancer"])
	name := state.Attributes["name"]

	client := testAccProvider.Meta().(*Client).lbaasClient.PolicyClient()

	if info, _ := client.GetPolicy(lb, name); info != nil {
		return fmt.Errorf("Policy '%s' still exists: %+v", name, info)
	}
	return nil
}

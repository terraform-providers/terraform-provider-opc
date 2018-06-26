package opc

import "testing"

func TestValidateIPPrefixCIDR(t *testing.T) {
	validPrefixes := []string{
		"10.0.1.0/24",
		"10.1.0.0/16",
		"192.168.0.1/32",
		"10.20.0.0/18",
		"10.0.12.0/24",
	}

	for _, v := range validPrefixes {
		_, errors := validateIPPrefixCIDR(v, "prefix")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IP Address Prefix: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"10.0.0.1/35",
		"192.168.1.256/16",
		"256.0.1/16",
	}

	for _, v := range invalidPrefixes {
		_, errors := validateIPPrefixCIDR(v, "prefix")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid IP Address", v)
		}
	}
}

func TestValidateAdminDistance(t *testing.T) {
	validDistances := []int{
		0,
		1,
		2,
	}

	for _, v := range validDistances {
		_, errors := validateAdminDistance(v, "distance")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Admin Distance: %q", v, errors)
		}
	}

	invalidDistances := []int{
		-1,
		4,
		3,
		42,
	}

	for _, v := range invalidDistances {
		_, errors := validateAdminDistance(v, "distance")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid Admin Distance", v)
		}
	}
}

func TestValidateIPProtocol(t *testing.T) {
	validProtocols := []string{
		"all",
		"ah",
		"esp",
		"icmp",
		"icmpv6",
		"igmp",
		"ipip",
		"gre",
		"mplsip",
		"ospf",
		"pim",
		"rdp",
		"sctp",
		"tcp",
		"udp",
	}

	for _, v := range validProtocols {
		_, errors := validateIPProtocol(v, "ip_protocol")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Admin Distance: %q", v, errors)
		}
	}

	invalidProtocols := []string{
		"bad",
		"real bad",
		"are you even trying at this point?",
	}
	for _, v := range invalidProtocols {
		_, errors := validateIPProtocol(v, "ip_protocol")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid IP Protocol", v)
		}
	}

}

func TestValidateComputeStorageAccount(t *testing.T) {
	validStorageAccounts := []string{
		"/Compute-hasicorp/cloud_storage",
		"/Compute-abc123456/cloud_storage",
		"/Compute-123456789/cloud_storage",
	}

	for _, v := range validStorageAccounts {
		_, errors := validateComputeStorageAccountName(v, "account")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid two part Storage Account name: %q", v, errors)
		}
	}

	invalidStorageAccounts := []string{
		"cloud_storage",
		"/Storage-mydomain",
		"/Compute-mydomain",
		"/Storage-mydomain:user@example.com",
		"/Storage-mydomain/user@example.com",
		"/Compute-mydomain/user@example.com",
		"/Compute-mydomain/cloud-storage",
		"/Compute-mydomain/cloud_storage/",
	}

	for _, v := range invalidStorageAccounts {
		_, errors := validateComputeStorageAccountName(v, "account")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid two part Storage Account name", v)
		}
	}
}

func TestValidateLoadBalancerResourceName(t *testing.T) {
	validNames := []string{
		"abc123ABC",
		"abc-123-ABC",
		"a",
	}

	for _, v := range validNames {
		_, errors := validateLoadBalancerResourceName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q is not a valid resource name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"under_score",
		"/resourcename",
		"resource/name",
		"-name",
		"name-",
	}

	for _, v := range invalidNames {
		_, errors := validateLoadBalancerResourceName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid resource name", v)
		}
	}
}

func TestValidateLoadBalancerPolicyName(t *testing.T) {
	validNames := []string{
		"abc123ABC",
		"abc-123-ABC",
		"abc_123_ABC",
		"a",
	}

	for _, v := range validNames {
		_, errors := validateLoadBalancerPolicyName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q is not a valid resource name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"/resourcename",
		"resource/name",
		"-name",
		"name-",
		"_name",
		"name_",
	}

	for _, v := range invalidNames {
		_, errors := validateLoadBalancerPolicyName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid resource name", v)
		}
	}
}

func TestValidateLoadBalancerID(t *testing.T) {
	validNames := []string{
		"us-central-1/lb1",
	}

	for _, v := range validNames {
		_, errors := validateLoadBalancerID(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q is not a valid Load Balancer ID: %q", v, errors)
		}
	}

	invalidNames := []string{
		"us-central-1",
		"/",
		"us-central-1/",
		"/lb1",
	}

	for _, v := range invalidNames {
		_, errors := validateLoadBalancerID(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid Load Balancer ID", v)
		}
	}
}

func TestValidateComputeResourceName(t *testing.T) {
	validResourceNames := []string{
		"/Compute-hasicorp/user/name",
		"/Compute-abc123456/user@example.com/resource-name",
		"/Compute-55555555/Some_Path/resource_name",
		"/Compute-55555555/Some_Path/resource_name/subitem",
	}

	for _, v := range validResourceNames {
		_, errors := validateComputeResourceFQDN(v, "account")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid three part resoure name: %q", v, errors)
		}
	}

	invalidResourceNames := []string{
		"/Storage-hasicorp/user/name",
		"/Compute-abc-123456/user@example.com/resource-name",
		"/Compute-idcs-blahblah/Some_Path/resource_name",
		"/Compute-/Some_Path/resource_name",
		"/Compute-/user/name",
		"/Compute-hasicorp//name",
		"/Compute-hasicorp/user/",
		"/Compute-hasicorp//",
	}

	for _, v := range invalidResourceNames {
		_, errors := validateComputeResourceFQDN(v, "account")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid three part resource name", v)
		}
	}
}

func TestValidateOriginServerPoolURI(t *testing.T) {
	validURIs := []string{
		"https://lbaas-148cba7050494081b95151c522617ba9.balancer.oraclecloud.com/vlbrs/uscom-central-1/lb1/originserverpool/pool1",
	}

	for _, v := range validURIs {
		_, errors := validateOriginServerPoolURI(v, "account")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid originserverpool URI %q", v, errors)
		}
	}

	invalidURIs := []string{
		"https://lbaas-148cba7050494081b95151c522617ba9.balancer.oraclecloud.com/vlbrs/uscom-central-1/lb1/listeners/pool1",
	}
	for _, v := range invalidURIs {
		_, errors := validateOriginServerPoolURI(v, "account")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid originserverpool URI", v)
		}
	}
}

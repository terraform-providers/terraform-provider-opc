package opc

import (
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/go-oracle-terraform/compute"
)

// Validate whether an IP Prefix CIDR is correct or not
func validateIPPrefixCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error while parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
		return
	}
	return
}

// Admin distance can either be a 0, 1, or a 2. Defaults to 0.
func validateAdminDistance(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 0 || value > 2 {
		errors = append(errors, fmt.Errorf(
			"%q can only be an interger between 0-2. Got: %d", k, value))
	}
	return
}

// Admin distance can either be a 0, 1, or a 2. Defaults to 0.
func validateIPProtocol(v interface{}, k string) (ws []string, errors []error) {
	validProtocols := map[string]struct{}{
		string(compute.All):    {},
		string(compute.AH):     {},
		string(compute.ESP):    {},
		string(compute.ICMP):   {},
		string(compute.ICMPV6): {},
		string(compute.IGMP):   {},
		string(compute.IPIP):   {},
		string(compute.GRE):    {},
		string(compute.MPLSIP): {},
		string(compute.OSPF):   {},
		string(compute.PIM):    {},
		string(compute.RDP):    {},
		string(compute.SCTP):   {},
		string(compute.TCP):    {},
		string(compute.UDP):    {},
	}

	value := v.(string)
	if _, ok := validProtocols[value]; !ok {
		errors = append(errors, fmt.Errorf(
			`%q must contain a valid Image owner , expected ["all",	"ah",	"esp", "icmp",	"icmpv6",	"igmp",	"ipip",	"gre",	"mplsip",	"ospf",	"pim",	"rdp",	"sctp",	"tcp",	"udp"] got %q`,
			k, value))
	}
	return
}

// Check storage account name matches required format
func validateComputeStorageAccountName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("^/Compute-([a-zA-Z0-9]*)/cloud_storage$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"%s is not a valid storage account name (/Compute-identity_domain/cloud_storage)", value))
	}
	return
}

// Name can contain only alphanumeric characters and hyphens. First and last characters cannot be hyphen.
func validateLoadBalancerResourceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("^[a-zA-Z0-9](([a-zA-Z0-9-]*)([a-zA-Z0-9]))?$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"Name \"%s\" must contain only alphanumeric characters and hyphens. First and last characters cannot be hyphen", value))
	}
	return
}

// Name can contain only alphanumeric characters, hyphens, and underscores. First and last characters cannot be hyphen or underscore.
func validateLoadBalancerPolicyName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("^[a-zA-Z0-9](([a-zA-Z0-9-_]*)([a-zA-Z0-9]))?$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"Name \"%s\" contain only alphanumeric characters, hyphens, and underscores. First and last characters cannot be hyphen or underscore", value))
	}
	return
}

// Check Load Balancer refernece matches the two part "Region/Name" format
func validateLoadBalancerID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("^([a-zA-Z0-9-]+)/([a-zA-Z0-9-]+)$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"Load Balancer ID \"%s\" must be in the format \"region/name\"", value))
	}
	return
}

// Check name matchs the Compute three part name format /Compute-{domain}/{container}/{name}
func validateComputeResourceFQDN(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("^/Compute-([a-zA-Z0-9]+)/([^/]+)/([a-zA-Z0-9-_/]+)$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"Name \"%s\" must alphanumeric characters and hyphen only, first and last characters cannot be hyphen", value))
	}
	return
}

// Check URI matches the Load Balancer Origin Server Pool URI format
func validateOriginServerPoolURI(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if match, _ := regexp.MatchString("https://([-A-Za-z0-9+@_.]+)/vlbrs/([a-zA-Z0-9-]+)/([a-zA-Z0-9-]+)/originserverpool/([a-zA-Z0-9-]+)$", value); match != true {
		errors = append(errors, fmt.Errorf(
			"\"%s\" must be a valid URI to an Load Balancer Origin Server Pool resource", value))
	}
	return
}

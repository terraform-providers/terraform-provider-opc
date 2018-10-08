---
layout: "opc"
page_title: "Oracle: opc_compute_vpn_endpoint_v2"
sidebar_current: "docs-opc-resource-vpn-endpoint-v2"
description: |-
  Creates and manages a VPN Endpoint V2 in an Oracle Cloud Infrastructure Compute Classic identity domain.
---

# opc\_compute\_vpn_endpoint_v2

The ``opc_compute_vpn_endpoint_v2`` resource creates and manages an VPN Endpoint V2 in an Oracle Cloud Infrastructure Compute Classic identity domain.

## Example Usage

```hcl
resource "opc_compute_ip_network" "default" {
	name = "default_ip_network"
	ip_address_prefix = "10.0.12.0/24"
}

resource "opc_compute_vpn_endpoint_v2" "default" {
    name = "default_vpn_endpoint_v2
    customer_vpn_gateway = "127.0.0.1"
    ip_network = "${opc_compute_ip_network.default.name}"
    pre_shared_key = "defaultpsk"
    reachable_routes = ["127.0.0.1/24"]
    vnic_sets = ["default"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPN Endpoint V2.

* `customer_vpn_gateway` - (Required) The ip address of the VPN gateway in your data center through which you want to connect to the Oracle Cloud VPN gateway.

* `ip_network` - (Required) The name of the IP network on which the cloud gateway is created by VPNaaS.

* `pre_shared_key` - (Required) The pre-shared VPN key

* `reachable_routes` - (Required) A list of routes (CIDR prefixes) that are reachable through this VPN tunnel.

* `vnic_sets` - (Required) A list of vnic sets that traffics is allowed to and from. 

* `description` - (Optional) A description of the VPN Endpoint V2.

* `enabled` - (Optional) Enables or disables the VPN Endpoint V2. Set to true by default.

* `ike_identifier` - (Optional) The Internet Key Exchange (IKE) ID. If you don't specify a value, the default value is the public IP address of the cloud gateway.

* `require_perfect_forward_secrecy` - (Optional) Boolean specificying whether Perfect Forward Secrecy is enabled. Set to true by default.

* `phase_one_settings` - (Optional) Settings for the phase one protocol (IKE). Phase One Settings are detailed below.

* `phase_two_settings` - (Optional) Settings for the phase two protocol (IPSEC). Phase Two Settings are detailed below.

* `tags` - (Optional) List of tags that may be applied to the VPN Endpoint V2.

Phase One Settings support the following:

* `encryption` - (Required) Encryption options for IKE.

* `hash` - (Required) Authentication options for IKE. 

* `dh_group` - (Required) Diffie-Hellman group for both IKE and ESP. 

Phase Two Settings support the following: 

* `encryption` - (Required) Encryption options for IKE.

* `hash` - (Required) Authentication options for IKE. 

In addition to the above, the following values are exported:

* `uri` - The Uniform Resource Identifier for the VPN Endpoint V2.

## Import

VPN Endpoint V2's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_compute_vpn_endpoint_v2.default example
```

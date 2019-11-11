---
subcategory: "Compute Classic"
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
resource "opc_compute_vnic_set" "vnicset1" {
  name = "vpnaas-vnicset1"
}

resource "opc_compute_ip_network" "ipnetwork1" {
	name              = "ipnetwork1"
	ip_address_prefix = "10.0.12.0/24"
}

resource "opc_compute_vpn_endpoint_v2" "vpnaas1" {
  name                 = "vpnaas1"
  customer_vpn_gateway = "${var.vpn_endpoint_public_ip}"
  ip_network           = "${opc_compute_ip_network.ipnetwork1.name}"
  pre_shared_key       = "${var.pre_shared_key}"
  reachable_routes     = ["172.16.4.0/24"]
  vnic_sets            = ["${opc_compute_vnic_set.vnicset1.name}"]
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

* `require_perfect_forward_secrecy` - (Optional) Boolean specifying whether Perfect Forward Secrecy is enabled. Set to true by default.

* `phase_one_settings` - (Optional) Settings for the phase one protocol (IKE). Phase One Settings are detailed below.

* `phase_two_settings` - (Optional) Settings for the phase two protocol (IPSEC). Phase Two Settings are detailed below.

* `tags` - (Optional) List of tags that may be applied to the VPN Endpoint V2.

Phase One Settings support the following:

* `encryption` - (Required) IKE Encryption. `aes128`, `aes192` or `aes256`  

* `hash` - (Required) IKE Hash. `sha1`, `sha2_256` or `md5`

* `dh_group` - (Required) Diffie-Hellman group for both IKE and ESP. `group2`, `group5`, `group14`, `group22`, `group23`, or `group24`

* `lifetime` - (Optional) IKE Lifetime in seconds.

Phase Two Settings support the following:

* `encryption` - (Required) ESP Encryption.  `aes128`, `aes192` or `aes256`  

* `hash` - (Required) ESP Hash. `sha1`, `sha2_256` or `md5`

* `lifetime` - (Optional) IPSEC Lifetime in seconds.

In addition to the above, the following values are exported:

* `local_gateway_ip_address` - Public IP Address of the Local Gateway.

* `local_gateway_private_ip_address` - Private IP Address of the Local Gateway.

* `uri` - The Uniform Resource Identifier for the VPN Endpoint V2.

## Import

VPN Endpoint V2's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_compute_vpn_endpoint_v2.vpnaas1 /Compute-mydomain/user/example
```

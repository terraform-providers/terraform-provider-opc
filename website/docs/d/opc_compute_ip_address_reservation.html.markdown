---
layout: "opc"
page_title: "Oracle: opc_compute_ip_address_reservation"
sidebar_current: "docs-opc-datasource-ip-address-reservation"
description: |-
  Gets information about an existing IP Network IP address reservation.
---

# opc\_compute\_ip\_address\_reservation

Use this data source to access the attributes of an existing IP Network IP Address Reservation.

## Example Usage

```hcl
data "opc_compute_ip_address_reservation" "example" {
  name = "/Compute-${var.domain}/${var.user}/ipaddressreservation1"
}

output "public_ip_address" {
  value = "${data.opc_compute_ip_address_reservation.example.ip_address}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip address reservation.

## Attributes Reference

* `ip_address_pool` - The IP address pool from which the IP address is allocated.

* `description` - A description of the ip address reservation.

* `tags` - List of tags that applied to the IP address reservation.

* `ip_address` - The reserved IPv4 Public IP address.

* `uri` - The Uniform Resource Identifier of the ip address reservation

* `fqdn` - The Fully Qualified Domain Name

---
layout: "opc"
page_title: "Oracle: opc_compute_ip_reservation"
sidebar_current: "docs-opc-datasource-ip-reservation"
description: |-
  Gets information about an existing IP reservation for the Shared Network.
---

# opc\_compute\_ip\_reservation

Use this data source to access the attributes of an existing Shared Network IP Reservation.

## Example Usage

```hcl
data "opc_compute_ip_reservation" "example" {
  name = "/Compute-${var.domain}/${var.user}/309dd783-552d-4b8a-a3aa-a71109f703df"
}`

output "public_ip_address" {
  value = "${data.opc_compute_ip_reservation.example.ip}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the IP Reservation.

## Attributes Reference

* `ip` - The reserved IPv4 Public IP address.

* `permanent` - Whether the IP address remains reserved even when it is no longer associated with an instance.

* `parent_pool` - The pool from which to allocate the IP address.

* `tags` - List of tags applied to the IP reservation.

* `used` - indicates that the IP reservation is associated with an instance.

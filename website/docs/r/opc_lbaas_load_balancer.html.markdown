---
layout: "opc"
page_title: "Oracle: opc_lbaas_load_balancer"
sidebar_current: "docs-opc-resource-lbaas-load-balancer"
description: |-
  Creates an Oracle Load Balancer Classic instance.
---

# opc\_lbaas\_load\_balancer

The `opc_lbaas_load_balancer` resource creates and manages a Load Balancer Classic instance in an Oracle Cloud Infrastructure Classic Compute region. You must define server pools, create at least one listener, and optionally define the policies for the load balancer before it will be operational.

## Example Usage

```hcl
resource "opc_lbaas_load_balancer" "lb1" {
  name        = "example-lb1"
  region      = "uscom-central-1"
  description = "My Example Load Balancer"
  scheme      = "INTERNET_FACING"

  permitted_methods = ["GET", "HEAD", "POST", "PUT"]  
}
```

## Argument Reference

* `name` - (Required) The name of the Load Balancer.

* `description` - (Optional) A short description for the load balancer. The description must not exceed 1000 characters.

* `enabled` - (Optional) Boolean flag to enable or disable the Load Balancer. Default is `true` (enabled).

* `ip_network` - (Optional) Fully qualified three part name of the IP network to be associated with the load balancer.

* `parent_load_balancer` - (Optional) Select a parent load balancer if you want to create a dependent load balancer.

* `permitted_clients` - (Optional) List of permitted client IP addresses or CIDR ranges which can connect to this load balancer on the configured Listener ports. If not set all connections are permitted.

* `permitted_methods` - (Optional) List of permitted HTTP methods. e.g. `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD` or you can also create your own custom methods. Requests with methods not listed in this field will result in a 403 (unauthorized access) response.

* `region` - (Required) The region in which to create the Load Balancer, e.g. `uscom-central-1`

* `scheme` - (Required) Set to either `INTERNET_FACING` or `INTERNAL`

  - `INTERNET_FACING` - Create an internet-facing load balancer in a given IP network.

  - `INTERNAL` - Create an internal load balancer in a given IP network for sole consumption of other clients inside the same network.

* `tags` - (Optional) List of tags.

## Additional Attributes

In addition to the above, the following values are exported:

* `operational_details` - Description of the operational state.

* `state` - State of the Origin Server Pool.

* `status` - Status of the Origin Server Pool.

* `uri` - The Uniform Resource Identifier for the Load Balancer.

## Import

Load Balancers can be imported using the a combinable of the resource `region` and `name` in the format `region/name`

```shell
$ terraform import opc_lbaas_load_balancer.lb1 uscom-central-1/example-lb1
```

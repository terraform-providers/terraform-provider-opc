---
layout: "opc"
page_title: "Oracle: opc_lbaas_listener"
sidebar_current: "docs-opc-resource-opc_lbaas_listener"
description: |-
  Creates an Oracle Load Balancer Classic Listener.
---

# opc\_lbaas\_listener

The `opc_lbaas_listener` resource creates and manages a Load Balancer Classic Listener for a Load Balancer Classic instance.

## Example Usage

```hcl
resource "opc_lbaas_listener" "listener1" {
  load_balancer     = "${opc_lbaas_load_balancer.lb1.id}"
  name              = "listener1"

  port              = 80
  balancer_protocol = "HTTP"
  server_protocol   = "HTTP"
  server_pool       = "${opc_lbaas_server_pool.serverpool1.uri}"
  virtual_hosts     = ["mywebapp.example.com"]
  policies          = ["${opc_lbaas_policy.load_balancing_mechanism_policy.uri}"]
}
```

## Argument Reference

* `name` - (Required) The name of the Listener.

* `load_balancer` - (Required) The parent Load Balancer the Listener.

* `balancer_protocol` - (Required)  transport protocol that will be accepted for all incoming requests to the selected load balancer listener. `HTTP` or `HTTPS`

* `enabled` - (Optional) Boolean flag to enable or disable the Listener. Default is `true` (enabled).

* `path_prefixes` - (Optional) List of paths to configure the listener to accept only requests that are targeted to a specific path within the URI of the request.

* `polices` - (Optional) List of the Load Balancer Policy URIs to apply to the listener.

* `port` - (Required) The port on which the Load Balancer is listening.

* `server_pool` - (Optional) URI of the Server Pool resource to which the load balancer distributes requests.

* `server_protocol` - (Required) The protocol to be used for routing traffic to the origin servers in the server pool. `HTTP` or `HTTPS`

* `ssl_certificates` - (Optional) The URI of the server security certificate. If the `balancer_protocol` is set to either HTTPS or SSL then you must select a server certificate.

* `tags` - (Optional) List of tags.

* `virtual_hosts` - (Optional) Configure the listener to only accept URI requests that include the host names listed in this field.

## Additional Attributes

In addition to the above, the following values are exported:

* `operational_details` - Description of the operational state.

* `parent_listener` - The parent Listener.

* `state` - State of the Origin Server Pool.

* `uri` - The Uniform Resource Identifier for the Listner.

## Import

Listeners can be imported using the a combinable of the resource region, load balancer name and policy name and in the format `region/loadbalancer/name`

```shell
$ terraform import opc_lbaas_listener.listener1 uscom-central-1/lb1/example-listener1
```

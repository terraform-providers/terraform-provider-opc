---
layout: "opc"
page_title: "Oracle: opc_lbaas_server_pool"
sidebar_current: "docs-opc-resource-lbaas_server_pool"
description: |-
  Creates an Oracle Load Balancer Classic Origin Server Pool.
---

# opc\_lbaas\_server\_pool

The `opc_lbaas_server_pool` resource creates and manages a Load Balancer Classic Origin Server Pool for a Load Balancer Classic instance. The Server Pool defines one or more servers (referred to as origin servers) to which the load balancer can distribute requests.

## Example Usage

```hcl
resource "opc_lbaas_server_pool" "serverpool1" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "serverpool1"

  servers = ["129.150.169.179:8080", "129.150.170.162:8080"]

  health_checks {
    type = "http"
    path = "/healthcheck"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Server Pool.

* `load_balancer` - (Required) The parent Load Balancer the Origin Server Pool.

* `enabled` - (Optional) Boolean flag to enable or disable the Server Pool. Default is `true` (enabled).

* `health_checks` - (Optional) Enables Load Balancer health checks, see [Health Check Attributes](#health-check-attributes)

* `servers` - (Required) List of servers in the Server Pool. To define the server in the server pool, provide IP address or DNS name of the compute instances, and port for load balancer to direct traffic to, in the format `host:port`

* `tags` - (Optional) List of tags.

* `vnic_set` - (Optional) Fully qualified three part name of a vNICSet to be associated with the server pool vNIC. Load Balancer uses this vNICSet to set the right ACLs to allow egress traffic from the load balancer.

### Heath Check Attributes

The load balancer can perform regular health checks of the origin servers and route inbound traffic to the healthy origin servers.

* `type` - (Optional) Health check mechanism to use to test the origin servers. Options

  - `http` - sends an HTTP HEAD to the set `path` and checks response is in one of the `accepted_return_codes`

  Default is `http`

* `path` - (Optional) The path to check. Set the '/' the check all paths. Default is '/'

* `accepted_return_codes` - (Optional) List of HTTP response status codes that indicate the origin server is healthy. Accepted return codes can be one or more of the `2xx`, `3xx`, `4xx`, or `5xx` codes. Default is `["2xx","3xx"]`.

* `enabled` - (Optional) Boolean flag to enable or disable the Health Checks. Default is `true` (enabled).

* `interval` - (Optional) The approximate interval, in seconds, that the load balancer will wait before sending the target request to each origin server, in the range 5 to 300 seconds. Default is `30`.

* `timeout` - (Optional) The amount of time, in seconds, that the load balancer will wait without a response before identifying the origin server as unavailable. The timeout value must be less than the `interval` value and it should range between 2 to 60. Default is `60`.

* `healthy_threshold` - (Optional) The number of consecutive successful health checks required before moving the origin server to the healthy state. The value of healthy threshold ranges from 2 to 10. Default is `6`.

* `unhealthy_threshold` - (Optional) The number of consecutive health check failures required before moving the origin server to the unhealthy state. The value of unhealthy threshold ranges from 2 to 10. Default is `3`.

## Additional Attributes

In addition to the above, the following values are exported:

* `operational_details` - Description of the operational state.

* `state` - State of the Origin Server Pool.

* `status` - Status of the Origin Server Pool.

* `uri` - The Uniform Resource Identifier for the Server Pool.

## Import

Origin Server Pools can be imported using the combination of the resource region, load balancer name, and policy name and in the format `region/loadbalancer/name`

```shell
$ terraform import opc_lbaas_server_pool.serverpool1 uscom-central-1/lb1/example-serverpool1
```

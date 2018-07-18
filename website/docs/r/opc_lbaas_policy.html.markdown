---
layout: "opc"
page_title: "Oracle: opc_lbaas_policy"
sidebar_current: "docs-opc-resource-lbaas-policy"
description: |-
  Creates an Oracle Load Balancer Classic Policy.
---

# opc\_lbaas\_policy

The `opc_lbaas_policy` resource creates and manages a Load Balancer Classic Policy for a Load Balancer Classic instance.

The Policy resource supports the definition of distinct Policy types:

- [Application Cookie Stickiness Policy](#application-cookie-stickiness-policy)
- [CloudGate Policy](#cloudgate-policy)
- [Load Balancer Cookie Stickiness Policy](#load-balancer-cookie-stickiness-policy)
- [Load Balancing Mechanism Policy](#load-balancing-mechanism-policy)
- [Rate Limiting Request Policy](#rate-limiting-request-policy)
- [Redirect Policy](#redirect-policy)
- [Resource Access Control Policy](#resource-access-control-policy)
- [Set Request Header Policy](#set-request-header-policy)
- [SSL Negotiation Policy](#set-negotiation-policy)
- [Trusted Certificate Policy](#trusted-certificate-policy)

## Example Usage


```hcl
resource "opc_lbaas_policy" "load_balancing_mechanism_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_load_balancing_mechanism_policy"

  load_balancing_mechanism_policy {
    load_balancing_mechanism = "round_robin"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Listener.

* `load_balancer` - (Required) The parent Load Balancer the Listener.

Include one, and only one, of:

* `application_cookie_stickiness_policy` - see [Application Cookie Stickiness Policy](#application-cookie-stickiness-policy)

* `cloudgate_policy` - see [CloudGate Policy](#cloudgate-policy)

* `load_balancer_cookie_stickiness_policy` - see [Load Balancer Cookie Stickiness Policy](#load-balancer-cookie-stickiness-policy)

* `load_balancing_mechanism_policy` - see [Load Balancing Mechanism Policy](#load-balancing-mechanism-policy)

* `rate_limiting_request_policy` - see [Rate Limiting Request Policy](#rate-limiting-request-policy)

* `redirect_policy` - see [Redirect Policy](#redirect-policy)

* `resource_access_control_policy` - see [Resource Access Control Policy](#resource-access-control-policy)

* `set_request_header_policy` - see [Set Request Header Policy](#set-request-header-policy)

* `ssl_negotiation_policy` - see [SSL Negotiation Policy](#set-negotiation-policy)

* `trusted_certificate_policy` - see [Trusted Certificate Policy](#trusted-certificate-policy)

### Application Cookie Stickiness Policy

Enable session stickiness (session affinity) for any request based on a given cookie name specified in the policy.

#### Example

```hcl
resource "opc_lbaas_policy" "application_cookie_stickiness_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_application_cookie_stickiness_policy"

  application_cookie_stickiness_policy {
    cookie_name = "MY_APP_COOKIE"
  }
}
```

#### Attributes

* `cookie_name` - (Required) Name of the application cookie used to control how long the load balancer will continue to route requests to the same origin server.

### CloudGate Policy

protect resources/applications with the help of CloudGate module available in Load Balancer. These headers will enable CloudGate to lookup for the application and the policy present under the appropriate IDCS Tenant containing information for the protection mechanism to be enforced.

#### Example

```hcl
resource "opc_lbaas_policy" "cloudgate_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_cloudgate_policy"

  cloudgate_policy {
    virtual_hostname_for_policy_attribution = "host1.example.com"
  }
}
```

#### Attributes

* `virtual_hostname_for_policy_attribution` - (Required) Host name needed by CloudGate to enforce OAuth policies.


### Load Balancer Cookie Stickiness Policy

Enables session stickiness (session affinity) for all requests for a given period of time specified in the policy.

#### Example

```hcl
resource "opc_lbaas_policy" "load_balancer_cookie_stickiness_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_load_balancer_cookie_stickiness_policy"

  load_balancer_cookie_stickiness_policy {
    cookie_expiration_period = 60
  }
}
```

#### Attributes

* `cookie_expiration_period` - (Required) The time period, in seconds, after which the cookie should be considered stale. If the value is zero or negative the stickiness session lasts for the duration of the browser session.

### Load Balancing Mechanism Policy

Specify a load balancing mechanism for distributing client requests across multiple origin servers

#### Example

```hcl
resource "opc_lbaas_policy" "load_balancing_mechanism_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_load_balancing_mechanism_policy"

  load_balancing_mechanism_policy {
    load_balancing_mechanism = "round_robin"
  }
}
```

#### Attributes

* `load_balancing_mechanism` - (Required) Supported options are `round_robin`,
`least_conn`, and `ip_hash`.

### Rate Limiting Request Policy

Limits the number of requests that can be processed per second by the load balancer.

#### Example

```hcl
resource "opc_lbaas_policy" "rate_limiting_request_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_rate_limiting_request_policy"

  rate_limiting_request_policy {
    requests_per_second      = 1
    burst_size               = 10
    delay_excessive_requests = true
    zone                     = "examplezone"
  }
}
```

#### Attributes

* `burst_size` - (Required) The number of requests that can be delayed until it exceeds the maximum number specified as burst size in which case the request is terminated.

* `delay_excessive_requests` - (Required) delay excessive requests while requests are being limited.

* `requests_per_second` - (Required) Maximum number of requests per second.

* `zone` - (Required) Name of the shared memory zone.

* `http_error_code` - (Optional) Status code to return in response to rejected requests. You can specify any status code between 405 to 599. Default is `503`

* `logging_level` - (Optional) Logging level for cases when the server refuses to process requests due to rate exceeding, or delays request processing. Can be one of `info`, `notice`, `warn`, or `error`.  Default is `warn`

* `rate_limiting_criteria` - (Optional) Criteria based on which requests will be throttled. Default is `server`

  - `server` - limit the requests processed by the virtual server
  - `remote_address` - limit the processing rate of requests coming from a single IP address.
  - `host` - limit the processing rate of requests coming from a host.

* `zone_memory_size` - (Optional) Size of the shared memory occupied by the zone. Default is `10`


### Redirect Policy

Redirects all requests to this load balancer to a specific URI.

#### Example

```hcl
resource "opc_lbaas_policy" "redirect_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_redirect_policy"

  redirect_policy {
    redirect_uri = "https://redirect.example.com"
    response_code = 306
  }
}
```

#### Attributes

* `redirect_uri` - (Required) redirected requests to the specified URI.

* `response_code` - (Required) The exact 3xx response code to use when redirecting

### Resource Access Control Policy

Controls what clients have access to the load balancer, based on the IP address or the CIDR range of the incoming request.

#### Example

```hcl
resource "opc_lbaas_policy" "resource_access_control_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_resource_access_control_policy"

  resource_access_control_policy {
    disposition = "DENY_ALL"
    permitted_clients = ["10.0.0.0/16"]
  }
}
```

#### Attributes

* `disposition` - (Required) Default policy. `DENY_ALL` or `PERMIT_ALL`.

* `denied_clients` - (Optional) List of IP address or CIDR ranges identifying clients from which requests must be accepted by the Load Balancer

* `permitted_clients` - (Optional) IP address or CIDR ranges identifying clients from which requests must be denied by the Load Balancer.

### Set Request Header Policy

Inserts additional information into the standard HTTP headers of requests forwarded to a server pool.

#### Example

```hcl
resource "opc_lbaas_policy" "set_request_header_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_set_request_header_policy"

  set_request_header_policy {
    header_name                     = "X-Custom-Header"
    value                           = "foo-bar  "
    action_when_header_exists       = "OVERWRITE"
    action_when_header_value_is     = ["bar", "foo"]
  }
}
```

#### Attributes

* `header_name` - (Required)

* `action_when_header_exists` - (Required) action to be taken when a header exists in the request. Options: `NOOP`, `PREPEND`, `APPEND`, `OVERWRITE`, `CLEAR`

* `action_when_header_value_is` - (Optional) List if header values. Action is taken only when the header exists in the request and the header value matches one of the values provided.

* `action_when_header_value_is_not` - (Optional) List if header values. Action is taken only when the header exists in the request and the header value does not match any of the values provided.


### SSL Negotiation Policy

Define specific SSL protocols, ciphers, and server order preference for the secure connection

#### Example

```hcl
resource "opc_lbaas_policy" "ssl_negotiation_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_ssl_negitiation_policy"

  ssl_negotiation_policy {
    port = 8022
    server_order_preference = "ENABLED"
    ssl_protocol = ["SSLv3", "TLSv1.2"]
    ssl_ciphers = ["AES256-SHA"]
  }
}
```

#### Attributes

* `ssl_protocol` - (Required) Security protocols supported for incoming secure client connections to the associated listener. Supported options are `SSLv2`, `SSLv3`, `TLSv1`, `TLSv1.1`, `TLSv1.2`

* `port` - (Optional) The load balancer port for the the SSL protocols and the SSL ciphers.

* `server_order_preference` - (Optional) enable or disable the server order preference for secure connections to associated Listener. `ENABLED` or `DISABLED`.

* `ssl_ciphers` - (Optional) List of SSL ciphers supported for incoming secure client connections to the associated Listener.

### Trusted Certificate Policy

Identifies a trusted certificate, which the load balancer uses when making a secure connection to the compute instances in the server pool. If you are configuring a secure connection (HTTPS or SSL) between the load balancer and the origin servers, you must add this policy to the load balancer.

#### Example

```hcl
resource "opc_lbaas_policy" "trusted_certificate_policy" {
  load_balancer = "${opc_lbaas_load_balancer.lb1.id}"
  name          = "example_trusted_certificate_policy"

  trusted_certificate_policy {
    trusted_certificate = "${opc_lbaas_ssl_certificate.cert1.uri}"
  }
}
```

#### Attributes

* `trusted_certificate` - (Required) URI of the SSL Certificate

## Additional Attributes

In addition to the above, the following values are exported:

* `state` - State of the Policy.

* `type` - The Type of the Policy.

* `uri` - The Uniform Resource Identifier for the Policy.

## Import

Policies can be imported using the a combinable of the resource region, load balancer name and policy name and in the format `region/loadbalancer/name`

```shell
$ terraform import opc_lbaas_policy.policy1 uscom-central-1/lb1/example-policy1
```

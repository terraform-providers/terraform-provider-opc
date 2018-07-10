---
layout: "opc"
page_title: "Oracle: opc_lbaas_certificate"
sidebar_current: "docs-opc-resource-lbaas-certificate"
description: |-
  Creates an Oracle Load Balancer Classic SSL Certificate.
---

# opc\_lbaas\_certificate

The `opc_lbaas_certificate` resource creates and manages an Load Balancer Classic TLS/SSL Digital Certificate.

Server certificates are used to secure the connection between clients and the load balancers. Trusted certificates are used to secure the connection between the load balancer and the origin servers in the server pool.

## Example Usage

```hcl
resource "opc_lbaas_certificate" "example" {
  name             = "example-cert1"
  type             = "SERVER"
  private_key      = "${var.private_key_pem}"
  certificate_body = "${var.cert_pem}"
}
```

## Argument Reference

* `name` - (Required) The name of the Certificate.

* `certificate_body` - (Required) The Certificate data in PEM format.

* `type` - (Required) Sets the Certificate Type. `TRUSTED` or `SERVER`.

* `certificate_chain` - (Optional) PEM encoded bodies of all certificates in the chain up to and including the CA certificate. This is not need when the certificate is self signed.

* `private_key` - (Optional) The private key data in PEM format. Only required for Server Certificates



## Additional Attributes

In addition to the above, the following values are exported:

* `state` - The State of the Digital Certificate resource.

* `uri` - The Uniform Resource Identifier for the Certificate resource.

## Import

Digital Certificates can be imported using the `name` of the resource.

```shell
$ terraform import opc_lbaas_ssl_certificate.cert1 example-cert1
```

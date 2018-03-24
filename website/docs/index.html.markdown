---
layout: "opc"
page_title: "Provider: Oracle Cloud Infrastructure Classic"
sidebar_current: "docs-opc-index"
description: |-
  The Oracle Cloud Infrastructure Classic provider is used to interact with the many resources supported by the Oracle Cloud Infrastructure Classic services. The provider needs to be configured with credentials for the Oracle Cloud Account.
---

# Oracle Cloud Infrastructure Classic Provider

The Oracle Cloud Infrastructure Classic provider (formerly know as the Oracle Public Cloud provider) is used to interact with the many resources supported by the [Oracle Cloud Infrastructure Classic](http://cloud.oracle.com/classic) and [Oracle Cloud at Customer](https://cloud.oracle.com/cloud-at-customer) infrastructure services. The provider needs to be configured with credentials for the Oracle Cloud Account.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Oracle Cloud Infrastructure Classic provider
provider "opc" {
  user            = "..."
  password        = "..."
  identity_domain = "..."
  endpoint        = "..."
}

# Create an IP Reservation
resource "opc_compute_ip_reservation" "production" {
  parent_pool = "/oracle/public/ippool"
  permanent = true
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Optional) The username to use, generally your email address. It can also
  be sourced from the `OPC_USERNAME` environment variable.

* `password` - (Optional) The password associated with the username to use. It can also be sourced from
  the `OPC_PASSWORD` environment variable.

* `identity_domain` - (Optional) The Identity Domain name (for Traditional accounts) or Identity Service ID (for IDCS accounts) of the environment to use. It can also be sourced from the `OPC_IDENTITY_DOMAIN` environment variable.  

* `endpoint` - (Optional) The Compute Classic API endpoint to use, associated with your Oracle Cloud Account. This is known as the `REST Endpoint` within the Oracle portal. It can also be sourced from the `OPC_ENDPOINT` environment variable.

* `storage_endpoint` - (Optional) The Storage Classic API endpoint to use, associated with your Oracle Storage Cloud account. This is known as the `REST Endpoint` within the Oracle portal. Can also be set via the `OPC_STORAGE_ENDPOINT` environment variable.

* `storage_service_id` - (Optional) The Storage Service ID for authentication with the `storage_endpoint`  If not set the `identity_domain` value is used. Can also be set via the `OPC_STORAGE_SERVICE_ID` environment variable.

* `max_retries` - (Optional) The maximum number of tries to make for a successful response when operating on resources. It can also be sourced from the `OPC_MAX_RETRIES` environment variable. Defaults to 1.

* `insecure` - (Optional) Skips TLS Verification for using self-signed certificates. Should only be used if absolutely needed. Can also via setting the `OPC_INSECURE` environment variable to `true`.

## Testing

Credentials must be provided via the `OPC_USERNAME`, `OPC_PASSWORD`,
`OPC_IDENTITY_DOMAIN` and `OPC_ENDPOINT` environment variables in order to run
acceptance tests.

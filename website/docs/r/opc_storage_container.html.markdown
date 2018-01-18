---
layout: "opc"
page_title: "Oracle: opc_storage_container"
sidebar_current: "docs-opc-resource-storage-container"
description: |-
  Creates and manages a Container in the OPC Storage Domain. `storage_endpoint` must be set in the
  provider or environment to manage these resources.
---

# opc\_storage\_container

Creates and manages a Container in the OPC Storage Domain. `storage_endpoint` must be set in the
provider or environment to manage these resources.

## Example Usage

```hcl
resource "opc_storage_container" "default" {
  name = "storage-container-1"
  read_acls = [ ".r:example.com", ".rlistings" ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Container.

* `read_acls` - (Optional) The list of ACLs that grant read access. See [Setting Container ACLs](#setting-container-acls).

* `write_acls` - (Optional) The list of ACLs that grant write access. See [Setting Container ACLs](#setting-container-acls).

* `primary_key` - (Optional) The primary secret key value for temporary URLs.

* `secondary_key` - (Optional) The secondary secret key value for temporary URLs.

* `allowed_origins` - (Optional) List of origins that are allowed to make cross-origin requests.

* `exposed_headers` - (Optional) List of headers exposed to the user agent (e.g. browser) in the actual request response

* `max_age` - (Optional) Maximum age in seconds for the origin to hold the preflight results.

* `quota_bytes` - (Optional) Maximum size of the container, in bytes

* `quota_count` - (Optional) Maximum object count of the container

* `metadata` - (Optional) Additional object metadata headers. See [Container Metadata ](#container-metadata) below for more information.

## Setting Container ACLs

The `read_acl` consists of a list of **roles** or **referrer designations**. The `write_acls` consists of a list of **roles**.

- **roles** can be built-in roles or custom roles. Custom roles are defined in the Users tab in the Oracle Cloud My Services console. For a role that was provisioned as part of another service instance, the format is `domainName.serviceName.roleName`. For a custom role, the format is `domainName.roleName`.  Default Storage roles include:

  - `${var.domain}.Storage.Storage_ReadOnlyGroup`
  - `${var.domain}.Storage.Storage_ReadWriteGroup`
  - `${var.domain}.Storage.Storage_Administrator`

- **referrer designation** indicates the host (or hosts) for which read access to the container should be allowed or denied. When the server receives a request for the container, it compares the referrer designations specified in the Read ACL with the value of the Referer header in the request, and determines whether access should be allowed or denied. The syntax of the referrer designation is: `.r:value`

  `value` indicates the host for which access to the container should be allowed. It can be a specific host name (example: `.r:www.example.com`), a domain (example: `.r:.example.com`), or an asterisk (`.r:*`) to indicate all hosts. Note that if `.r:*` is specified, objects in the container will be publicly readable without authentication.

  A minus sign (-) before value (example: `.r:-temp.example.com`) indicates that the host specified in the value field must be denied access to the container.

By default, read access to a container does not include permission to list the objects in the container. To allow listing of objects as well, include the `.rlistings` directive in the ACL.


## Container Metadata

The `metadata` config defines a map of additional meta data header name value pairs. The additional meta data items are set as HTTP Headers on the container in the form `X-Container-Meta-{name}: {value}`, where `{name}` is the name of the metadata item  `{value}` is the header content. For example:

```hcl
metadata {
  "Foo-Bar" = "barfoo"
  "Tags" = "abc 123 xyz"
}
```

## Import

Container's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_storage_container.default example
```

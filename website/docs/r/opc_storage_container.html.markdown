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
  name        = "storage-container-1"
  read_acls   = ["read_acl_1", "read_acl_2"]
  max_age = 60
  primary_key = "primary-key-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Container.

* `read_acls` - (Optional) The list of ACLs that grant read access.

* `write_acls` - (Optional) The list of ACLs that grant write access.

* `primary_key` - (Optional) The primary secret key value for temporary URLs.

* `secondary_key` - (Optional) The secondary secret key value for temporary URLs.

* `allowed_origins` - (Optional) List of origins that are allowed to make cross-origin requests.

* `max_age` - (Optional) Maximum age in seconds for the origin to hold the preflight results.

* `quota_bytes` - (Optional) Maximum size of the container, in bytes

* `quota_count` - (Optional) Maximum object count of the container


## Import

Container's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_storage_container.default example
```

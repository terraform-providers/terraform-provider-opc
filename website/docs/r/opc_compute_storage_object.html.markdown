---
layout: "opc"
page_title: "Oracle: opc_storage_object"
sidebar_current: "docs-opc-resource-storage-object"
description: |-
  Creates and manages a Container Object in the OPC Storage Domain. `storage_endpoint` must be set in the
  provider or environment to manage these resources.
---

# opc\_storage\_object

Creates and manages a Container in the OPC Storage Domain. `storage_endpoint` must be set in the
provider or environment to manage these resources.

## Example Usage

```hcl
resource "opc_storage_container" "foo" {
  name = "my-storage-container"
  max_age = 50
  primary_key = "test-key"
  allowed_origins = ["origin-1"]
}

resource "opc_storage_object" "default" {
  name        = "my-storage-object"
  container   = "${opc_storage_container.foo.name}"
  content     = <<EOF
FOO BAR BAZ
File Contents that will be supplied as the storage object
EOF
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Object.

* `container` - (Required) The name of the Storage Container in which to place the object.

**Note:** One of `content`, `file`, or `copy_from` must be specified

* `content` - (Optional) A raw string value for the Storage Object's contents.

* `file` - (Optional) The path to the file to use as the Storage Object's contents.

* `copy_from` - (Optional) The source Storage Object to copy from. (`container_name/object_name`).

* `content_disposition` - (Optional) Overrides the behavior of the browser.

* `content_encoding` - (Optional) Set the content-encoding metadata for the Storage Object.

* `content_type` - (Optional) Sets the MIME type for the object. Will be computed via the API if not specified.

* `delete_at` - (Optional) Specify the number of seconds after which the system deletes the storage object.

* `etag` - (Optional) MD5 checksum value of the request body. Unquoted, strongly recommended, but not required.

* `transfer_encoding` - (Optional) Sets the transfer encoding. Can only be `chunked` if set. Requires `content_length` to be `0` if set.

## Attributes Reference

The following attributes are exported:

* `accept_ranges` - Type of ranges that the object accepts.

* `content_length` - Length of the Storage Object in bytes.

* `last_modified` - Date and time that the object was created/modified in ISO 8601.

* `object_manifest` - The dynamic large-object manifest object.

* `timestamp` - Date and Time in UNIX EPOCH when the account, container, or object was initially created at the current version.

* `transaction_id` - Transaction ID of the request. Used for bug reports.

## Import

Storage Object's can be imported using `container_name/object_name`, e.g.
```shell
$ terraform import opc_storage_object.test my_container/my_object
```

Please note though, importing a Storage Object does _not_ allow a user to modify the content, or attributes for the Storage Object.
It is, however, possible to import a Storage Object, and replace the object with new content, or a copy of another Storage Object.
It is also possible to import a Storage Object into Terraform in order to delete the object.
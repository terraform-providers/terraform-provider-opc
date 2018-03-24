---
layout: "opc"
page_title: "Oracle: opc_storage_object"
sidebar_current: "docs-opc-resource-storage-object"
description: |-
  Creates and manages a Object in an Oracle Cloud Infrastructure Storage Classic container. `storage_endpoint` must be set in the provider or environment to manage these resources.
---

# opc\_storage\_object

Creates and manages a Object in an Oracle Cloud Infrastructure Storage Classic container. `storage_endpoint` must be set in the provider or environment to manage these resources.

## Example Usage

```hcl
resource "opc_storage_object" "default" {
  name         = "storage-object-1"
  container    = "${opc_storage_container.container.name}"
  file         = "${"./source_file.txt"}"
  etag         = "${md5(file("./source_file.txt"))}"
  content_type = "text/plain;charset=utf-8"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Object.

* `container` - (Required) The name of Storage Container the store the object in.

**Note:** One of `content`, `file`, or `copy_from` must be specified

* `content` - (Optional) Raw content in string-form of the data.

* `file` - (Optional) File path for the content to use for data.

* `copy_from` - (Optional) name of an existing object used to create the new object as a copy. The value is in form `container/object`. You must UTF-8-encode and then URL-encode the names of the container and object.

* `content_disposition` - (Optional) Set the HTTP `Content-Disposition` header to specify the override behaviour for the browser, e.g. `inline` or `attachment`.

* `content_encoding` - (Optional) set the HTTP `Content-Encoding` for the object.

* `content_type` - (Optional) set the MIME type for the object.

* `delete_at` - (Optional) The date and time in UNIX Epoch time stamp format when the system removes the object.

* `etag` - (Optional) MD5 checksum value of the request body. Strongly Recommended.

* `transfer_encoding` - (Optional) Set to `chunked` to enable chunked transfer encoding.

* `metadata` - (Optional) Additional object metadata headers. See [Object Metadata ](#object-metadata) below for more information.

## Attributes

In addition to the attributes listed above, the following attributes are exported:

* `id` - The combined container and object name path of the object.
* `accept_ranges` - Type of ranges that the object accepts.
* `content_length` - Length of the object in bytes.
* `last_modified` - Date and Time that the object was created/modified in ISO 8601.
* `object_manifest` - The dynamic large-object manifest object.
* `timestamp` - Date and Time in UNIX EPOCH when the account, container, or object was initially created at the current version.
* `transaction_id` - Transaction ID of the request.

## Object Metadata

The `metadata` config defines a map of additional meta data header name value pairs. The additional meta data items are set HTTP Headers on the object in the form `X-Object-Meta-{name}: {value}`, where `{name}` is the name of the metadata item  `{value}` is the header content. For example:

```hcl
metadata {
  "Foo-Bar" = "barfoo",
  "Sha256" = "e91ed4f93637379a7539cb5d8d0b5bca3972755de4f9371ab2e123e7b4c53680"
}
```

## Import

Object's can be imported using the `resource id`, e.g.

```shell
$ terraform import opc_storage_object.default container/example
```

Please note though, importing a Storage Object does _not_ allow a user to modify the content, or attributes for the Storage Object. It is, however, possible to import a Storage Object, and replace the object with new content, or a copy of another Storage Object. It is also possible to import a Storage Object into Terraform in order to delete the object.

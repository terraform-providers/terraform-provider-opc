---
layout: "opc"
page_title: "Oracle: opc_compute_storage_volume_attachment"
sidebar_current: "docs-opc-resource-storage-volume-attachment"
description: |-
  Creates and manages a storage volume attachment in an Oracle Cloud Infrastructure Compute Classic identity domain.
---

# opc\_compute\_storage\_volume\_attachment

The `opc_compute_storage_volume_attachment` resource creates and manages a storage volume attachment in an Oracle Cloud Infrastructure Compute Classic identity domain.

## Example Usage

```hcl
resource "opc_compute_storage_volume" "default" {
  name = "storage-volume-1"
  size = 1
}

resource "opc_compute_instance" "default" {
  name = "instance-1"
  label = "instance-1"
  shape = "oc3"
  image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
}

resource "opc_compute_storage_attachment" "test" {
  instance = "${opc_compute_instance.default.name}"
  storage_volume = "${opc_compute_storage_volume.default.name}"
  index = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the instance the volume will be attached to.

* `storage_volume` - (Required) The name of the storage volume that will be attached to the
 instance

* `index` - (Required) The index on the instance that the storage volume will be attached to.

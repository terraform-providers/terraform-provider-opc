---
layout: "opc"
page_title: "Oracle: opc_compute_storage_volume"
sidebar_current: "docs-opc-resource-storage-volume-type"
description: |-
  Creates and manages a storage volume in an OPC identity domain.
---

# opc\_compute\_storage\_volume

The ``opc_compute_storage_volume`` resource creates and manages a storage volume in an OPC identity domain.

~> **Caution:** The ``opc_compute_storage_volume`` resource can completely delete your storage volume just as easily as it can create it. To avoid costly accidents, consider setting [``prevent_destroy``](/docs/configuration/resources.html#prevent_destroy) on your storage volume resources as an extra safety measure.

## Example Usage

```hcl
resource "opc_compute_storage_volume" "test" {
  name        = "storageVolume1"
  description = "Description for the Storage Volume"
  size        = 10
  tags        = ["bar", "foo"]
}
```

## Example Usage (Bootable Volume)
```hcl
data "opc_compute_image_list_entry" "test" {
  image_list = "my_image_list"
  version    = 1
}

resource "opc_compute_storage_volume" "test" {
  name             = "storageVolume1"
  description      = "Description for the Bootable Storage Volume"
  size             = 30
  tags             = ["first", "second"]
  bootable         = true
  image_list       = "${data.opc_compute_image_list_entry.test.image_list}"
  image_list_entry = "${data.opc_compute_image_list_entry.test.version}"
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) The name for the Storage Account.
* `description` (Optional) The description of the storage volume.
* `size` (Required) The size of this storage volume in GB. The allowed range is from 1 GB to 2 TB (2048 GB).
* `storage_type` - (Optional) - The Type of Storage to provision. Possible values are `/oracle/public/storage/latency` or `/oracle/public/storage/default`. Defaults to `/oracle/public/storage/default`.
* `bootable` - (Optional) Is the Volume Bootable? Defaults to `false`.
* `image_list` - (Optional) Defines an image list.
* `image_list_entry` - (Optional) Defines an image list entry.
* `snapshot` - (Optional) The name of the parent snapshot from which the storage volume is restored or cloned.
* `snapshot_id` - (Optional) The Id of the parent snapshot from which the storage volume is restored or cloned.
* `snapshot_account` - (Optional) The Account of the parent snapshot from which the storage volume is restored.
* `tags` - (Optional) Comma-separated strings that tag the storage volume.

## Attributes Reference

The following attributes are exported:

* `hypervisor` - The hypervisor that this volume is compatible with.
* `machine_image` - Name of the Machine Image - available if the volume is a bootable storage volume.
* `managed` - Is this a Managed Volume?
* `platform` - The OS platform this volume is compatible with.
* `readonly` - Can this Volume be attached as readonly?
* `status` - The current state of the storage volume.
* `storage_pool` - The storage pool from which this volume is allocated.
* `uri` - Unique Resource Identifier of the Storage Volume.

## Import

Storage Volume's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_compute_storage_volume.volume1 example
```

<a id="timeouts"></a>
## Timeouts

`opc_compute_storage_volume` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for Creating Storage Volumes.
- `update` - (Default `30 minutes`) Used for Modifying Storage Volumes.
- `delete` - (Default `30 minutes`) Used for Deleting Storage Volumes.


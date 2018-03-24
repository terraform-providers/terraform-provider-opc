---
layout: "opc"
page_title: "Oracle: opc_compute_image_list_entry"
sidebar_current: "docs-opc-datasource-image-list-entry"
description: |-
  Gets information about the configuration of an Image List Entry within an Oracle Cloud Infrastructure Compute Classic domain.
---

# opc\_compute\_image\_list\_entry

Use this data source to access the configuration of an Image List Entry.

## Example Usage

```hcl
data "opc_compute_image_list_entry" "foo" {
  image_list = "my_image_list"
  version = "version_of_my_list"
}

output "machine_images" {
  value = "${data.opc_compute_image_list_entry.foo.machine_images}"
}

output "single_image" {
  value = "${data.opc_compute_image_list_entry.foo.machine_images[1]}"
}
```

## Argument Reference
* `image_list` - (Required) - The name of the image list to lookup.
* `version` - (Required) - The version (integer) of the Image List to use.
* `entry` - (Optional) - Which machine image to use. See [Entry](#entry) below for more details

## Entry
The `entry` argument is fully optional when configuring the Data Source. If specified, however,
the returned array of machine images will have a length of 1, and only contain the desired image.

Thus, if "my_image_list" is an image list that contains the following images:

```
["image1", "image2", "image3", "image4", "image5"]
```

Then specifing an `entry` of `3`, the returned `machine_images` array will have a sigle element:
`"image3"`. If `entry` was omitted, or set to `0`, the returned `machine_images` array will contain
all of the images for that image list version.

If the supplied `entry` value is invalid for the image list, Terraform will exit with an error,
that the desired image list entry was not found.

## Attributes Reference

* `dns` - Array of DNS servers for the interface.
* `attributes` - JSON object of all of the image list's attributes
* `machine_images` - An array of machine images as strings
* `uri` - The URI of the image list

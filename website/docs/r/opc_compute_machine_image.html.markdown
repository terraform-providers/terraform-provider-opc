---
layout: "opc"
page_title: "Oracle: opc_compute_machine_image"
sidebar_current: "docs-opc-resource-machine-image"
description: |-
  Creates and manages a Machine Image in a Oracle Cloud Infrastructure Compute Classic identity domain.
---

# opc\_compute\_machine\_image

The ``opc_compute_machine_image`` resource creates and manages a machine image template of a virtual hard disk of a specific size with an installed operating system.

Before performing this creating the Machine Image, you must upload your machine image file to Oracle Cloud Infrastructure Object Storage Classic `compute_images` container


## Example Usage

```hcl
resource "opc_compute_machine_image" "centos" {
  account     = "/Compute-${var.domain}/cloud_storage"
  name        = "CentOS_7"
  file        = "CentOS-7-x86_64-OracleCloud.raw.tar.gz"
  description = "CentOS 7"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required) The two part name of the compute object storage account in the format `/Compute-{identity_domain}/cloud_storage`

* `name` - (Required) The name of the Machine Image.

* `file` - (Required) The name of the Machine Image .tar.gz file in the `compute_images` storage container.

* `description` - (Optional) A description of the Machine Image.

* `attributes` - (Optional) An optional JSON object of arbitrary attributes to be made available to the instance. These are user-defined tags. After defining attributes, you can view them from within an instance at http://192.0.0.192/

In addition to the above, the following values are exported:

* `error_reason` - Description of the state of the machine image if there is an error.

* `hypervisor` -  Dictionary of hypervisor-specific attributes.

* `image_format` - The format of the image.

* `platform` - The OS platform of the image.

* `state` - The state of the uploaded machine image.

* `uri` - The Uniform Resource Identifier for the Machine Image.

## Import

Machine Images can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_compute_machine_image.machine_image1 example
```

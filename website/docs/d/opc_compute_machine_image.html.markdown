---
layout: "opc"
page_title: "Oracle: opc_compute_machine_image"
sidebar_current: "docs-opc-datasource-machine-image"
description: |-
  Gets information about the configuration of an Machine Image in a Compute Classic identity domain.
---

# opc\_compute\_machine\_image

Use this data source to access the configuration of an Machine Image.

## Example Usage

```hcl
data "opc_compute_machine_image" "foo" {
  account = "/Compute-${var.domain}/cloud_storage"
  name = "Microsoft_Windows_Server_2012_R2"
}

output "attributes" {
  value = "${data.opc_compute_machine_image.foo.attributes}"
}
```

## Argument Reference

* `account` - (Required) The two part name of the compute object storage account in the format `/Compute-{identity_domain}/cloud_storage`

* `name` - (Required) The name of the Machine Image.

## Attributes Reference

* `file` - The name of the Machine Image .tar.gz file in the `compute_images` storage container.

* `description` - A description of the Machine Image.

* `attributes` - An optional JSON object of arbitrary attributes to be made available to the instance. These are user-defined tags. After defining attributes, you can view them from within an instance at http://192.0.0.192/

* `error_reason` - Description of the state of the machine image if there is an error.

* `hypervisor` -  Dictionary of hypervisor-specific attributes.

* `image_format` - The format of the image.

* `platform` - The OS platform of the image.

* `state` - The state of the uploaded machine image.

* `uri` - The Uniform Resource Identifier for the Machine Image.

* `fqdn` - The Fully Qualified Domain Name
---
layout: "opc"
page_title: "Oracle: opc_compute_ssh_key"
sidebar_current: "docs-opc-datasource-ssh-key"
description: |-
  Gets information about an existing SSH key.
---

# opc\_compute\_ssh_key

Use this data source to access the attributes of an SSH Key.

## Example Usage

```hcl
data "opc_compute_ssh_key" "test" {
  name    = "/Compute-${var.domain}/${var.user}/test-key"
}

output "public_ssh_key" {
  value = "${data.opc_compute_ssh_key.test.key}"
}
```

## Argument Reference

* `name` - (Required) The unique (within this identity domain) name of the SSH key.

## Attributes Reference


* `key` - The public SSH key.

* `enabled` - Whether or not the key is enabled.

* `fqdn` - The Fully Qualified Domain Name
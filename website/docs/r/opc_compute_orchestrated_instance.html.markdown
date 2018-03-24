---
layout: "opc"
page_title: "Oracle: opc_compute_orchestrated_instance"
sidebar_current: "docs-opc-resource-orchestrated-instance"
description: |-
  Creates and manages an Orchestration containing a number of instances in an Oracle Cloud Infrastructure Compute Classic identity domain.
---

# opc\_compute\_orchestrated\_instance

The `opc_compute_orchestrated_instance` resource creates and manages an orchestration containing a number of
instances in an Oracle Cloud Infrastructure Compute Classic identity domain.

## Example Usage

```hcl
resource "opc_compute_orchestrated_instance" "default" {
  name          = "default_orchestration"
  desired_state = "active"

  instance {
    name       = "default-orchestrated-instance"
    label      = "Orchestrated Instance 1 Label"
    shape      = "oc3"
    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
  }
}
```

## Example Usage with Multiple Instances

```hcl
resource "opc_compute_orchestrated_instance" "default" {
  name          = "default_orchestration"
  desired_state = "active"

  instance {
    name       = "default-instance-1"
    label      = "Instance One"
    shape      = "oc3"
    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
  }

  instance {
    name       = "default-instance-2"
    label      = "Instance Two"
    shape      = "oc3"
    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"
  }
}
```

## Example Usage with IP Networking

```hcl
resource "opc_compute_ip_network" "default" {
  name              = "default-ip-network"
  description       = "testing-ip-network-instance"
  ip_address_prefix = "10.1.12.0/24"
}

resource "opc_compute_orchestrated_instance" "default" {
  name          = "default_orchestration"
  desired_state = "active"

  instance {
    name       = "default-instance"
    label      = "Default Instance"
    shape      = "oc3"
    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"

    networking_info {
      index          = 0
      ip_network     = "${opc_compute_ip_network.default.id}"
      vnic           = "default-ip-network"
      shared_network = false
    }
  }
}
```

## Example Usage with Storage

```hcl
resource "opc_compute_storage_volume" "foo" {
  name = "acc-test-orchestration-%d"
  size = 1
}

resource "opc_compute_storage_volume" "bar" {
  name = "acc-test-orchestration-2-%d"
  size = 1
}

resource "opc_compute_orchestrated_instance" "default" {
  name          = "test_orchestration-%d"
  desired_state = "active"

  instance {
    name       = "default-instance"
    label      = "Default Instance"
    shape      = "oc3"
    image_list = "/oracle/public/OL_7.2_UEKR4_x86_64"

    storage {
      volume = "${opc_compute_storage_volume.foo.name}"
      index  = 1
    }

    storage {
      volume = "${opc_compute_storage_volume.bar.name}"
      index  = 2
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the orchestration.

* `desired_state` - (Required) The desired state of the orchestration. Permitted values are:

  - `active`: all resource (instances) declared in the orchestration are created

  - `suspend`: all resources (instances) declared in the orchestration are removed unless the instances has
`persistent = true`

  - `inactive`:  all resources (instances) declared in the orchestration are removed including the instances that have
`persistent = true`

* `instance` - (Required) The information pertaining to creating an instance through the orchestration API.

* `description` - (Optional) The description of the orchestration.

## Instance

Instance supports the arguments found in [opc_compute_instance](https://www.terraform.io/docs/providers/opc/r/opc_compute_instance.html)
and the following:

* `persistent` - (Optional) Determines whether the instance will persist when the orchestration is suspended.
Defaults to false.

In addition to the above, the following values are exported:

* `uri` - The Uniform Resource Identifier for the Orchestration

* `version` - (Optional) The version of the orchestration.

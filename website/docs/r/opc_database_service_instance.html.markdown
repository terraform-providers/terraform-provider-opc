---
layout: "opc"
page_title: "Oracle: opc_database_service_instance"
sidebar_current: "docs-opc-resource-service-instance"
description: |-
  Creates and manages a service instance in an OPC identity domain.

---

# opc\_database\_service\_instance

The ``opc_database_service_instance`` resource creates and manages a service instance in an OPC identity domain.

## Example Usage

```hcl
resource "opc_database_service_instance" "default" {
  name        = "service-instance-1"
  description = "This is a description for an service instance"
  edition = "EE_EP"
  level = "PAAS"
  shape = "oc1m"
  subscription_type = "HOURLY"
  version = "12.2.0.1"
  vm_public_key = "A ssh public key"
  parameter {
    admin_password = "Test_String7"
    backup_destination = "NONE"
    sid = "ORCL"
    usable_storage = 15
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Instance.

* `edition` - (Required) Database edition for the service instance. Possible values are `SE`, `EE`, `EE_HP`, or `EE_EP`.

* `level` - (Required) Service level for the service instance. Possible values are `BASIC` or `PAAS`.

* `shape` - (Required) Desired compute shape. Possible values are `oc3`, `oc4`, `oc5`, `oc6`, `oc1m`, `oc2m`, `oc3m`, or `oc4m`.

* `subscription_type` - (Required) Billing unit. Possible values are `HOURLY` or `MONTHLY`.

* `version` - (Required) Oracle Database software version; one of: `12.2.0.1`, `12.1.0.2`, or `11.2.0.4`.

* `vm_public_key` - (Required) Public key for the secure shell (SSH). This key will be used for authentication when connecting to the Database Cloud Service instance using an SSH client.

* `parameter` - (Optional) Additional configuration for a service instance. This set is required if level is PAAS. Parameter is documented below.

* `ibkup` - (Optional) Specify if the service instance's database should, after the instance is created, be replaced by a database stored in an existing cloud backup that was created using Oracle Database Backup Cloud Service. IBKUP is documented below.

* `cloud_storage` - (Optional) Provides Cloud Storage for service instance backups. Cloud Storage
is documented below

* `description` - (Optional) A description of the Service Instance.

Parameter supports the following:

* `admin_password` - (Required) Password for Oracle Database administrator users sys and system. The password must meet the following requirements: Starts with a letter. Is between 8 and 30 characters long. Contains letters, at least one number, and optionally, any number of these special characters: dollar sign `$`, pound sign `#`, and underscore `_`.

* `backup_destination` - (Required) Backup Destination. Possible values are `BOTH`, `OSS`, `NONE`.

* `char_set` - (Required) Character Set for the Database Cloud Service Instance. All possible values are listed under the [parameters section documentation](http://docs.oracle.com/en/cloud/paas/database-dbaas-cloud/csdbr/op-paas-service-dbcs-api-v1.1-instances-%7BidentityDomainId%7D-post.html). Default value is `AL32UTF8`.

* `usable_storage` - (Required) Storage size for data (in GB). Minimum value is `15`. Maximum value depends on the backup destination: if `BOTH` is specified, the maximum value is `1200`; if `OSS` or `NONE` is specified, the maximum value is `2048`.

* `disaster_recovery` - (Optional) Specify if an Oracle Data Guard configuration is created using the Disaster Recovery option or the High Availability option. Valid values are `yes` and `no`. Default value is no.

* `failover_database` - Specify if an Oracle Data Guard configuration comprising a primary database and a standby database is created. Valid values are `yes` and `no`. Default value is `no`.

* `golden_gate` - (Optional) Specify if the database should be configured for use as the replication database of an Oracle GoldenGate Cloud Service instance. Valid values are `yes` and `no`. Default value is `no`. You cannot set `goldenGate` to `yes` if either `is_rac` or `failoverDatabase` is set to `yes`.

* `is_rac` - (Optional) Specify if a cluster database using Oracle Real Application Clusters should be configured. Valid values are `yes` and `no`. Default value is `no`.

* `n_char_set` - (Optional) National Character Set for the Database Cloud Service instance. Valid values are `AL16UTF16` and `UTF8`. Default value is `AL16UTF16`.

* `pdb_name` - (Optional) This attribute is valid when Database Cloud Service instance is configured with version 12c. Pluggable Database Name for the Database Cloud Service instance. Default value is `pdb1`.

* `sid` - (Optional) Database Name for the Database Cloud Service instance. Default value is `ORCL`.

* `source_service_name` - (Optional) Indicates that the service instance should be created as a "snapshot clone" of another service instance. Provide the name of the existing service instance whose snapshot is to be used.

* `snapshot_name` - (Optional) The name of the snapshot of the service instance specified by sourceServiceName that is to be used to create a "snapshot clone". This parameter is valid only if source_service_name is specified.

* `timezone` - (Optional) Time Zone for the Database Cloud Service instance. Default value is `UTC`.

* `type` - (Optional) Component type to which the set of parameters applies. Defaults to `db`

* `db_demo` - (Optional) Indicates whether to include the Demos PDB. Valid values are `yes` or `no`.

IBKUP supports the following:

* `cloud_storage_username` - (Required) Username of the Oracle Cloud user.

* `cloud_storage_password` - (Required) Password of the Oracle Cloud user specified in `ibkup_cloud_storage_user`.

* `database_id` - (Required) Database id of the database from which the existing cloud backup was created.

* `decryption_key` - (Optional) Password used to create the existing, password-encrypted cloud backup. This password is used to decrypt the backup. Specify either `ibkup_decryption_key` or `ibkup_wallet_file_content` for decrypting the backup.

* `wallet_file_content` - (Optional) String containing the xsd:base64Binary representation of the cloud backup's wallet file. This wallet is used to decrypt the backup. Specify either `ibkup_decryption_key` or `ibkup_wallet_file_content` for decrypting the backup.

Cloud Storage supports the following:

* `container` - (Required) Name of the Oracle Storage Cloud Service container used to provide storage for your service instance backups. Use the following format to specify the container name: `<storageservicename>-<storageidentitydomain>/<containername>`

* `username` - (Required) Username for the Oracle Storage Cloud Service administrator.

* `password` - (Required) Password for the Oracle Storage Cloud Service administrator.

* `create_if_missing` - (Optional) Specify if the given cloud_storage_container is to be created if it does not already exist. Default value is `false`.


In addition to the above, the following values are exported:

* `apex_url` - The URL to use to connect to Oracle Application Express on the service instance.

* `backup_supported_version` - The version of cloud tooling for backup and recovery supported by the service instance.

* `compute_site_name` - The Oracle Cloud location housing the service instance.

* `connect_descriptor_with_public_ip` - The connection descriptor for Oracle Net Services (SQL*Net) with IP addresses instead of host names.

* `created_by` - The user name of the Oracle Cloud user who created the service instance.

* `creation_time` - The date-and-time stamp when the service instance was created.

* `current_version` - The Oracle Database version on the service instance, including the patch level.

* `dbaasmonitor_url`- The URL to use to connect to Oracle DBaaS Monitor on the service instance.

* `em_url` - The URL to use to connect to Enterprise Manager on the service instance.

* `glassfish_url` - The URL to use to connect to the Oracle GlassFish Server Administration Console on the service instance.

* `hdg_prem_ip` - Data Guard Role of the on-premise instance in Oracle Hybrid Disaster Recovery configuration.

* `hybrid_dg` - Indicates whether the service instance hosts an Oracle Hybrid Disaster Recovery configuration.

* `identity_domain` - The identity domain housing the service instance.

* `ip_network` - This attribute is only applicable to accounts where regions are supported. The three-part name of an IP network to which the service instance is added. For example: /Compute-identity_domain/user/object

* `ip_reservations` - Groups one or more IP reservations in use on this service instance. This attribute is only applicable to accounts where regions are supported.

* `jaas_instances_using_service` - The Oracle Java Cloud Service instances using this Database Cloud Service instance.

* `listener_port` - The listener port for Oracle Net Services (SQL*Net) connections.

* `num_ip_reservations` - The number of Oracle Compute Service IP reservations assigned to the service instance.

* `num_nodes` - The number of compute nodes in the service instance.

* `rac_database` - Indicates whether the service instance hosts an Oracle RAC database.

* `region` - This attribute is only applicable to accounts where regions are supported. Location where the service instance is provisioned (only for accounts where regions are supported).

* `sm_plugin_version` - The version of the cloud tooling service manager plugin supported by the service instance.

* `total_shared_storage` - For service instances hosting an Oracle RAC database, the size in GB of the storage shared and accessed by the nodes of the RAC database.

* `uri` - The Uniform Resource Identifier for the Service Instance

## Import

Service Instance's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_database_service_instance.instance1 example
```

---
layout: "opc"
page_title: "Oracle: opc_java_service_instance"
sidebar_current: "docs-opc-resource-service-instance"
description: |-
  Creates and manages a java service instance in an OPC identity domain.

---

# opc\_java\_service\_instance

The ``opc_java_service_instance`` resource creates and manages a java service instance in an OPC identity domain.

## Example Usage

```hcl
resource "opc_database_service_instance" "default" {
  ...
}

resource "opc_java_service_instance" "test" {
  name = "test-java-service-instance-%d"
  type = "weblogic"
  shape = "oc1m"
  version = "12.2.1"
  public_key = "ssh-public-key"
  cloud_storage {
    container = "Storage-identitydomain/default"
    create_if_missing = true
  }
  database {
    name = "${opc_database_service_instance.default.name}"
    username = "sys"
    password = "Pass_Test9"
  }
  admin {
    username = "terraform-user"
    password = "Pass_Test8"
  }
}
```

The following is an example of how to provision a service instance with the Oracle Traffic Director:

```hcl
resource "opc_database_service_instance" "default" {
  ...
}

resource "opc_java_service_instance" "default" {
  name = "test-java-service-instance-%d"
  subscription_type = "HOURLY"
  weblogic {
    version = "12.2.1"
    edition = "EE"
    public_key = "ssh-public-key"
    shape = "oc1m"
    database {
      name = "${opc_database_service_instance.default.name}"
      username = "sys"
      password = "Pass_Test9"
    }
    admin {
      username = "terraform-user"
      password = "Pass_Test8"
    }
  }
  cloud_storage {
    container = "Storage-identitydomain/default"
    create_if_missing = true
  }
  otd {
    admin {
      username = "otd-admin"
      password = "Pass_Test8"
    }
    shape = "oc1m"
    public_key = "ssh-public-key"
  }
}
```

The following is an example of how to provision a service instance with the DataGrid Scaling Units:

```hcl
resource "opc_database_service_instance" "default" {
  ...
}

resource "opc_java_service_instance" "default" {
  name = "test-java-service-instance-%d"
  subscription_type = "HOURLY"
  weblogic {
    version = "12.2.1"
    edition = "EE"
    public_key = "ssh-public-key"
    shape = "oc1m"
    database {
      name = "${opc_database_service_instance.default.name}"
      username = "sys"
      password = "Pass_Test9"
    }
    admin {
      username = "terraform-user"
      password = "Pass_Test8"
    }
  }
  cloud_storage {
    container = "Storage-identitydomain/default"
    create_if_missing = true
  }
  datagrid {
    cluster_name = "test-datagrid"
    scaling_unit_count = 1
    scaling_unit {
      name = "SMALL"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Instance.

* `cloud_storage` - (Required) Provides Cloud Storage for service instance backups. Cloud Storage
is documented below

* `subscription_type` - (Required) Billing unit. Possible values are `HOURLY` or `MONTHLY`.

* `weblogic` - (Required) The attributes required to create a WebLogic server alongside the java service instance.
WebLogic is documented below.

* `otd` - (Optional) The attributes required to create an Oracle Traffic Director (Load balancer). OTD is
documented below.

* `datagrid` - (Optional) The attributes required to enable DataGrid on the service instance. DataGrid is
documented below.

* `level` - (Optional) Service level for the service instance. Possible values are `BASIC` or `PAAS`. Default
value is `PAAS`.

* `backup_destination` - (Optional) Specifies whether to enable backups for this Oracle Java Cloud Service instance.
Valid values are `BOTH` or `NONE`. Defaults to `BOTH`.

* `description` - (Optional) Provides additional on the java service instance.

* `enable_admin_console` - (Optional) Flag that specifies whether to enable (true) or disable (false) the access
rules that control external communication to the WebLogic Server Administration Console, Fusion Middleware Control,
and Load Balancer Console.

* `ip_network` - (Optional) The three-part name of a custom IP network to attach this service instance to. For example: `/Compute-identity_domain/user/object`.

* `public_network` - (Optional) Path to the network from which the Oracle Java Cloud Service REST API will be accessed.

* `region` - (Optional) Name of the region where the Oracle Java Cloud Service instance is to be provisioned.
This attribute is only applicable to accounts where regions are supported. A region name must be specified if you
want to use ipReservations or ipNetwork.

* `sample_app_deployment_requested` - (Optional) Flag that specifies whether to automatically deploy and start the
sample application, sample-app.war, to the Managed Server in your service instance. The default value is false.

* `ip_reservations` - (Optional) A single IP reservation name. This can be specified multiple times. `region` must be
specified is using this attribute.

Cloud Storage supports the following:

* `container` - (Required) Name of the Oracle Storage Cloud Service container used to provide storage for your service instance backups. Use the following format to specify the container name: `<storageservicename>-<storageidentitydomain>/<containername>`

* `username` - (Optional) Username for the Oracle Storage Cloud Service administrator. If left unspecified,
the username for Oracle Public Cloud is used.

* `password` - (Optional) Password for the Oracle Storage Cloud Service administrator. If left unspecified,
the password for Oracle Public Cloud is used.

* `create_if_missing` - (Optional) Specify if the given cloud_storage_container is to be created if it does not already exist. Default value is `false`.

WebLogic supports the following:

* `edition` - (Required) The edition for the service instance. Possible values are `SE`, `EE`, or `SUITE`.

* `database` - (Required) Information about the database deployment on Oracle Database Cloud Service. Database
is documented below.

* `shape` - (Required) Desired compute shape.

* `version` - (Required) Oracle java software version.

* `admin` - (Required) Admin information for the WebLogic Server. Admin is documented below.

* `public_key` - (Required) The public key for the secure shell (SSH). This key will be used for authentication
when connecting to the Oracle Java Cloud Service instance using an SSH client.

* `app_db` - (Optional) Details of Database Cloud Service database deployments that host application schemas.
Multiple can be specified. App DB is specified below.

* `backup_volume_size` - (Optional) Size of the backup volume for the service. The value must be a multiple of GBs.
You can specify this value in bytes or GBs. If specified in GBs, use the following format:
nG, where n specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs.
For example: 100000000000 or 10G. This value defaults to the system configured volume size.

* `cluster_name` - (Optional) - Specifies the name of the cluster that contains the Managed Servers
for the service instance.

* `connect_string` - (Optional) - Connection string for the database. The connection string must be entered using one
of the following formats: host:port:SID, host:port/serviceName.

* `content_port` - (Optional) - Port for accessing the deployed applications using HTTP. Default value is 8001.

* `deployment_channel_port` - Port for accessing the Administration Server using WLST. Default value is 9001.

* `domain` - (Optional) Information about the WebLogic domain. Domain is documented below.

* `managed_servers` - (Optional) Details information about the managed servers the java service instance will
look after. Managed Servers is documented below.

* `mw_volume_size` - (Optional) Size of the MW_HOME disk volume for the service (/u01/app/oracle/middleware).
The value must be a multiple of GBs. You can specify this value in bytes or GBs.
If specified in GBs, use the following format: nG, where n specifies the number of GBs.
For example, you can express 10 GBs as bytes or GBs. For example: 100000000000 or 10G.
This value defaults to the system configured volume size.

* `node_manager` - (Optional) Node Manager is a WebLogic Server utility that enables you to start, shut down,
and restart Administration Server and Managed Server instances from a remote location. Node Manager is documented
below.

* `pdb_service_name` - (Optional) Name of the pluggable database for Oracle Database 12c. The default value is the
pluggable database name when the database was created.

* `privileged_ports` - (Optional) A block of privileged port specifications. Privileged ports are specified below.

* `secured_admin_port` -  (Optional) Port for accessing the Administration Server using HTTPS. The default value is
7002.

* `upper_stack_product_name` - (Optional) The Oracle Fusion Middleware product installer to add to this Oracle Java
Cloud Service instance. Valid values are `ODI` (Oracle Data Integrator) or `WCP` (Oracle WebCenter Portal)

OTD supports the following:

* `admin` - (Required) Admin information for the Oracle Traffic Director. Admin is documented below.

* `shape` - (Required) Desired compute shape.

* `public_key` - (Required) The public key for the secure shell (SSH). This key will be used for authentication
when connecting to the Oracle Java Cloud Service instance using an SSH client.

* `high_availability` - (Optional) Flag that specifies whether load balancer HA is enabled.
This value defaults to false (that is, HA is not enabled).

* `listener` - (Optional) Specifies the type and number of the listener port. Listener is documented below.

* `load_balancing_policy` - (Optional) Policy to use for routing requests to the load balancer. Valid policies
include: `least_connection_count`, `least_response_time`, `round_robin`. The default value is
`least_connection_count`.

* `privileged_ports` - (Optional) A block of privileged port specifications. Privileged ports are specified below.

* `secured_listener_port` - (Optional) Secured listener port for accessing the deployed applications using HTTPS. The
default value is 8081.

DataGrid supports the following:

* `cluster_name` - (Required) - Specifies the name of the cluster that contains the Managed Servers
for the service instance.

* `scaling_unit_count` - (Required) The number of scaling units that will be managed.

* `scaling_unit` - (Optional) Groups attributes for a custom capacity unit. Scaling Unit is specified below


Database supports the following:

* `username` - (Required) Username for the database administrator.

* `password` - (Required) Password for the database administrator.

* `name` - (Required) Name of the database on the Database Cloud Service.

* `network` - (Optional) Path to the network through which the Oracle Java Cloud Service instance will access the database.

* `uri` - (Exported) URI for the database on the Database Service Cloud.

Admin supports the following:

* `username` - (Required) Username for the WebLogic Server or Oracle Traffic Director administrator.

* `password` - (Required) Password for the WebLogic Server or Oracle Traffic Director administrator.

* `port` - (Optional) Port for accessing the WebLogic Server or Oracle Traffic Director using HTTP. The default values are 7001 for WebLogic Server or 8989 for Oracle Traffic Director.

App DB supports the following:

* `username` - (Required) Username for the database administrator.

* `password` - (Required) Password for the database administrator.

* `name` - (Required) Name of the database deployment on the Database Cloud Service.

* `pdb_name` - (Optional) Name of the pluggable database for Oracle Database 12c. If not specified, the pluggable database name configured when the database was created will be used.

Domain supports the following:

* `mode` - (Optional) Mode of the domain. Valid values are `DEVELOPMENT`  or `PRODUCTION`. Default value is
`PRODUCTION`.

* `name` - (Optional) Name of the WebLogic domain. By default, the domain name will be generated from the first
eight characters of the Oracle Java Cloud Service instance name (serviceName), using the
following format: first8charsOfServiceInstanceName_domain.

* `partition_count` - (Optional) Number of partitions to enable in the domain for WebLogic Server 12.2.1.
Valid values include: 0 (no partitions), 1, 2, and 4.

* `volume_size` - (Optional) Size of the domain volume for the service. The value must be a multiple of GBs.
You can specify this value in bytes or GBs. If specified in GBs, use the following format:
nG, where n specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs.
For example: 100000000000 or 10G.

Listener supports the following:

* `port` - (Optional) Listener port for the load balancer for accessing deployed applications using HTTP.
The default value is 8080.

* `type` - (Optional) Protocol used for the load balancer listener port. The default value is http.

* `enabled` - (Optional) Boolean on whether to enable the listener port. Default is true.

Managed Server supports the following:

* `server_count` - (Optional) Number of Managed Servers in the domain. Valid values include: 1, 2, 4, and 8.
The default value is 1.

* `initial_heap_size` - (Optional) Initial Java heap size for a Managed Server JVM, specified in megabytes.

* `max_heap_size` - (Optional) Maximum Java heap size for a Managed Server JVM, specified in megabytes.

* `jvm_args` - (Optional) One or more Managed Server JVM arguments separated by a space.

* `initial_permanent_generation` - (Optional) Initial Permanent Generation space in Java heap memory.

* `max_permanent_generation` - (Optional) Maximum Permanent Generation space in Java heap memory.

* `overwrite_jvm_args` - (Optional) Flag that determines whether the user defined Managed Server JVM arguments
specified in msJvmArgs should replace the server start arguments (true), or append the server start arguments
(false). Default is false.

Node Manager supports the following:

* `username` - (Optional) User name for the Node Manager. This value defaults to the WebLogic administrator user name.

* `password` - (Optional) Password for the Node Manager. This value defaults to the WebLogic administrator password.

* `port` - (Optional) Port for the Node Manager. This value defaults to 5556.

Privileged Port supports the following:

* `content_port` - (Optional) Privileged content port for accessing the deployed applications using HTTP.
To disable the privileged content port, set the value to 0. The default value is 80.

* `listener_port` - (Optional) Privileged content port for accessing the deployed applications using HTTPS.
To disable the privileged listener port, set the value to 0. The default value is 443.

* `secured_content_port` - (Optional) Privileged content port for accessing the deployed applications using HTTPS.
To disable the privileged secured content port, set the value to 0. The default value is 443.

* `secured_listener_port` - (Optional) Privileged listener port for accessing the deployed applications using HTTPS.
To disable the privileged listener port, set the value to 0. The default value is 443.

Scaling Unit supports the following:

* `name` - (Required) The name of the scaling unit. Valid values include `BASIC`, `SMALL`, `MEDIUM`, `LARGE`.

* `unit` - (Required) Specific details about each unit. Unit is documented below.

Unit supports the following:

* `heap_size` (Required) Heap size to configure with each JVM, based on the memory available from the chosen compute shape. Value must be between 1 and 16.

* `jvm_count` - (Required) Number of JVMs to start on each VM. Value must be between 1 and 8.

* `shape` - (Required) Desired compute shape.

* `vm_count` - (Required) Number of VMs to configure for a custom capacity unit. Value must be between 1 and 3.

Secured Ports supports the following:

* `admin_port` -  (Optional) Port for accessing the Administration Server using HTTPS. The default value is 7002.

* `content_port` - (Optional) Port for accessing the Administration Server using HTTPS. The default value is 8002.



In addition to the above, the following values are exported:

* `auto_update` - Flag that specifies whether updates to the Oracle Cloud Tools are automatically applied to the
Oracle Java Cloud Service instance during the maintenance window.

* `compliance_status` - Status indicating whether the version of Oracle Cloud Tools is out of compliance.

* `compliance_status_description` - Description that provides more details about the compliance status of the
Oracle Cloud Tools, used to manage the lifecycle of your Oracle Java Cloud Service instance.

* `created_by` - The user name of the Oracle Cloud user who created the service instance.

* `creation_time` - The date-and-time stamp when the service instance was created.

* `db_info` - Database that is used to host the Oracle Required Schema.

* `fmw_control_url` - URL to Enterprise Manager Fusion Middleware Control.

* `identity_domain` - Identity domain ID for the Oracle Java Cloud Service account (on Oracle Public Cloud).

* `is_app_2_cloud` - Flag that specifies whether this service instance is created with AppToCloud artifacts.

* `life_cycle_control_job_id` - Job ID of a lifecycle control request.

* `memory_size` - Total amount of memory in GBs allocated across all nodes in the service instance.

* `ocpu_count` - Total number of Oracle Compute Units (OCPUs) allocated across all nodes in the service instance.

* `otd_admin_url` - URL to load balancer Administration Console.

* `otd_provisioned` - Flag that specifies whether the load balancer is enabled.

* `otd_shape` - Desired compute shape for the load balancer.

* `otd_storage_size` - Storage size of the load balancer in GBs.

* `psm_plugin_version` - Version of the PaaS Service Manager.

* `sample_app_url` - URL for accessing the sample application, if it was installed and deployed when the service
instance was provisioned.

* `secure_content_url` - URL for accessing the deployed applications using HTTPS.

* `storage_size` - Total amount of block storage in GBs allocated across all nodes in the service instance.

* `wls_admin_url` - URL to the WebLogic Administration Console.

* `wls_deployment_channel_port` - Port for accessing the Administration Server using WLST.

* `wls_version` - Oracle WebLogic Server software version.

* `oracle_middleware_version` - Oracle Fusion Middleware software version.

* `uri` - The Uniform Resource Identifier for the Service Instance

## Import

Service Instance's can be imported using the `resource name`, e.g.

```shell
$ terraform import opc_java_service_instance.instance1 example
```

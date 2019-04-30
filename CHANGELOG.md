## 1.3.6 (Unreleased)
## 1.3.5 (April 30, 2019)

NOTES:

This release includes a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#168](https://github.com/terraform-providers/terraform-provider-opc/issues/168))

## 1.3.4 (April 23, 2019)

IMPROVEMENTS:

* Add ability to use different username/identity domains without forcing a recreate of the resource ([#167](https://github.com/terraform-providers/terraform-provider-opc/issues/167))

## 1.3.3 (April 04, 2019)

IMPROVEMENTS:

* Add ability to set `instance_attributes` for orchestrated instances ([#165](https://github.com/terraform-providers/terraform-provider-opc/issues/165))

## 1.3.2 (January 16, 2019)

IMPROVEMENTS:

* Add timeout support to snapshot ([#160](https://github.com/terraform-providers/terraform-provider-opc/issues/160))

## 1.3.1 (December 12, 2018)

IMPROVEMENTS:

* r/opc_vpn_endpoint_v2: Many attributes are now updatable/more robust ([#158](https://github.com/terraform-providers/terraform-provider-opc/issues/158))

## 1.3.0 (October 11, 2018)

FEATURES:

* **New Resource:** `r/opc_compute_vpn_endpoint_v2` ([#157](https://github.com/terraform-providers/terraform-provider-opc/issues/157))

## 1.2.1 (August 10, 2018)

IMPROVEMENTS:

* r/opc_lbaas_certificate: Now supports `TRUSTED` type ([#147](https://github.com/terraform-providers/terraform-provider-opc/issues/147))

* r/opc_ip_address_reservation: Adding support for custom `ip_address_pools` ([#152](https://github.com/terraform-providers/terraform-provider-opc/issues/152))

BUG FIXES: 

* r/opc_lbaas_server_pool: additional validation around origin servers ([#143](https://github.com/terraform-providers/terraform-provider-opc/issues/143))

* r/opc_lbaas_certificate: required attribute `certificate_chain` -> to optional ([#143](https://github.com/terraform-providers/terraform-provider-opc/issues/143))

## 1.2.0 (July 18, 2018)

FEATURES:

* **New Resource:** `r/opc_lbaas_certificate` [[#132](https://github.com/terraform-providers/terraform-provider-opc/issues/132)]    
* **New Resource:** `r/opc_lbaas_listener` [[#132](https://github.com/terraform-providers/terraform-provider-opc/issues/132)]          
* **New Resource:** `r/opc_lbaas_load_balancer` [[#132](https://github.com/terraform-providers/terraform-provider-opc/issues/132)]            
* **New Resource:** `r/opc_lbaas_policy` [[#132](https://github.com/terraform-providers/terraform-provider-opc/issues/132)]                 
* **New Resource:** `r/opc_lbaas_server_pool` ([#132](https://github.com/terraform-providers/terraform-provider-opc/issues/132))

BUG FIXES: 

* All resources confirm that the client has been properly initialized ([#136](https://github.com/terraform-providers/terraform-provider-opc/issues/136))

* Don't unqualify security lists in orchestrated instances ([#140](https://github.com/terraform-providers/terraform-provider-opc/issues/140))

## 1.1.2 (June 20, 2018)

FEATURES:

* **New Datasource:** `d/opc_compute_ssh_key` ([#129](https://github.com/terraform-providers/terraform-provider-opc/issues/129))
* **New Datasource:** `d/opc_compute_ip_reservation` ([#130](https://github.com/terraform-providers/terraform-provider-opc/issues/130))
* **New Datasource:** `d/opc_compute_ip_address_reservation ` ([#130](https://github.com/terraform-providers/terraform-provider-opc/issues/130))

## 1.1.1 (May 24, 2018)

IMPROVEMENTS:

* r/opc_compute_ip_reservation: Allows users to specify specific `parent_pool` ([#127](https://github.com/terraform-providers/terraform-provider-opc/issues/127))

## 1.1.0 (January 18, 2018)

FEATURES: 

* **New Resource:** `r/opc_storage_volume_attachment` ([#112](https://github.com/terraform-providers/terraform-provider-opc/issues/112))

* **New Resource:** `r/opc_compute_machine_image` ([#109](https://github.com/terraform-providers/terraform-provider-opc/issues/109))

BUG FIXES:

* r/opc_orchestrated_instance: Fixed silent failures when creating an orchestrated instance ([#108](https://github.com/terraform-providers/terraform-provider-opc/issues/108))

IMPROVEMENTS:

* r/opc_storage_container: Add `quota_bytes`, `quota_count` and `metadata` attributes ([#96](https://github.com/terraform-providers/terraform-provider-opc/issues/96))

* r/opc_storage_object: Add `metadata` attributes ([#96](https://github.com/terraform-providers/terraform-provider-opc/issues/96))

## 1.0.1 (December 20, 2017)

BUG FIXES

* Fixing broken link ([#103](https://github.com/terraform-providers/terraform-provider-opc/issues/103))

## 1.0.0 (December 20, 2017)

NEW RESOURCE:

* r/opc_compute_orchestrated_instance [#92]

IMPROVEMENTS:

* r/opc_compute_instance: Add `is_default_gatway` variable [#87]

* d/network_interface: Add `is_default_gateway` variable [#87]

* provider: Add `storage_service_id` variable to provider [#99]

* provider: Add useragent [#97]

* docs: Various documentation improvements [#95]

## 0.1.3 (September 15, 2017)

FEATURES:

* **New Resource:** `opc_storage_object` ([#55](https://github.com/terraform-providers/terraform-provider-opc/issues/55))

BUG FIXES:

* r/ip_network: Allow changing the name of an IP Network ([#73](https://github.com/terraform-providers/terraform-provider-opc/issues/73))
* r/opc_compute_image_list_entry: Fix resource imports ([#66](https://github.com/terraform-providers/terraform-provider-opc/issues/66))
* r/storage_container: Fixed `allowed_origins` ([#62](https://github.com/terraform-providers/terraform-provider-opc/issues/62))
* r/storage_volume: Allow errors to surface from Create ([#69](https://github.com/terraform-providers/terraform-provider-opc/issues/69))
* r/vnic_set: Make `virtual_nics` Computed [#52](https://github.com/terraform-providers/terraform-provider-opc/issues/52)
* r/storage_volume_snapshot: Increase timeout for larger snapshots ([#79](https://github.com/terraform-providers/terraform-provider-opc/issues/79))
* r/storage_volume: Remove validation around storage_type ([#80](https://github.com/terraform-providers/terraform-provider-opc/issues/80))

NOTES:

* Various doc fixes/updates

## 0.1.2 (August 02, 2017)

FEATURES:

 * **New Datasource:** `opc_compute_image_list_entry` ([#50](https://github.com/terraform-providers/terraform-provider-opc/issues/50))
 * **New Datasource:** `opc_compute_storage_volume_snapshot` ([#46](https://github.com/terraform-providers/terraform-provider-opc/issues/46))
 * **New Resource:** `opc_compute_storage_container` ([#23](https://github.com/terraform-providers/terraform-provider-opc/issues/23))
 * Add timeout configuration: ([#41](https://github.com/terraform-providers/terraform-provider-opc/issues/41))

BUG FIXES:
 * `opc_storage_volume_snapshot`: Fix crash on import ([#10](https://github.com/terraform-providers/terraform-provider-opc/issues/10))
 * `opc_compute_storage_volume`: bootable volumes can be added without an image list/image list entry ([#19](https://github.com/terraform-providers/terraform-provider-opc/issues/19))
 * r/security_list: Suppress case diffs for security_list ([#27](https://github.com/terraform-providers/terraform-provider-opc/issues/27))

## 0.1.1 (June 21, 2017)

NOTES:

* Bumping the provider version to get around provider caching issues - still same functionality

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)

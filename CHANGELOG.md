## 0.1.2 (Unreleased)

FEATURES:

 * **New Datasource:** `opc_compute_image_list_entry` [GH-50]
 * **New Datasource:** `opc_compute_storage_volume_snapshot` [GH-46]
 * **New Resource:** `opc_compute_storage_container` [GH-23]
 * Add timeout configuration: [GH-41]
  
BUG FIXES:
 * `opc_storage_volume_snapshot`: Fix crash on import [GH-10]
 * `opc_compute_storage_volume`: bootable volumes can be added without an image list/image list entry [GH-19]
 * r/security_list: Suppress case diffs for security_list [GH-27]

## 0.1.1 (June 21, 2017)

NOTES:

* Bumping the provider version to get around provider caching issues - still same functionality 

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)

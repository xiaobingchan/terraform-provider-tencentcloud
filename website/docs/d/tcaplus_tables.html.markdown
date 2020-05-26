---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_tables"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_tables"
description: |-
  Use this data source to query tcaplus tables
---

# tencentcloud_tcaplus_tables

Use this data source to query tcaplus tables

## Example Usage

```hcl
data "tencentcloud_tcaplus_tables" "null" {
  cluster_id = "19162256624"
}

data "tencentcloud_tcaplus_tables" "group" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
}

data "tencentcloud_tcaplus_tables" "name" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
  table_name = "guagua"
}

data "tencentcloud_tcaplus_tables" "id" {
  cluster_id = "19162256624"
  table_id   = "tcaplus-faa65eb7"
}
data "tencentcloud_tcaplus_tables" "all" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
  table_id   = "tcaplus-faa65eb7"
  table_name = "guagua"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Id of the tcaplus cluster to be query.
* `group_id` - (Optional) Group id to be query.
* `result_output_file` - (Optional) Used to save results.
* `table_id` - (Optional) Table id to be query.
* `table_name` - (Optional) Table name to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus groups. Each element contains the following attributes.
  * `create_time` - Create time of the tcaplus table.
  * `description` - Description of this table.
  * `error` - Show if this table  create error.
  * `group_id` - Group of this table belongs.
  * `idl_id` - Idl file id for this table.
  * `reserved_read_qps` - Table reserved read QPS.
  * `reserved_volume` - Table reserved capacity(GB).
  * `reserved_write_qps` - Table reserved write QPS.
  * `status` - Status of this table.
  * `table_id` - Id of this table.
  * `table_idl_type` - Type of this table idl.
  * `table_name` - Name of this table.
  * `table_size` - Size of this table.
  * `table_type` - Type of this table.



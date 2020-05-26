---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateways"
sidebar_current: "docs-tencentcloud-datasource-vpn_gateways"
description: |-
  Use this data source to query detailed information of VPN gateways.
---

# tencentcloud_vpn_gateways

Use this data source to query detailed information of VPN gateways.

## Example Usage

```hcl
data "tencentcloud_vpn_gateways" "foo" {
  name              = "main"
  id                = "vpngw-8ccsnclt"
  public_ip_address = "1.1.1.1"
  zone              = "ap-guangzhou-3"
  vpc_id            = "vpc-dk8zmwuf"
  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the VPN gateway.
* `name` - (Optional) Name of the VPN gateway. The length of character is limited to 1-60.
* `public_ip_address` - (Optional) Public ip address of the VPN gateway.
* `result_output_file` - (Optional) Used to save results.
* `tags` - (Optional) Tags of the VPN gateway to be queried.
* `vpc_id` - (Optional) ID of the VPC.
* `zone` - (Optional) Zone of the VPN gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_list` - Information list of the dedicated gateways.
  * `bandwidth` - The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100. Default is 5.
  * `charge_type` - Charge Type of the VPN gateway, valid values are `PREPAID`, `POSTPAID_BY_HOUR` and default is `POSTPAID_BY_HOUR`.
  * `create_time` - Create time of the VPN gateway.
  * `expired_time` - Expired time of the VPN gateway when charge type is `PREPAID`.
  * `id` - ID of the VPN gateway.
  * `is_address_blocked` - Indicates whether ip address is blocked.
  * `name` - Name of the VPN gateway.
  * `new_purchase_plan` - The plan of new purchase, valid value is `PREPAID_TO_POSTPAID`.
  * `prepaid_renew_flag` - Flag indicates whether to renew or not, valid values are `NOTIFY_AND_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `NOT_NOTIFY_AND_NOT_RENEW`.
  * `public_ip_address` - Public ip of the VPN gateway.
  * `restrict_state` - Restrict state of VPN gateway, valid values are `PRETECIVELY_ISOLATED`, `NORMAL`.
  * `state` - State of the VPN gateway, valid values are `PENDING`, `DELETING`, `AVAILABLE`.
  * `tags` - A list of tags used to associate different resources.
  * `type` - Type of gateway instance, valid values are `IPSEC`, `SSL` and `CCN`.
  * `vpc_id` - ID of the VPC.
  * `zone` - Zone of the VPN gateway.



/*
Provides a resource to create an entry of a routing table.

Example Usage

```hcl
variable "availability_zone" {
  default = "na-siliconvalley-1"
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  vpc_id            = tencentcloud_vpc.foo.id
  name              = "terraform test subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = var.availability_zone
  route_table_id    = tencentcloud_route_table.foo.id
}

resource "tencentcloud_route_table" "foo" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"
}

resource "tencentcloud_route_table_entry" "instance" {
  route_table_id         = tencentcloud_route_table.foo.id
  destination_cidr_block = "10.4.4.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = "ci-test-route-table-entry"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/errors"
)

func resourceTencentCloudVpcRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRouteEntryCreate,
		Read:   resourceTencentCloudVpcRouteEntryRead,
		Delete: resourceTencentCloudVpcRouteEntryDelete,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of routing table to which this entry belongs.",
			},
			"destination_cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "Destination address block.",
			},
			"next_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(ALL_GATE_WAY_TYPES),
				Description:  "Type of next-hop, and available values include CVM, VPN, DIRECTCONNECT, PEERCONNECTION, SSLVPN, NAT, NORMAL_CVM, EIP and CCN.",
			},
			"next_hub": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of next-hop gateway. Note: when 'next_type' is EIP, GatewayId should be '0'.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the routing table entry.",
			},
		},
	}
}

func resourceTencentCloudVpcRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_entry.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		description          = ""
		routeTableId         = ""
		destinationCidrBlock = ""
		nextType             = ""
		nextHub              = ""
	)

	if temp, ok := d.GetOk("description"); ok {
		description = temp.(string)
	}
	if temp, ok := d.GetOk("route_table_id"); ok {
		routeTableId = temp.(string)
	}
	if temp, ok := d.GetOk("destination_cidr_block"); ok {
		destinationCidrBlock = temp.(string)
	}
	if temp, ok := d.GetOk("next_type"); ok {
		nextType = temp.(string)
	}
	if temp, ok := d.GetOk("next_hub"); ok {
		nextHub = temp.(string)
	}

	if routeTableId == "" || destinationCidrBlock == "" || nextType == "" || nextHub == "" {
		return fmt.Errorf("some needed fields is empty string")
	}

	if nextType == GATE_WAY_TYPE_EIP && nextHub != "0" {
		return fmt.Errorf("if next_type is %s, next_hub can only be \"0\" ", GATE_WAY_TYPE_EIP)
	}

	entryId, err := service.CreateRoutes(ctx, routeTableId, destinationCidrBlock, nextType, nextHub, description)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d.%s", entryId, routeTableId))

	return nil
}

func resourceTencentCloudVpcRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_entry.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), ".")
	if len(items) != 2 {
		return fmt.Errorf("entry id be destroyed, we can not get route table id")
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeRouteTable(ctx, items[1])
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		if has != 1 {
			e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return resource.NonRetryableError(e)
		}

		for _, v := range info.entryInfos {
			if fmt.Sprintf("%d", v.routeEntryId) == items[0] {
				_ = d.Set("description", v.description)
				_ = d.Set("route_table_id", items[1])
				_ = d.Set("destination_cidr_block", v.destinationCidr)
				_ = d.Set("next_type", v.nextType)
				_ = d.Set("next_hub", v.nextBub)
				return nil
			}
		}
		d.SetId("")
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudVpcRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_entry.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), ".")
	if len(items) != 2 {
		return fmt.Errorf("entry id be destroyed, we can not get route table id")
	}

	routeTableId := items[1]
	entryId, err := strconv.ParseUint(items[0], 10, 64)
	if err != nil {
		return fmt.Errorf("entry id be destroyed, we can not get route entry id")
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteRoutes(ctx, routeTableId, entryId); err != nil {
			if sdkErr, ok := err.(*errors.TceCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}

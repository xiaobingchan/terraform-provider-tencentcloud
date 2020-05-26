// +build tencentcloud

/*
Provides an available EIP for the user.

The EIP data source fetch proper EIP from user's EIP pool.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_eips.

Example Usage

```hcl
data "tencentcloud_eip" "my_eip" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
```
*/
package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentyun/tcecloud-sdk-go/tcecloud/vpc/v20170312"
)

var (
	errEIPNotFound = errors.New("eip not found")
)

func dataSourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.20.0. Please use 'tencentcloud_eips' instead.",
		Read:               dataSourceTencentCloudEipRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "One or more name/value pairs to filter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the filter, valid keys: `address-id`,`address-name`,`address-ip`.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Value of the filter.",
						},
					},
				},
			},
			"include_arrears": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the IP is arrears.",
			},
			"include_blocked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the IP is blocked.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An EIP id indicate the uniqueness of a certain EIP,  which can be used for instance binding or network interface binding.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An public IP address for the EIP.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EIP, there are several status like `BIND`, `UNBIND`, and `BIND_ENI`.",
			},
		},
	}
}

func dataSourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eip.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string][]string)
	filters, ok := d.GetOk("filter")
	if ok {
		for _, v := range filters.(*schema.Set).List() {
			vv := v.(map[string]interface{})
			name := vv["name"].(string)
			filter[name] = []string{}
			for _, vvv := range vv["values"].([]interface{}) {
				filter[name] = append(filter[name], vvv.(string))
			}
		}
	}

	var eips []*vpc.Address
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eips, errRet = vpcService.DescribeEipByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	includeArrears := false
	if v, ok := d.GetOk("include_arrears"); ok {
		includeArrears = v.(bool)
	}
	includeBlocked := false
	if v, ok := d.GetOk("include_blocked"); ok {
		includeBlocked = v.(bool)
	}

	if len(eips) == 0 {
		return errEIPNotFound
	}

	var filteredEips []*vpc.Address
	for _, eip := range eips {
		if *eip.IsArrears && !includeArrears {
			continue
		}
		if *eip.IsBlocked && !includeBlocked {
			continue
		}
		filteredEips = append(filteredEips, eip)
	}

	if len(filteredEips) == 0 {
		return errEIPNotFound
	}

	eip := filteredEips[0]
	d.SetId(*eip.AddressId)
	_ = d.Set("public_ip", *eip.AddressIp)
	_ = d.Set("status", *eip.AddressStatus)

	return nil
}

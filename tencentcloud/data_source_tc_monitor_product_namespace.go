/*
Use this data source to query product namespace in monitor)

Example Usage

```hcl
data "tencentcloud_monitor_product_namespace" "instances" {
  name = "Redis"
}
```

*/
package tencentcloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentyun/tcecloud-sdk-go/tcecloud/monitor/v20180724"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentMonitorProductNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorProductNamespaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name for filter, eg:`Load Banlancer`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list product namespaces. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English name of this product.",
						},
						"product_chinese_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese name of this product.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace of each cloud product in monitor system.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorProductNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_product_namespace.read")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribeProductListRequest()
		response       *monitor.DescribeProductListResponse
		products       []*monitor.ProductSimple
		offset         uint64 = 0
		limit          uint64 = 20
		err            error
		filterName     = d.Get("name").(string)
	)

	request.Offset = &offset
	request.Limit = &limit
	request.Module = helper.String("monitor")

	var finish = false
	for {

		if finish {
			break
		}

		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribeProductList(request); err != nil {
				return retryError(err, InternalError)
			}
			products = append(products, response.Response.ProductList...)
			if len(response.Response.ProductList) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}
		offset = offset + limit
	}

	var list = make([]interface{}, 0, len(products))

	for _, product := range products {
		var listItem = map[string]interface{}{}
		listItem["product_name"] = product.ProductEnName
		listItem["product_chinese_name"] = product.ProductName
		listItem["namespace"] = product.Namespace
		if filterName == "" {
			list = append(list, listItem)
			continue
		}
		if product.ProductEnName != nil && strings.Contains(*product.ProductEnName, filterName) {
			list = append(list, listItem)
			continue
		}
		if product.ProductName != nil && strings.Contains(*product.ProductName, filterName) {
			list = append(list, listItem)
			continue
		}

	}
	if err = d.Set("list", list); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("product_namespace_%s", filterName))
	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}
	return nil
}

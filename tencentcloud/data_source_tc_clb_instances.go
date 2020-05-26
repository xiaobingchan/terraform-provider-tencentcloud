// +build tencentcloud

/*
Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "myclb"
  project_id         = 0
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentyun/tcecloud-sdk-go/tcecloud/clb/v20180317"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CLB to be queried.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'.",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CLB to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project id of the CLB.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"clb_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of CLB.",
						},
						"clb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CLB.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CLB.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id of the project.",
						},
						"clb_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The virtual service address table of the CLB.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of CLB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the CLB.",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest state transition time of CLB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the subnet.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Id set of the security groups.",
						},
						"target_region_info_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region information of backend service are attached the CLB.",
						},
						"target_region_info_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId information of backend service are attached the CLB.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CLB.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["project_id"] = v.(int)
	}
	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var clbs []*clb.LoadBalancer
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		clbs = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB instances failed, reason:%+v", logId, err)
		return err
	}
	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clbInstance := range clbs {
		mapping := map[string]interface{}{
			"clb_id":                    *clbInstance.LoadBalancerId,
			"clb_name":                  *clbInstance.LoadBalancerName,
			"network_type":              *clbInstance.LoadBalancerType,
			"status":                    *clbInstance.Status,
			"create_time":               *clbInstance.CreateTime,
			"status_time":               *clbInstance.StatusTime,
			"project_id":                *clbInstance.ProjectId,
			"vpc_id":                    *clbInstance.VpcId,
			"subnet_id":                 *clbInstance.SubnetId,
			"clb_vips":                  helper.StringsInterfaces(clbInstance.LoadBalancerVips),
			"target_region_info_region": *clbInstance.TargetRegionInfo.Region,
			"target_region_info_vpc_id": *clbInstance.TargetRegionInfo.VpcId,
			"security_groups":           helper.StringsInterfaces(clbInstance.SecureGroups),
		}
		if clbInstance.Tags != nil {
			tags := make(map[string]interface{}, len(clbInstance.Tags))
			for _, t := range clbInstance.Tags {
				tags[*t.TagKey] = *t.TagValue
			}
			mapping["tags"] = tags
		}
		clbList = append(clbList, mapping)
		ids = append(ids, *clbInstance.LoadBalancerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("clb_list", clbList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), clbList); e != nil {
			return e
		}
	}

	return nil
}

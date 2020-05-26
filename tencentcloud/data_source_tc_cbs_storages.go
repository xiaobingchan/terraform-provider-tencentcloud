/*
Use this data source to query detailed information of CBS storages.

Example Usage

```hcl
data "tencentcloud_cbs_storages" "storages" {
  storage_id         = "disk-kdt0sq6m"
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCbsStorages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsStoragesRead,

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CBS to be queried.",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CBS to be queried.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project with which the CBS is associated.",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_TYPE),
				Description:  "Types of storage medium, and available values include CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.",
			},
			"storage_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Types of CBS, and available values include SYSTEM_DISK and DATA_DISK.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"storage_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of storage. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CBS.",
						},
						"storage_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CBS.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of storage medium.",
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CBS.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone of CBS.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project.",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume of CBS.",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the CBS is mounted the CVM.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CVM instance that be mounted by this CBS.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether CBS is encrypted.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of CBS.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of CBS.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CBS.",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The way that CBS instance will be renew automatically or not when it reach the end of the prepaid tenancy.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of the CBS instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsStoragesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cbs_storages.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]string)
	if v, ok := d.GetOk("storage_id"); ok {
		params["disk-id"] = v.(string)
	}
	if v, ok := d.GetOk("storage_name"); ok {
		params["disk-name"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["project-id"] = fmt.Sprintf("%d", v.(int))
	}
	if v, ok := d.GetOk("storage_type"); ok {
		params["disk-type"] = v.(string)
	}
	if v, ok := d.GetOk("storage_usage"); ok {
		params["disk-usage"] = v.(string)
	}

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		storages, e := cbsService.DescribeDisksByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		ids := make([]string, 0, len(storages))
		storageList := make([]map[string]interface{}, 0, len(storages))
		for _, storage := range storages {
			mapping := map[string]interface{}{
				"storage_id":         storage.DiskId,
				"storage_name":       storage.DiskName,
				"storage_usage":      storage.DiskUsage,
				"storage_type":       storage.DiskType,
				"availability_zone":  storage.Placement.Zone,
				"project_id":         storage.Placement.ProjectId,
				"storage_size":       storage.DiskSize,
				"attached":           storage.Attached,
				"instance_id":        storage.InstanceId,
				"encrypt":            storage.Encrypt,
				"create_time":        storage.CreateTime,
				"status":             storage.DiskState,
				"prepaid_renew_flag": storage.RenewFlag,
				"charge_type":        storage.DiskChargeType,
			}
			if storage.Tags != nil {
				tags := make(map[string]interface{}, len(storage.Tags))
				for _, t := range storage.Tags {
					tags[*t.Key] = *t.Value
				}
				mapping["tags"] = tags
			}
			storageList = append(storageList, mapping)
			ids = append(ids, *storage.DiskId)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		if e = d.Set("storage_list", storageList); e != nil {
			log.Printf("[CRITAL]%s provider set storage list fail, reason:%s\n ", logId, e.Error())
			return resource.NonRetryableError(e)
		}

		output, ok := d.GetOk("result_output_file")
		if ok && output.(string) != "" {
			if e := writeToFile(output.(string), storageList); e != nil {
				return resource.NonRetryableError(e)
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs storages failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

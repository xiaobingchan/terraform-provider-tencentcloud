/*
Provides a resource to create a CBS.

Example Usage

```hcl
resource "tencentcloud_cbs_storage" "storage" {
  storage_name      = "mystorage"
  storage_type      = "CLOUD_SSD"
  storage_size      = "50"
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    test = "tf"
  }
}
```

Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.storage disk-41s6jwy4
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cbs/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCbsStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageCreate,
		Read:   resourceTencentCloudCbsStorageRead,
		Update: resourceTencentCloudCbsStorageUpdate,
		Delete: resourceTencentCloudCbsStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_TYPE),
				Description:  "Type of CBS medium, and available values include CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.",
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(10, 16000),
				Description:  "Volume of CBS, and unit is GB. If storage type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 36),
				Description:  "The purchased usage period of CBS, and value range [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"storage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Name of CBS. The maximum length can not exceed 60 bytes.",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the snapshot. If specified, created the CBS by this snapshot.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the project to which the instance belongs.",
			},
			"encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether CBS is encrypted.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this CBS.",
			},

			// computed
			"storage_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of CBS, and available values include UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.",
			},
			"attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the CBS is mounted the CVM.",
			},
		},
	}
}

func resourceTencentCloudCbsStorageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cbs.NewCreateDisksRequest()
	request.DiskName = helper.String(d.Get("storage_name").(string))
	request.DiskType = helper.String(d.Get("storage_type").(string))
	request.DiskSize = helper.IntUint64(d.Get("storage_size").(int))
	request.Placement = &cbs.Placement{
		Zone: helper.String(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.Placement.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = helper.String(v.(string))
	}
	if _, ok := d.GetOk("encrypt"); ok {
		request.Encrypt = helper.String("ENCRYPT")
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		request.Tags = make([]*cbs.Tag, 0, len(tags))
		for key, value := range tags {
			tag := cbs.Tag{
				Key:   helper.String(key),
				Value: helper.String(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	request.DiskChargeType = helper.String("POSTPAID_BY_HOUR")

	storageId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e, InternalError)
		}

		if len(response.Response.DiskIdSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("storage id is nil"))
		}

		storageId = *response.Response.DiskIdSet[0]
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(storageId)

	// must wait for finishing creating disk
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return retryError(e, InternalError)
		}
		if storage == nil {
			return resource.RetryableError(fmt.Errorf("storage is still creating..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudCbsStorageRead(d, meta)
}

func resourceTencentCloudCbsStorageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var storage *cbs.Disk
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		storage, e = cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if storage == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("storage_type", storage.DiskType)
	_ = d.Set("storage_size", storage.DiskSize)
	_ = d.Set("availability_zone", storage.Placement.Zone)
	_ = d.Set("storage_name", storage.DiskName)
	_ = d.Set("project_id", storage.Placement.ProjectId)
	_ = d.Set("encrypt", storage.Encrypt)
	_ = d.Set("tags", flattenCbsTagsMapping(storage.Tags))
	_ = d.Set("storage_status", storage.DiskState)
	_ = d.Set("attached", storage.Attached)

	return nil
}

func resourceTencentCloudCbsStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)
	storageId := d.Id()
	storageName := ""
	projectId := -1
	changed := false

	if d.HasChange("storage_name") {
		changed = true
		storageName = d.Get("storage_name").(string)
	}

	if d.HasChange("project_id") {
		changed = true
		projectId = d.Get("project_id").(int)
	}

	if changed {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifyDiskAttributes(ctx, storageId, storageName, projectId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if d.HasChange("storage_name") {
			d.SetPartial("storage_name")
		}
		if d.HasChange("project_id") {
			d.SetPartial("project_id")
		}
	}

	if d.HasChange("storage_size") {
		oldInterface, newInterface := d.GetChange("storage_size")
		oldValue := oldInterface.(int)
		newValue := newInterface.(int)
		if oldValue > newValue {
			return fmt.Errorf("storage size must be greater than current storage size")
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ResizeDisk(ctx, storageId, newValue)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return retryError(e)
			}
			if storage != nil && *storage.DiskState == CBS_STORAGE_STATUS_EXPANDING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("storage_size")
	}

	if d.HasChange("snapshot_id") {
		snapshotId := d.Get("snapshot_id").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ApplySnapshot(ctx, storageId, snapshotId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return retryError(e)
			}
			if storage != nil && *storage.DiskState == CBS_STORAGE_STATUS_ROLLBACKING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("snapshot_id")
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cvm:%s:uin/:volume/%s", region, storageId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudCbsStorageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskById(ctx, storageId)
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

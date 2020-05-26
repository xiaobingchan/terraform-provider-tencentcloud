/*
Provide a resource to attach an existing  cvm to kubernetes cluster.

Example Usage

```hcl

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

variable "default_instance_type" {
  default = "SA1.LARGE8"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}


data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  cpu_core_count = 8
  memory_size    = 16
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-auto-test-1-1"
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = var.default_instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = tencentcloud_kubernetes_cluster.managed_cluster.id
  instance_id = tencentcloud_instance.foo.id
  password    = "Lo4wbdit"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/errors"
	cvm "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cvm/v20170312"
	tke "github.com/tencentyun/tcecloud-sdk-go/tcecloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudTkeClusterAttachment() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "ID of the cluster.",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the CVM instance, this cvm will reinstall the system.",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateAsConfigPassword,
			Description:  "Password to access, should be set if `key_ids` not set.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.",
		},

		//compute
		"security_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "A list of security group ids after attach to cluster.",
		},
	}

	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterAttachmentCreate,
		Read:   resourceTencentCloudTkeClusterAttachmentRead,
		Delete: resourceTencentCloudTkeClusterAttachmentDelete,
		Schema: schemaBody,
	}
}

func resourceTencentCloudTkeClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	/*tke has been deleted*/
	_, has, err := tkeService.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = tkeService.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return nil
	}
	if !has {
		d.SetId("")
		return nil
	}

	/*cvm has been deleted*/
	var instance *cvm.Instance
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, err = cvmService.DescribeInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	/*attachment has been  deleted*/
	_, workers, err := tkeService.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, clusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has = false
	for _, worker := range workers {
		if worker.InstanceId == instanceId {
			has = true
		}
	}

	if !has {
		d.SetId("")
		return nil
	}

	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	}
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	return nil
}

func resourceTencentCloudTkeClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	request := tke.NewAddExistedInstancesRequest()

	instanceId := helper.String(d.Get("instance_id").(string))
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceIds = []*string{instanceId}
	request.LoginSettings = &tke.LoginSettings{}

	var loginSettingsNumbers = 0

	if v, ok := d.GetOk("key_ids"); ok {
		request.LoginSettings.KeyIds = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		loginSettingsNumbers++
	}

	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
		loginSettingsNumbers++
	}

	if loginSettingsNumbers != 1 {
		return fmt.Errorf("parameters `key_ids` and `password` must set and only set one")
	}

	/*cvm has been  attached*/
	var err error
	_, workers, err := tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has := false
	for _, worker := range workers {
		if worker.InstanceId == *instanceId {
			has = true
		}
	}
	if has {
		return fmt.Errorf("instance %s has been attached to cluster %s,can not attach again", *instanceId, *request.ClusterId)
	}

	var response *tke.AddExistedInstancesResponse

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = tkeService.client.UseTkeClient().AddExistedInstances(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("add existed instance %s to cluster %s error,reason %v", *instanceId, *request.ClusterId, err)
	}
	var success = false
	for _, v := range response.Response.SuccInstanceIds {
		if *v == *instanceId {
			d.SetId(*instanceId + "_" + *request.ClusterId)
			success = true
		}
	}

	if !success {
		return fmt.Errorf("add existed instance %s to cluster %s error, instance not in success instanceIds", *instanceId, *request.ClusterId)
	}

	/*wait for cvm status*/
	if err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, *instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance %s status is %s, retry...", *instanceId, *instance.InstanceState))
	}); err != nil {
		return err
	}

	/*wait for tke init ok */
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
		if err != nil {
			return retryError(err, InternalError)
		}
		has := false
		for _, worker := range workers {
			if worker.InstanceId == *instanceId {
				has = true
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						*instanceId, *request.ClusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						*instanceId, worker.InstanceState))
				}

			}
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s not exist in tke instance list", *instanceId))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTkeClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	request := tke.NewDeleteClusterInstancesRequest()

	request.ClusterId = &clusterId
	request.InstanceIds = []*string{
		&instanceId,
	}
	request.InstanceDeleteMode = helper.String("retain")

	var err error

	if err = resource.Retry(4*writeRetryTimeout, func() *resource.RetryError {
		_, err := tkeService.client.UseTkeClient().DeleteClusterInstances(request)
		if e, ok := err.(*errors.TceCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
			if e.GetCode() == "InternalError.Param" &&
				strings.Contains(e.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
				return nil
			}
		}

		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

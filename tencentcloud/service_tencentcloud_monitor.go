// +build tencentcloud

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MonitorService struct {
	client *connectivity.TencentCloudClient
}

func (me *MonitorService) CheckCanCreateMysqlROInstance(ctx context.Context, mysqlId string) (can bool, errRet error) {

	logId := getLogId(ctx)

	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		errRet = fmt.Errorf("Can not load  time zone `Asia/Chongqing`, reason %s", err.Error())
		return
	}

	request := monitor.NewGetMonitorDataRequest()

	request.Namespace = helper.String("QCE/CDB")
	request.MetricName = helper.String("RealCapacity")
	request.Period = helper.Uint64(60)

	now := time.Now()
	request.StartTime = helper.String(now.Add(-5 * time.Minute).In(loc).Format("2006-01-02T15:04:05+08:00"))
	request.EndTime = helper.String(now.In(loc).Format("2006-01-02T15:04:05+08:00"))

	request.Instances = []*monitor.Instance{
		{
			Dimensions: []*monitor.Dimension{{
				Name:  helper.String("InstanceId"),
				Value: &mysqlId,
			}},
		},
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().GetMonitorData(request)
	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.DataPoints) == 0 {
		return
	}
	dataPoint := response.Response.DataPoints[0]
	if len(dataPoint.Values) == 0 {
		return
	}
	can = true
	return
}

func (me *MonitorService) FullRegions() (regions []string, errRet error) {
	request := cvm.NewDescribeRegionsRequest()
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err := me.client.UseCvmClient().DescribeRegions(request); err != nil {
			return retryError(err, InternalError)
		} else {
			for _, region := range response.Response.RegionSet {
				regions = append(regions, *region.Region)
			}
		}
		return nil
	}); err != nil {
		errRet = err
		return
	}
	return
}

func (me *MonitorService) DescribePolicyGroupDetailInfo(ctx context.Context, groupId int64) (response *monitor.DescribePolicyGroupInfoResponse, errRet error) {

	var (
		request = monitor.NewDescribePolicyGroupInfoRequest()
		err     error
	)
	request.GroupId = &groupId
	request.Module = helper.String("monitor")

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = me.client.UseMonitorClient().DescribePolicyGroupInfo(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		errRet = err
		return
	}
	return
}

func (me *MonitorService) DescribePolicyGroup(ctx context.Context, groupId int64) (info *monitor.DescribePolicyGroupListGroup, errRet error) {

	var (
		request       = monitor.NewDescribePolicyGroupListRequest()
		offset  int64 = 0
		limit   int64 = 20
		finish  bool
	)
	request.Module = helper.String("monitor")
	request.Offset = &offset
	request.Limit = &limit

	for {
		if finish {
			break
		}
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := me.client.UseMonitorClient().DescribePolicyGroupList(request)
			if err != nil {
				return retryError(err, InternalError)
			}
			if len(response.Response.GroupList) < int(limit) {
				finish = true
			}
			for _, v := range response.Response.GroupList {
				if *v.GroupId == groupId {
					info = v
					return nil
				}
			}
			return nil
		}); err != nil {
			errRet = err
			return
		}
		if info != nil {
			return
		}
		offset = offset + limit
	}
	return
}
func (me *MonitorService) DescribeBindingPolicyObjectList(ctx context.Context, groupId int64) (objects []*monitor.DescribeBindingPolicyObjectListInstance, errRet error) {

	var (
		requestList  = monitor.NewDescribeBindingPolicyObjectListRequest()
		responseList *monitor.DescribeBindingPolicyObjectListResponse
		offset       int64 = 0
		limit        int64 = 20
		finish       bool
		err          error
	)

	requestList.GroupId = &groupId
	requestList.Module = helper.String("monitor")
	requestList.Offset = &offset
	requestList.Limit = &limit

	for {
		if finish {
			break
		}
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(requestList.GetAction())
			if responseList, err = me.client.UseMonitorClient().DescribeBindingPolicyObjectList(requestList); err != nil {
				return retryError(err, InternalError)
			}
			objects = append(objects, responseList.Response.List...)
			if len(responseList.Response.List) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			errRet = err
			return
		}
		offset = offset + limit
	}

	return
}

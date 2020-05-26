data "tencentcloud_cam_groups" "groups" {
  //You should first create a user group with CAM
}
data "tencentcloud_instances" "instances" {
}

resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}

#for cvm
resource "tencentcloud_monitor_binding_object" "binding" {
  group_id = tencentcloud_monitor_policy_group.group.id
  dimensions {
    dimensions_json = "{\"unInstanceId\":\"${data.tencentcloud_instances.instances.instance_list[0].instance_id}\"}"
  }
}

resource "tencentcloud_monitor_binding_receiver" "receiver" {
  group_id = tencentcloud_monitor_policy_group.group.id
  receivers {
    start_time          = 0
    end_time            = 86399
    notify_way          = ["SMS"]
    receiver_type       = "group"
    receiver_group_list = [data.tencentcloud_cam_groups.groups.group_list[0].group_id]
    receive_language    = "en-US"
  }
}

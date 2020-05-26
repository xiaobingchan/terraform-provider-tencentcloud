package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMonitorPolicyGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorPolicyGroups(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_policy_groups.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_policy_groups.name",
						"list.#"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorPolicyGroups() string {
	return `data "tencentcloud_monitor_policy_groups" "name" {}`
}

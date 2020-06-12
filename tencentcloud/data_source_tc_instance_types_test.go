package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.availability_zone", "yf-1"),
				),
			},
		},
	})
}

const testAccTencentCloudInstanceTypesDataSourceConfigBasic = `
data "tencentcloud_instance_types" "t4c8g" {
  availability_zone = "yf-1"
  cpu_core_count = 4
  memory_size    = 8
}
`

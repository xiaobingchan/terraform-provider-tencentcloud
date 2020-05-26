package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMongodbInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "engine_version", "MONGO_36_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "machine_type", "GIO"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "available_zone", "ap-guangzhou-2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test", "test"),
				),
			},
			{
				Config: testAccMongodbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "200"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "password"},
			},
		},
	})
}

func testAccCheckMongodbInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mongodbService := MongodbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_instance" {
			continue
		}

		_, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("mongodb instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMongodbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mongodb instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance id is not set")
		}
		mongodbService := MongodbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccMongodbInstance = `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "GIO"
  available_zone = "ap-guangzhou-2"
  project_id     = 0
  password       = "test1234"

  tags = {
    "test" = "test"
  }
}
`

const testAccMongodbInstance_update = `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = 8
  volume         = 200
  engine_version = "MONGO_36_WT"
  machine_type   = "GIO"
  available_zone = "ap-guangzhou-2"
  project_id     = 0
  password       = "tests1234"

  tags = {
    "abc" = "abc"
  }
}
`

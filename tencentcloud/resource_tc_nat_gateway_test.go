package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	vpc "github.com/tencentyun/tcecloud-sdk-go/tcecloud/vpc/v20170312"
)

func TestAccTencentCloudNatGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "3000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "500"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "10000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_nat_gateway" {
			continue
		}
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 0 {
			return fmt.Errorf("nat gateway id is still exists")
		}

	}
	return nil
}

func testAccCheckNatGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("nat gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("nat gateway id is not set")
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 1 {
			return fmt.Errorf("nat gateway id is not found")
		}
		return nil
	}
}

const testAccNatGatewayConfig = `
data "tencentcloud_vpc_instances" "foo" {
	name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name             = "terraform_test"
  max_concurrent   = 3000000
  bandwidth        = 500

  assigned_eip_set = [
	  tencentcloud_eip.eip_dev_dnat.public_ip,
	  tencentcloud_eip.eip_test_dnat.public_ip,
	]
}
`
const testAccNatGatewayConfigUpdate = `
data "tencentcloud_vpc_instances" "foo" {
	name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
	name = "terraform_test"
  }
resource "tencentcloud_eip" "eip_test_dnat" {
	name = "terraform_test"
}
resource "tencentcloud_eip" "new_eip" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name             = "new_name"
  max_concurrent   = 10000000
  bandwidth        = 1000
  assigned_eip_set = [
	  tencentcloud_eip.eip_dev_dnat.public_ip,
	  tencentcloud_eip.new_eip.public_ip,
	]
}
`

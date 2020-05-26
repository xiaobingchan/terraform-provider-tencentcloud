package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/errors"
	vpc "github.com/tencentyun/tcecloud-sdk-go/tcecloud/vpc/v20170312"
)

func TestAccTencentCloudVpnCustomerGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("tencentcloud_vpn_customer_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "public_ip_address", "1.1.1.2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "tags.test", "tf"),
				),
			},
			{
				Config: testAccVpnCustomerGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("tencentcloud_vpn_customer_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "public_ip_address", "1.1.1.2"),
				),
			},
		},
	})
}

func testAccCheckVpnCustomerGatewayDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_customer_gateway" {
			continue
		}
		request := vpc.NewDescribeCustomerGatewaysRequest()
		request.CustomerGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeCustomerGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeCustomerGateways(request)
			if e != nil {
				ee, ok := e.(*errors.TceCloudSDKError)
				if !ok {
					return retryError(e)
				}
				if ee.Code == VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return retryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*errors.TceCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound" {
				return nil
			} else {
				return err
			}
		} else {
			if len(response.Response.CustomerGatewaySet) != 0 {
				return fmt.Errorf("VPN customer gateway id is still exists")
			}
		}

	}
	return nil
}

func testAccCheckVpnCustomerGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN customer gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN customer gateway id is not set")
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeCustomerGatewaysRequest()
		request.CustomerGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeCustomerGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeCustomerGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.CustomerGatewaySet) != 1 {
			return fmt.Errorf("VPN customer gateway id is not found")
		}
		return nil
	}
}

const testAccVpnCustomerGatewayConfig = `
resource "tencentcloud_vpn_customer_gateway" "my_cgw" {
  name              = "terraform_test"
  public_ip_address = "1.1.1.2" 

  tags = {
    test = "tf"
  }
}
`
const testAccVpnCustomerGatewayConfigUpdate = `
resource "tencentcloud_vpn_customer_gateway" "my_cgw" {
  name              = "terraform_update"
  public_ip_address = "1.1.1.2"

  tags = {
    test = "test"
  }
}
`

package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClbListener_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "listener_name", "listener_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "scheduler", "WRR"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbListener_tcp(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListener_tcp, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
				),
			},
			{
				Config: testAccClbListener_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbListener_https(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListener_https, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "77"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_id", defaultSshCertificate),
				),
			},
			{
				Config: fmt.Sprintf(testAccClbListener_https_update, defaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "33"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_id", defaultSshCertificateB),
				),
			},
		},
	})
}

func TestAccTencentCloudClbListener_tcpssl(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListener_tcpssl, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcpssl"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcpssl", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "protocol", "TCP_SSL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "listener_name", "listener_tcpssl"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_id", defaultSshCertificate),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_unhealth_num", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccClbListener_tcpssl_update, defaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcpssl"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcpssl", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "protocol", "TCP_SSL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "listener_name", "listener_tcpssl_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_id", defaultSshCertificateB),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_unhealth_num", "3"),
				),
			},
		},
	})
}

func testAccCheckClbListenerDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_listener" {
			continue
		}
		time.Sleep(5 * time.Second)
		clbId := rs.Primary.Attributes["clb_id"]
		instance, err := clbService.DescribeListenerById(ctx, rs.Primary.ID, clbId)
		if instance != nil && err == nil {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB listener][Destroy] check: CLB listener still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbListenerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB listener][Exists] check: CLB listener %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB listener][Exists] check: CLB listener id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		clbId := rs.Primary.Attributes["clb_id"]
		instance, err := clbService.DescribeListenerById(ctx, rs.Primary.ID, clbId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB listener][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbListener_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  protocol            = "TCP"
  listener_name       = "listener_basic"
  session_expire_time = 30
  scheduler           = "WRR"
}
`

const testAccClbListener_tcp = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp"
}

resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
}
`

const testAccClbListener_tcp_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp"
}

resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp_update"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 60
  scheduler                  = "WRR"
}
`

const testAccClbListener_tcpssl = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-tcpssl"
}

resource "tencentcloud_clb_listener" "listener_tcpssl" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcpssl"
  port                       = 44
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  scheduler                  = "WRR"
}
`
const testAccClbListener_tcpssl_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-tcpssl"
}

resource "tencentcloud_clb_listener" "listener_tcpssl"{
  clb_id = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcpssl_update"
  port                       = 44
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
}
`
const testAccClbListener_https = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-https"
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch           = true
}
`

const testAccClbListener_https_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-https"
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https_update"
  port                 = 33
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch           = true
}
`

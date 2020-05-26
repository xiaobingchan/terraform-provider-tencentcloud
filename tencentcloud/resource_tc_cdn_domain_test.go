package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCdnDomain(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "domain", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "service_type", "web"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "area", "mainland"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
				),
			},
			{
				Config: testAccCdnDomainFull,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "domain", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "service_type", "web"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "area", "mainland"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.server_name", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.https_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.http2_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.ocsp_stapling_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.spdy_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.verify_client", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.message", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.deploy_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.expire_time"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "tags.hello", "world"),
				),
			},
			{
				ResourceName:            "tencentcloud_cdn_domain.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"https_config"},
			},
		},
	})
}

func testAccCheckCdnDomainDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cdnService := CdnService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cdn_domain" {
			continue
		}

		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				domainConfig, err = cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if domainConfig != nil {
			return fmt.Errorf("cdn domain still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCdnDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cdn domain %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cdn domain id is not set")
		}
		cdnService := CdnService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				domainConfig, err = cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if domainConfig == nil {
			return fmt.Errorf("cdn domain is not found")
		}
		return nil
	}
}

const testAccCdnDomain = `
resource "tencentcloud_cdn_domain" "foo" {
  domain = "test.zhaoshaona.com"
  service_type = "web"
  area = "mainland"
  origin {
	origin_type = "ip"
	origin_list = ["139.199.199.140"]
  }
}
`

const testAccCdnDomainFull = `
resource "tencentcloud_cdn_domain" "foo" {
  domain = "test.zhaoshaona.com"
  service_type = "web"
  area = "mainland"
  full_url_cache = false

  origin {
	origin_type = "ip"
	origin_list = ["139.199.199.140"]
    server_name = "test.zhaoshaona.com"
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch = "on"
    http2_switch = "on"
    ocsp_stapling_switch = "on"
    spdy_switch = "on"
	verify_client = "off"

	server_certificate_config {
      certificate_content = <<EOT
-----BEGIN CERTIFICATE-----
MIIDuDCCAqACCQDJd98Shn/cJTANBgkqhkiG9w0BAQsFADCBnTELMAkGA1UEBhMC
Q04xEDAOBgNVBAgMB1RpYW5qaW4xEDAOBgNVBAcMB1RpYW5qaW4xDjAMBgNVBAoM
BU1vY2hhMRcwFQYDVQQLDA5Nb2NoYSBTb2Z0d2FyZTEcMBoGA1UEAwwTdGVzdC56
aGFvc2hhb25hLmNvbTEjMCEGCSqGSIb3DQEJARYUeWFsaW5wZWlAdGVuY2VudC5j
b20wHhcNMjAwNTIwMDcyNDQyWhcNMzAwNTE4MDcyNDQyWjCBnTELMAkGA1UEBhMC
Q04xEDAOBgNVBAgMB1RpYW5qaW4xEDAOBgNVBAcMB1RpYW5qaW4xDjAMBgNVBAoM
BU1vY2hhMRcwFQYDVQQLDA5Nb2NoYSBTb2Z0d2FyZTEcMBoGA1UEAwwTdGVzdC56
aGFvc2hhb25hLmNvbTEjMCEGCSqGSIb3DQEJARYUeWFsaW5wZWlAdGVuY2VudC5j
b20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgndm2xEWL7CaVQ/lb
TO6Gj4EqEp1tWygjdfqkUXADfsgMGPukYaZY+klV6AJzLcj8VD5iWgKa+V4kLHtf
yh66c45nZrdUVoF9CFTw2+B/LTa/UzsvbLTVOnEjVBjI1V5kVzliF5cK5OlQ258d
w6yFaccOgXqSkp9i57Y9pT1FIb691hsf2VHiVLizPYy3vvLQeN8RnXS3vK56BcQk
o+49H11TAsrIh0C5maF0jp/7poSQkrX0kjfX4+gK/mC4Dn3PgK464Ko5OR45IGji
D368/klCK1bqIshlv4owEfgzAEQMPUQ0CfuvXTX85aojM48RiYiDmYveaICtYnSR
04MTAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAHWUpfePVt3LjZVDS3OmQ7rTG8zc
zwZgJfxP0f4ZNo/9t53SNsQ0UM/+7sqnBKOjsWfgyFqSh9cfN0Bnsn3gmvPXmD5R
nCa9qr9IO+FP9Ke4KKn0Ndx1sWJN3B6D8bUTnvFIkJoRsvsqNi24o2uKrUdcAYHL
5BVtrVe8E55i0A5WosC8KWv4ZJxTacvuxVjfyroKzxsLwOQvCqBNSuZLg1HYUeG6
XIj0/acmysb8S82Lxm39E82DbPdUO3Z0TlGL7umlAV947/6eGvPhszjnhBlxVo3p
tmHdyqfHxWbkTW4bnO/Gu+Sll6a3n1uyQ/onXuXH3pBZoXLp3Jj+CV1+N6E=
-----END CERTIFICATE-----
EOT

      private_key         = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoJ3ZtsRFi+wmlUP5W0zuho+BKhKdbVsoI3X6pFFwA37IDBj7
pGGmWPpJVegCcy3I/FQ+YloCmvleJCx7X8oeunOOZ2a3VFaBfQhU8Nvgfy02v1M7
L2y01TpxI1QYyNVeZFc5YheXCuTpUNufHcOshWnHDoF6kpKfYue2PaU9RSG+vdYb
H9lR4lS4sz2Mt77y0HjfEZ10t7yuegXEJKPuPR9dUwLKyIdAuZmhdI6f+6aEkJK1
9JI31+PoCv5guA59z4CuOuCqOTkeOSBo4g9+vP5JQitW6iLIZb+KMBH4MwBEDD1E
NAn7r101/OWqIzOPEYmIg5mL3miArWJ0kdODEwIDAQABAoIBAQCW2uuLX9k6THkI
pSlleWJm/A4C6Cz4cy/F7p+VCfA9OCzIQAbKI/VLiPisORdj+tLVPILDeWsNB75G
F4lhNMObt8E+mRkDm6RPPS4ac0nt6ReMp63lIyLNSvDMj8Yfi1f2wn3hBesVjl8d
VMmj+Q7m16zgkPgBBrmw+ZUPXU2oyUW4+0RvGYvuWnVUdtm/34PD1LC0NKBKaX9T
MDHrSIns0WpQ7P4vNVQyHW7MGgEl81uzIitSWuT/k+zH6YxBlxd7d66vmhNoxz9c
aeEf7DE3wAb4819UYWt0/ciMJwSLPkBOaTeAsktKUHVsrMLVELWcWqSIS+PYbSX8
g3tY1DlxAoGBANSiDKNjfr1rZRtpZvxnVkssLY/586UaHs+dFfyFyd0unr/rAPf/
GO/BIO0NbBdRb3XORMuiLQN3xf+qgKfoS0kXYglDMGKbEAC/5o6ZMV6E2E/aFrxh
xmgKTZxCBVnOxlAy33UFs+qR8tpOnR4auAc0pNPA9QB4I7q17vGJRMyHAoGBAMFf
7nF2aJ/k0Fcl53Cabs/FIaAwL/GBvok6Ny8wWLwiMZCtsGUUywnUdN/qbfr2GwC5
g0w2iaxGqQPI+qw2qn0utAIfZ0Tz2VAH+P3aUTuG+M4XWHObHVXxBUqO61X9zgV2
sXRXcbDOx3HgZeDCjk0otcGVJoC3zgzaaEZi5mQVAoGAQer+2gQ1PUm27XmOmL78
bI+EjHbjhpKDbL95GnDrdKtIQZz8DuXBeEo6B+M6WDxBvpa0kyByrfmKo0jbW7JS
7JTYKqDuthL2MhVLx3dMa83pNVAZ7kqtdIGFL+TzvbSxnBk5VxDuhtC6Jd1rLfMA
jBNQ6eiOy5dzFCXkrnJspq8CgYAO4ISFsihmdMIakk31+cugrHfjzRFDMUopYJMy
TDPndXH+wX4aqLjeLrw3JeAEOL7nFV6mlGOPH3iNU/8FFMeVDezHZQca5O/JGnPr
g8pQHBg0MtOZQUvGet5/V/N/ECGzhegtHTUf9yic+DieTBmKkiE5nXHy4TE3B+6R
y7YR6QKBgQDUoNAFOnMZB4BQMeCb/pQQnzNkNTG+Y02eMKjo5eZZDfyusqIui29l
KKcVGqvwVh2r8ocP7OnrQPVK9ZW7BcoYiqM2DjdKyl7AtQKnvWfPMai++oXKzo0y
8sg7m1Ic26sKO9W9t87cfZtFKcbKVcImLWucd9R7Ny4M4r6xlRKWpA==
-----END RSA PRIVATE KEY-----
EOT
      message = "test"
    }
  }

  tags = {
    hello = "world"
  }
}
`

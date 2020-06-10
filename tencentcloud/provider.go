/*
The TencentCloud provider is used to interact with many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.9.0 (June 18, 2019), the provider start to support Terraform 0.12.x.

Example Usage

```hcl
# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
}

#Configure the TencentCloud Provider with STS
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
  assume_role {
    role_arn         = var.assume_role_arn
    session_name     = var.session_name
    session_duration = var.session_duration
    policy           = var.policy
  }
}
```

Resources List

Provider Data Sources
  tencentcloud_availability_regions
  tencentcloud_availability_zones

Anti-DDoS(Dayu)
  Data Source
    tencentcloud_dayu_cc_http_policies
    tencentcloud_dayu_cc_https_policies
    tencentcloud_dayu_ddos_policies
    tencentcloud_dayu_ddos_policy_attachments
    tencentcloud_dayu_ddos_policy_cases
    tencentcloud_dayu_l4_rules
    tencentcloud_dayu_l7_rules

  Resource
    tencentcloud_dayu_cc_http_policy
    tencentcloud_dayu_cc_https_policy
    tencentcloud_dayu_ddos_policy
    tencentcloud_dayu_ddos_policy_attachment
    tencentcloud_dayu_ddos_policy_case
    tencentcloud_dayu_l4_rule
    tencentcloud_dayu_l7_rule

Auto Scaling(AS)
  Data Source
    tencentcloud_as_scaling_configs
    tencentcloud_as_scaling_groups
    tencentcloud_as_scaling_policies

  Resource
    tencentcloud_as_scaling_config
    tencentcloud_as_scaling_group
    tencentcloud_as_attachment
    tencentcloud_as_scaling_policy
    tencentcloud_as_schedule
    tencentcloud_as_lifecycle_hook
    tencentcloud_as_notification

CDN
  Resource
	tencentcloud_cdn_domain

Cloud Access Management(CAM)
  Data Source
    tencentcloud_cam_group_memberships
    tencentcloud_cam_group_policy_attachments
    tencentcloud_cam_groups
    tencentcloud_cam_policies
    tencentcloud_cam_role_policy_attachments
    tencentcloud_cam_roles
    tencentcloud_cam_saml_providers
    tencentcloud_cam_user_policy_attachments
    tencentcloud_cam_users

  Resource
    tencentcloud_cam_role
    tencentcloud_cam_role_policy_attachment
    tencentcloud_cam_policy
    tencentcloud_cam_user
    tencentcloud_cam_user_policy_attachment
    tencentcloud_cam_group
    tencentcloud_cam_group_policy_attachment
    tencentcloud_cam_group_membership
    tencentcloud_cam_saml_provider

Cloud Block Storage(CBS)
  Data Source
    tencentcloud_cbs_snapshots
    tencentcloud_cbs_storages
    tencentcloud_cbs_snapshot_policies

  Resource
    tencentcloud_cbs_storage
    tencentcloud_cbs_storage_attachment
    tencentcloud_cbs_snapshot
    tencentcloud_cbs_snapshot_policy
    tencentcloud_cbs_snapshot_policy_attachment

Cloud Connect Network(CCN)
  Data Source
    tencentcloud_ccn_bandwidth_limits
    tencentcloud_ccn_instances

  Resource
    tencentcloud_ccn
    tencentcloud_ccn_attachment
    tencentcloud_ccn_bandwidth_limit

Cloud File Storage(CFS)
  Data Source
    tencentcloud_cfs_access_groups
    tencentcloud_cfs_access_rules
    tencentcloud_cfs_file_systems

  Resource
    tencentcloud_cfs_file_system
    tencentcloud_cfs_access_group
    tencentcloud_cfs_access_rule

Container Cluster
  Data Source
    tencentcloud_container_cluster_instances
    tencentcloud_container_clusters

  Resource
    tencentcloud_container_cluster
    tencentcloud_container_cluster_instance

CLB
  Data Source
    tencentcloud_clb_attachments
    tencentcloud_clb_instances
    tencentcloud_clb_listener_rules
    tencentcloud_clb_listeners
    tencentcloud_clb_redirections

  Resource
    tencentcloud_clb_instance
    tencentcloud_clb_listener
    tencentcloud_clb_listener_rule
    tencentcloud_clb_attachment
    tencentcloud_clb_redirection
    tencentcloud_lb
    tencentcloud_alb_server_attachment

COS
  Data Source
    tencentcloud_cos_bucket_object
    tencentcloud_cos_buckets

  Resource
    tencentcloud_cos_bucket
    tencentcloud_cos_bucket_object

CVM
  Data Source
    tencentcloud_image
    tencentcloud_images
    tencentcloud_instance_types
    tencentcloud_instances
    tencentcloud_key_pairs
    tencentcloud_eip
    tencentcloud_eips
    tencentcloud_placement_groups
    tencentcloud_reserved_instance_configs
    tencentcloud_reserved_instances

  Resource
    tencentcloud_instance
    tencentcloud_eip
    tencentcloud_eip_association
    tencentcloud_key_pair
    tencentcloud_placement_group
    tencentcloud_reserved_instance

Direct Connect(DC)
  Data Source
    tencentcloud_dc_instances
    tencentcloud_dcx_instances

  Resource
    tencentcloud_dcx

Direct Connect Gateway(DCG)
  Data Source
    tencentcloud_dc_gateway_ccn_routes
    tencentcloud_dc_gateway_instances

  Resource
    tencentcloud_dc_gateway
    tencentcloud_dc_gateway_ccn_route

Elasticsearch
  Data Source
    tencentcloud_elasticsearch_instances

  Resource
    tencentcloud_elasticsearch_instance

Global Application Acceleration(GAAP)
  Data Source
    tencentcloud_gaap_certificates
    tencentcloud_gaap_http_domains
    tencentcloud_gaap_http_rules
    tencentcloud_gaap_layer4_listeners
    tencentcloud_gaap_layer7_listeners
    tencentcloud_gaap_proxies
    tencentcloud_gaap_realservers
    tencentcloud_gaap_security_policies
    tencentcloud_gaap_security_rules
    tencentcloud_gaap_domain_error_pages

  Resource
    tencentcloud_gaap_proxy
    tencentcloud_gaap_realserver
    tencentcloud_gaap_layer4_listener
    tencentcloud_gaap_layer7_listener
    tencentcloud_gaap_http_domain
    tencentcloud_gaap_http_rule
    tencentcloud_gaap_certificate
    tencentcloud_gaap_security_policy
    tencentcloud_gaap_security_rule
    tencentcloud_gaap_domain_error_page

Kubernetes
  Data Source
    tencentcloud_kubernetes_clusters

  Resource
    tencentcloud_kubernetes_cluster
    tencentcloud_kubernetes_scale_worker
    tencentcloud_kubernetes_as_scaling_group
    tencentcloud_kubernetes_cluster_attachment

MongoDB
  Data Source
    tencentcloud_mongodb_instances
    tencentcloud_mongodb_zone_config

  Resource
    tencentcloud_mongodb_instance
    tencentcloud_mongodb_sharding_instance

MySQL
  Data Source
    tencentcloud_mysql_backup_list
    tencentcloud_mysql_instance
    tencentcloud_mysql_parameter_list
    tencentcloud_mysql_zone_config

  Resource
    tencentcloud_mysql_instance
    tencentcloud_mysql_readonly_instance
    tencentcloud_mysql_account
    tencentcloud_mysql_privilege
    tencentcloud_mysql_account_privilege
    tencentcloud_mysql_backup_policy

Monitor
  Data Source
	tencentcloud_monitor_policy_conditions
	tencentcloud_monitor_data
	tencentcloud_monitor_product_event
	tencentcloud_monitor_binding_objects
	tencentcloud_monitor_policy_groups
	tencentcloud_monitor_product_namespace

  Resource
    tencentcloud_monitor_policy_group
    tencentcloud_monitor_binding_object
    tencentcloud_monitor_binding_receiver

Redis
  Data Source
    tencentcloud_redis_zone_config
    tencentcloud_redis_instances

  Resource
    tencentcloud_redis_instance
    tencentcloud_redis_backup_config

Serverless Cloud Function(SCF)
  Data Source
    tencentcloud_scf_functions
    tencentcloud_scf_logs
    tencentcloud_scf_namespaces

  Resource
    tencentcloud_scf_function
    tencentcloud_scf_namespace

SSL Certificates
  Data Source
    tencentcloud_ssl_certificates

  Resource
    tencentcloud_ssl_certificate

TcaplusDB
  Data Source
    tencentcloud_tcaplus_clusters
    tencentcloud_tcaplus_idls
    tencentcloud_tcaplus_tables
    tencentcloud_tcaplus_tablegroups

  Resource
    tencentcloud_tcaplus_cluster
    tencentcloud_tcaplus_tablegroup
    tencentcloud_tcaplus_idl
    tencentcloud_tcaplus_table

VPC
  Data Source
    tencentcloud_route_table
    tencentcloud_security_group
    tencentcloud_security_groups
    tencentcloud_subnet
    tencentcloud_vpc
    tencentcloud_vpc_instances
    tencentcloud_vpc_route_tables
    tencentcloud_vpc_subnets
    tencentcloud_dnats
    tencentcloud_enis
    tencentcloud_ha_vip_eip_attachments
    tencentcloud_ha_vips
    tencentcloud_nat_gateways
    tencentcloud_nats

  Resource
    tencentcloud_eni
    tencentcloud_eni_attachment
    tencentcloud_vpc
    tencentcloud_subnet
    tencentcloud_security_group
    tencentcloud_security_group_rule
    tencentcloud_security_group_lite_rule
    tencentcloud_route_table
    tencentcloud_route_entry
    tencentcloud_route_table_entry
    tencentcloud_dnat
    tencentcloud_nat_gateway
    tencentcloud_ha_vip
    tencentcloud_ha_vip_eip_attachment

VPN
  Data Source
    tencentcloud_vpn_connections
    tencentcloud_vpn_customer_gateways
    tencentcloud_vpn_gateways

  Resource
    tencentcloud_vpn_customer_gateway
    tencentcloud_vpn_gateway
    tencentcloud_vpn_connection
*/
package tencentcloud

import (
	"net/url"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common"
	sts "github.com/tencentyun/tcecloud-sdk-go/tcecloud/sts/v20180813"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	PROVIDER_SECRET_ID                    = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY                   = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN               = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION                       = "TENCENTCLOUD_REGION"
	PROVIDER_PROTOCOL                     = "TENCENTCLOUD_PROTOCOL"
	PROVIDER_DOMAIN                       = "TENCENTCLOUD_DOMAIN"
	PROVIDER_ASSUME_ROLE_ARN              = "TENCENTCLOUD_ASSUME_ROLE_ARN"
	PROVIDER_ASSUME_ROLE_SESSION_NAME     = "TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME"
	PROVIDER_ASSUME_ROLE_SESSION_DURATION = "TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION"
)

type TencentCloudClient struct {
	apiV3Conn *connectivity.TencentCloudClient
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "This is the TencentCloud access key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "This is the TencentCloud secret key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.",
				Sensitive:   true,
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECURITY_TOKEN, nil),
				Description: "TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).",
				Sensitive:   true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description:  "This is the TencentCloud region. It must be provided, but it can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is ap-guangzhou.",
				InputDefault: "ap-guangzhou",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_PROTOCOL, "HTTPS"),
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				Description:  "The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_DOMAIN, nil),
				Description: "The root domain of the API request, Default is `tencentcloudapi.com`.",
			},
			"assume_role": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "The `assume_role` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_ARN, nil),
							Description: "The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`.",
						},
						"session_name": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_SESSION_NAME, nil),
							Description: "The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`.",
						},
						"session_duration": {
							Type:         schema.TypeInt,
							Required:     true,
							InputDefault: "7200",
							ValidateFunc: validateIntegerInRange(0, 43200),
							Description:  "The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A more restrictive policy when making the AssumeRole call. Its content must not contains `principal` elements. Notice: more syntax references, please refer to: [policies syntax logic](https://intl.cloud.tencent.com/document/product/598/10603).",
						},
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"tencentcloud_availability_regions": dataSourceTencentCloudAvailabilityRegions(),
			"tencentcloud_availability_zones":   dataSourceTencentCloudAvailabilityZones(),
			"tencentcloud_instances":            dataSourceTencentCloudInstances(),
			/*
				"tencentcloud_reserved_instances":           dataSourceTencentCloudReservedInstances(),
			*/
			"tencentcloud_placement_groups": dataSourceTencentCloudPlacementGroups(),
			"tencentcloud_key_pairs":        dataSourceTencentCloudKeyPairs(),
			"tencentcloud_image":            dataSourceTencentCloudImage(),
			"tencentcloud_images":           dataSourceTencentCloudImages(),
			"tencentcloud_instance_types":   dataSourceInstanceTypes(),
			/*
				"tencentcloud_reserved_instance_configs":    dataSourceTencentCloudReservedInstanceConfigs(),
			*/
			"tencentcloud_vpc_instances":    dataSourceTencentCloudVpcInstances(),
			"tencentcloud_vpc_subnets":      dataSourceTencentCloudVpcSubnets(),
			"tencentcloud_vpc_route_tables": dataSourceTencentCloudVpcRouteTables(),
			"tencentcloud_vpc":              dataSourceTencentCloudVpc(),
			"tencentcloud_subnet":           dataSourceTencentCloudSubnet(),
			"tencentcloud_route_table":      dataSourceTencentCloudRouteTable(),
			/*
				"tencentcloud_eip":                          dataSourceTencentCloudEip(),
				"tencentcloud_eips":                         dataSourceTencentCloudEips(),
				"tencentcloud_enis":                         dataSourceTencentCloudEnis(),
				"tencentcloud_nats":                         dataSourceTencentCloudNats(),
				"tencentcloud_dnats":                        dataSourceTencentCloudDnats(),
				"tencentcloud_nat_gateways":                 dataSourceTencentCloudNatGateways(),
				"tencentcloud_vpn_customer_gateways":        dataSourceTencentCloudVpnCustomerGateways(),
				"tencentcloud_vpn_gateways":                 dataSourceTencentCloudVpnGateways(),
				"tencentcloud_vpn_connections":              dataSourceTencentCloudVpnConnections(),
				"tencentcloud_ha_vips":                      dataSourceTencentCloudHaVips(),
				"tencentcloud_ha_vip_eip_attachments":       dataSourceTencentCloudHaVipEipAttachments(),
				"tencentcloud_ccn_instances":                dataSourceTencentCloudCcnInstances(),
				"tencentcloud_ccn_bandwidth_limits":         dataSourceTencentCloudCcnBandwidthLimits(),
				"tencentcloud_dc_instances":                 dataSourceTencentCloudDcInstances(),
				"tencentcloud_dcx_instances":                dataSourceTencentCloudDcxInstances(),
				"tencentcloud_dc_gateway_instances":         dataSourceTencentCloudDcGatewayInstances(),
				"tencentcloud_dc_gateway_ccn_routes":        dataSourceTencentCloudDcGatewayCCNRoutes(),
			*/
			"tencentcloud_security_group":  dataSourceTencentCloudSecurityGroup(),
			"tencentcloud_security_groups": dataSourceTencentCloudSecurityGroups(),
			/*
				"tencentcloud_kubernetes_clusters":          dataSourceTencentCloudKubernetesClusters(),
				"tencentcloud_container_clusters":           dataSourceTencentCloudContainerClusters(),
				"tencentcloud_container_cluster_instances":  dataSourceTencentCloudContainerClusterInstances(),
				"tencentcloud_mysql_backup_list":            dataSourceTencentMysqlBackupList(),
				"tencentcloud_mysql_zone_config":            dataSourceTencentMysqlZoneConfig(),
				"tencentcloud_mysql_parameter_list":         dataSourceTencentCloudMysqlParameterList(),
				"tencentcloud_mysql_instance":               dataSourceTencentCloudMysqlInstance(),
			*/
			"tencentcloud_cos_bucket_object": dataSourceTencentCloudCosBucketObject(),
			"tencentcloud_cos_buckets":       dataSourceTencentCloudCosBuckets(),
			/*
				"tencentcloud_cfs_file_systems":             dataSourceTencentCloudCfsFileSystems(),
				"tencentcloud_cfs_access_groups":            dataSourceTencentCloudCfsAccessGroups(),
				"tencentcloud_cfs_access_rules":             dataSourceTencentCloudCfsAccessRules(),
			*/
			"tencentcloud_redis_zone_config":     dataSourceTencentRedisZoneConfig(),
			"tencentcloud_redis_instances":       dataSourceTencentRedisInstances(),
			"tencentcloud_as_scaling_configs":    dataSourceTencentCloudAsScalingConfigs(),
			"tencentcloud_as_scaling_groups":     dataSourceTencentCloudAsScalingGroups(),
			"tencentcloud_as_scaling_policies":   dataSourceTencentCloudAsScalingPolicies(),
			"tencentcloud_cbs_storages":          dataSourceTencentCloudCbsStorages(),
			"tencentcloud_cbs_snapshots":         dataSourceTencentCloudCbsSnapshots(),
			"tencentcloud_cbs_snapshot_policies": dataSourceTencentCloudCbsSnapshotPolicies(),
			/*
				"tencentcloud_clb_instances":                dataSourceTencentCloudClbInstances(),
				"tencentcloud_clb_listeners":                dataSourceTencentCloudClbListeners(),
				"tencentcloud_clb_listener_rules":           dataSourceTencentCloudClbListenerRules(),
				"tencentcloud_clb_attachments":              dataSourceTencentCloudClbServerAttachments(),
				"tencentcloud_clb_redirections":             dataSourceTencentCloudClbRedirections(),
				"tencentcloud_mongodb_zone_config":          dataSourceTencentCloudMongodbZoneConfig(),
				"tencentcloud_mongodb_instances":            dataSourceTencentCloudMongodbInstances(),
				"tencentcloud_dayu_cc_https_policies":       dataSourceTencentCloudDayuCCHttpsPolicies(),
				"tencentcloud_dayu_cc_http_policies":        dataSourceTencentCloudDayuCCHttpPolicies(),
				"tencentcloud_dayu_ddos_policies":           dataSourceTencentCloudDayuDdosPolicies(),
				"tencentcloud_dayu_ddos_policy_cases":       dataSourceTencentCloudDayuDdosPolicyCases(),
				"tencentcloud_dayu_ddos_policy_attachments": dataSourceTencentCloudDayuDdosPolicyAttachments(),
				"tencentcloud_dayu_l4_rules":                dataSourceTencentCloudDayuL4Rules(),
				"tencentcloud_dayu_l7_rules":                dataSourceTencentCloudDayuL7Rules(),
				"tencentcloud_gaap_proxies":                 dataSourceTencentCloudGaapProxies(),
				"tencentcloud_gaap_realservers":             dataSourceTencentCloudGaapRealservers(),
				"tencentcloud_gaap_layer4_listeners":        dataSourceTencentCloudGaapLayer4Listeners(),
				"tencentcloud_gaap_layer7_listeners":        dataSourceTencentCloudGaapLayer7Listeners(),
				"tencentcloud_gaap_http_domains":            dataSourceTencentCloudGaapHttpDomains(),
				"tencentcloud_gaap_http_rules":              dataSourceTencentCloudGaapHttpRules(),
				"tencentcloud_gaap_security_policies":       dataSourceTencentCloudGaapSecurityPolices(),
				"tencentcloud_gaap_security_rules":          dataSourceTencentCloudGaapSecurityRules(),
				"tencentcloud_gaap_certificates":            dataSourceTencentCloudGaapCertificates(),
				"tencentcloud_gaap_domain_error_pages":      dataSourceTencentCloudGaapDomainErrorPageInfoList(),
				"tencentcloud_ssl_certificates":             dataSourceTencentCloudSslCertificates(),
				"tencentcloud_cam_roles":                    dataSourceTencentCloudCamRoles(),
				"tencentcloud_cam_users":                    dataSourceTencentCloudCamUsers(),
				"tencentcloud_cam_groups":                   dataSourceTencentCloudCamGroups(),
				"tencentcloud_cam_group_memberships":        dataSourceTencentCloudCamGroupMemberships(),
				"tencentcloud_cam_policies":                 dataSourceTencentCloudCamPolicies(),
				"tencentcloud_cam_role_policy_attachments":  dataSourceTencentCloudCamRolePolicyAttachments(),
				"tencentcloud_cam_user_policy_attachments":  dataSourceTencentCloudCamUserPolicyAttachments(),
				"tencentcloud_cam_group_policy_attachments": dataSourceTencentCloudCamGroupPolicyAttachments(),
				"tencentcloud_cam_saml_providers":           dataSourceTencentCloudCamSAMLProviders(),
				"tencentcloud_scf_functions":                dataSourceTencentCloudScfFunctions(),
				"tencentcloud_scf_namespaces":               dataSourceTencentCloudScfNamespaces(),
				"tencentcloud_scf_logs":                     dataSourceTencentCloudScfLogs(),
				"tencentcloud_tcaplus_clusters":             dataSourceTencentCloudTcaplusClusters(),
				"tencentcloud_tcaplus_tablegroups":          dataSourceTencentCloudTcaplusTableGroups(),
				"tencentcloud_tcaplus_tables":               dataSourceTencentCloudTcaplusTables(),
				"tencentcloud_tcaplus_idls":                 dataSourceTencentCloudTcaplusIdls(),
				"tencentcloud_monitor_policy_conditions":    dataSourceTencentMonitorPolicyConditions(),
				"tencentcloud_monitor_data":                 dataSourceTencentMonitorData(),
				"tencentcloud_monitor_product_event":        dataSourceTencentMonitorProductEvent(),
				"tencentcloud_monitor_binding_objects":      dataSourceTencentMonitorBindingObjects(),
				"tencentcloud_monitor_policy_groups":        dataSourceTencentMonitorPolicyGroups(),
				"tencentcloud_monitor_product_namespace":    dataSourceTencentMonitorProductNamespace(),
				"tencentcloud_elasticsearch_instances":      dataSourceTencentCloudElasticsearchInstances(),
			*/
		},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_instance": resourceTencentCloudInstance(),
			/*
				"tencentcloud_reserved_instance":              resourceTencentCloudReservedInstance(),
			*/
			"tencentcloud_key_pair":                       resourceTencentCloudKeyPair(),
			"tencentcloud_placement_group":                resourceTencentCloudPlacementGroup(),
			"tencentcloud_cbs_snapshot":                   resourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_snapshot_policy":            resourceTencentCloudCbsSnapshotPolicy(),
			"tencentcloud_cbs_storage":                    resourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_attachment":         resourceTencentCloudCbsStorageAttachment(),
			"tencentcloud_cbs_snapshot_policy_attachment": resourceTencentCloudCbsSnapshotPolicyAttachment(),
			"tencentcloud_vpc":                            resourceTencentCloudVpcInstance(),
			"tencentcloud_subnet":                         resourceTencentCloudVpcSubnet(),
			"tencentcloud_route_entry":                    resourceTencentCloudRouteEntry(),
			"tencentcloud_route_table_entry":              resourceTencentCloudVpcRouteEntry(),
			"tencentcloud_route_table":                    resourceTencentCloudVpcRouteTable(),
			/*
				"tencentcloud_dnat":                           resourceTencentCloudDnat(),
				"tencentcloud_nat_gateway":                    resourceTencentCloudNatGateway(),
				"tencentcloud_eip":                            resourceTencentCloudEip(),
				"tencentcloud_eip_association":                resourceTencentCloudEipAssociation(),
				"tencentcloud_eni":                            resourceTencentCloudEni(),
				"tencentcloud_eni_attachment":                 resourceTencentCloudEniAttachment(),
				"tencentcloud_ccn":                            resourceTencentCloudCcn(),
				"tencentcloud_ccn_attachment":                 resourceTencentCloudCcnAttachment(),
				"tencentcloud_ccn_bandwidth_limit":            resourceTencentCloudCcnBandwidthLimit(),
				"tencentcloud_dcx":                            resourceTencentCloudDcxInstance(),
				"tencentcloud_dc_gateway":                     resourceTencentCloudDcGatewayInstance(),
				"tencentcloud_dc_gateway_ccn_route":           resourceTencentCloudDcGatewayCcnRouteInstance(),
				"tencentcloud_vpn_customer_gateway":           resourceTencentCloudVpnCustomerGateway(),
				"tencentcloud_vpn_gateway":                    resourceTencentCloudVpnGateway(),
				"tencentcloud_vpn_connection":                 resourceTencentCloudVpnConnection(),
				"tencentcloud_ha_vip":                         resourceTencentCloudHaVip(),
				"tencentcloud_ha_vip_eip_attachment":          resourceTencentCloudHaVipEipAttachment(),
			*/
			"tencentcloud_security_group":           resourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":      resourceTencentCloudSecurityGroupRule(),
			"tencentcloud_security_group_lite_rule": resourceTencentCloudSecurityGroupLiteRule(),
			/*
				"tencentcloud_lb":                            resourceTencentCloudLB(),
				"tencentcloud_alb_server_attachment":         resourceTencentCloudAlbServerAttachment(),
				"tencentcloud_clb_instance":                  resourceTencentCloudClbInstance(),
				"tencentcloud_clb_listener":                  resourceTencentCloudClbListener(),
				"tencentcloud_clb_listener_rule":             resourceTencentCloudClbListenerRule(),
				"tencentcloud_clb_attachment":                resourceTencentCloudClbServerAttachment(),
				"tencentcloud_clb_redirection":               resourceTencentCloudClbRedirection(),
				"tencentcloud_container_cluster":             resourceTencentCloudContainerCluster(),
				"tencentcloud_container_cluster_instance":    resourceTencentCloudContainerClusterInstance(),
				"tencentcloud_kubernetes_cluster":            resourceTencentCloudTkeCluster(),
				"tencentcloud_kubernetes_as_scaling_group":   ResourceTencentCloudKubernetesAsScalingGroup(),
				"tencentcloud_kubernetes_scale_worker":       resourceTencentCloudTkeScaleWorker(),
				"tencentcloud_kubernetes_cluster_attachment": resourceTencentCloudTkeClusterAttachment(),
				"tencentcloud_mysql_backup_policy":           resourceTencentCloudMysqlBackupPolicy(),
				"tencentcloud_mysql_account":                 resourceTencentCloudMysqlAccount(),
				"tencentcloud_mysql_account_privilege":       resourceTencentCloudMysqlAccountPrivilege(),
				"tencentcloud_mysql_privilege":               resourceTencentCloudMysqlPrivilege(),
				"tencentcloud_mysql_instance":                resourceTencentCloudMysqlInstance(),
				"tencentcloud_mysql_readonly_instance":       resourceTencentCloudMysqlReadonlyInstance(),
			*/
			"tencentcloud_cos_bucket":        resourceTencentCloudCosBucket(),
			"tencentcloud_cos_bucket_object": resourceTencentCloudCosBucketObject(),
			/*
				"tencentcloud_cfs_file_system":             resourceTencentCloudCfsFileSystem(),
				"tencentcloud_cfs_access_group":            resourceTencentCloudCfsAccessGroup(),
				"tencentcloud_cfs_access_rule":             resourceTencentCloudCfsAccessRule(),
			*/
			"tencentcloud_redis_instance":      resourceTencentCloudRedisInstance(),
			"tencentcloud_redis_backup_config": resourceTencentCloudRedisBackupConfig(),
			"tencentcloud_as_scaling_config":   resourceTencentCloudAsScalingConfig(),
			"tencentcloud_as_scaling_group":    resourceTencentCloudAsScalingGroup(),
			"tencentcloud_as_attachment":       resourceTencentCloudAsAttachment(),
			"tencentcloud_as_scaling_policy":   resourceTencentCloudAsScalingPolicy(),
			"tencentcloud_as_schedule":         resourceTencentCloudAsSchedule(),
			"tencentcloud_as_lifecycle_hook":   resourceTencentCloudAsLifecycleHook(),
			"tencentcloud_as_notification":     resourceTencentCloudAsNotification(),
			/*
				"tencentcloud_mongodb_instance":            resourceTencentCloudMongodbInstance(),
				"tencentcloud_mongodb_sharding_instance":   resourceTencentCloudMongodbShardingInstance(),
				"tencentcloud_dayu_cc_http_policy":         resourceTencentCloudDayuCCHttpPolicy(),
				"tencentcloud_dayu_cc_https_policy":        resourceTencentCloudDayuCCHttpsPolicy(),
				"tencentcloud_dayu_ddos_policy":            resourceTencentCloudDayuDdosPolicy(),
				"tencentcloud_dayu_ddos_policy_case":       resourceTencentCloudDayuDdosPolicyCase(),
				"tencentcloud_dayu_ddos_policy_attachment": resourceTencentCloudDayuDdosPolicyAttachment(),
				"tencentcloud_dayu_l4_rule":                resourceTencentCloudDayuL4Rule(),
				"tencentcloud_dayu_l7_rule":                resourceTencentCloudDayuL7Rule(),
				"tencentcloud_gaap_proxy":                  resourceTencentCloudGaapProxy(),
				"tencentcloud_gaap_realserver":             resourceTencentCloudGaapRealserver(),
				"tencentcloud_gaap_layer4_listener":        resourceTencentCloudGaapLayer4Listener(),
				"tencentcloud_gaap_layer7_listener":        resourceTencentCloudGaapLayer7Listener(),
				"tencentcloud_gaap_http_domain":            resourceTencentCloudGaapHttpDomain(),
				"tencentcloud_gaap_http_rule":              resourceTencentCloudGaapHttpRule(),
				"tencentcloud_gaap_certificate":            resourceTencentCloudGaapCertificate(),
				"tencentcloud_gaap_security_policy":        resourceTencentCloudGaapSecurityPolicy(),
				"tencentcloud_gaap_security_rule":          resourceTencentCloudGaapSecurityRule(),
				"tencentcloud_gaap_domain_error_page":      resourceTencentCloudGaapDomainErrorPageInfo(),
				"tencentcloud_ssl_certificate":             resourceTencentCloudSslCertificate(),
				"tencentcloud_cam_role":                    resourceTencentCloudCamRole(),
				"tencentcloud_cam_user":                    resourceTencentCloudCamUser(),
				"tencentcloud_cam_policy":                  resourceTencentCloudCamPolicy(),
				"tencentcloud_cam_role_policy_attachment":  resourceTencentCloudCamRolePolicyAttachment(),
				"tencentcloud_cam_user_policy_attachment":  resourceTencentCloudCamUserPolicyAttachment(),
				"tencentcloud_cam_group_policy_attachment": resourceTencentCloudCamGroupPolicyAttachment(),
				"tencentcloud_cam_group":                   resourceTencentCloudCamGroup(),
				"tencentcloud_cam_group_membership":        resourceTencentCloudCamGroupMembership(),
				"tencentcloud_cam_saml_provider":           resourceTencentCloudCamSAMLProvider(),
				"tencentcloud_scf_function":                resourceTencentCloudScfFunction(),
				"tencentcloud_scf_namespace":               resourceTencentCloudScfNamespace(),
				"tencentcloud_tcaplus_cluster":             resourceTencentCloudTcaplusCluster(),
				"tencentcloud_tcaplus_group":               resourceTencentCloudTcaplusGroup(),
				"tencentcloud_tcaplus_idl":                 resourceTencentCloudTcaplusIdl(),
				"tencentcloud_tcaplus_table":               resourceTencentCloudTcaplusTable(),
				"tencentcloud_cdn_domain":                  resourceTencentCloudCdnDomain(),
				"tencentcloud_monitor_policy_group":        resourceTencentMonitorPolicyGroup(),
				"tencentcloud_monitor_binding_object":      resourceTencentMonitorBindingObject(),
				"tencentcloud_monitor_binding_receiver":    resourceTencentMonitorBindingAlarmReceiver(),
				"tencentcloud_elasticsearch_instance":         resourceTencentCloudElasticsearchInstance(),
			*/
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	secretId := d.Get("secret_id").(string)
	secretKey := d.Get("secret_key").(string)
	securityToken := d.Get("security_token").(string)
	region := d.Get("region").(string)

	protocol := d.Get("protocol").(string)
	domain := d.Get("domain").(string)

	// standard client
	var tcClient TencentCloudClient
	tcClient.apiV3Conn = &connectivity.TencentCloudClient{
		Credential: common.NewTokenCredential(
			secretId,
			secretKey,
			securityToken,
		),
		Region:   region,
		Protocol: protocol,
		Domain:   domain,
	}

	// assume role client
	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		assumeRoleArn := assumeRole["role_arn"].(string)
		assumeRoleSessionName := assumeRole["session_name"].(string)
		assumeRoleSessionDuration := assumeRole["session_duration"].(int)
		assumeRolePolicy := assumeRole["policy"].(string)
		if assumeRoleSessionDuration == 0 {
			var err error
			if duration := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_DURATION); duration != "" {
				assumeRoleSessionDuration, err = strconv.Atoi(duration)
				if err != nil {
					return nil, err
				}
				if assumeRoleSessionDuration == 0 {
					assumeRoleSessionDuration = 7200
				}
			}
		}
		// applying STS credentials
		request := sts.NewAssumeRoleRequest()
		request.RoleArn = helper.String(assumeRoleArn)
		request.RoleSessionName = helper.String(assumeRoleSessionName)
		request.DurationSeconds = helper.IntUint64(assumeRoleSessionDuration)
		if assumeRolePolicy != "" {
			request.Policy = helper.String(url.QueryEscape(assumeRolePolicy))
		}
		ratelimit.Check(request.GetAction())
		response, err := tcClient.apiV3Conn.UseStsClient().AssumeRole(request)
		if err != nil {
			return nil, err
		}
		// using STS credentials
		tcClient.apiV3Conn.Credential = common.NewTokenCredential(
			*response.Response.Credentials.TmpSecretId,
			*response.Response.Credentials.TmpSecretKey,
			*response.Response.Credentials.Token,
		)
	}

	return &tcClient, nil
}

package connectivity

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	as "github.com/tencentyun/tcecloud-sdk-go/tcecloud/as/v20180419"

	//cam "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cam/v20190116"
	cbs "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cbs/v20170312"
	//cdb "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cdb/v20170320"
	//cdn "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cdn/v20180606"
	cfs "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cfs/v20190719"
	clb "github.com/tencentyun/tcecloud-sdk-go/tcecloud/clb/v20180317"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common"
	"github.com/tencentyun/tcecloud-sdk-go/tcecloud/common/profile"
	cvm "github.com/tencentyun/tcecloud-sdk-go/tcecloud/cvm/v20170312"

	//dayu "github.com/tencentyun/tcecloud-sdk-go/tcecloud/dayu/v20180709"
	dc "github.com/tencentyun/tcecloud-sdk-go/tcecloud/dc/v20180410"
	//gaap "github.com/tencentyun/tcecloud-sdk-go/tcecloud/gaap/v20180529"
	//mongodb "github.com/tencentyun/tcecloud-sdk-go/tcecloud/mongodb/v20180408"
	monitor "github.com/tencentyun/tcecloud-sdk-go/tcecloud/monitor/v20180724"
	redis "github.com/tencentyun/tcecloud-sdk-go/tcecloud/redis/v20180412"

	//scf "github.com/tencentyun/tcecloud-sdk-go/tcecloud/scf/v20180416"
	//tag "github.com/tencentyun/tcecloud-sdk-go/tcecloud/tag/v20180813"
	//tcaplusdb "github.com/tencentyun/tcecloud-sdk-go/tcecloud/tcaplusdb/v20190823"
	tke "github.com/tencentyun/tcecloud-sdk-go/tcecloud/tke/v20180525"
	vpc "github.com/tencentyun/tcecloud-sdk-go/tcecloud/vpc/v20170312"
	//ssl "github.com/tencentyun/tcecloud-sdk-go/tcecloud/wss/v20180426"
	//sts "github.com/tencentyun/tcecloud-sdk-go/tcecloud/sts/v20180813"
)

// TencentCloudClient is client for all TencentCloud service
type TencentCloudClient struct {
	Credential *common.Credential
	Region     string
	Protocol   string
	Domain     string

	cosConn *s3.S3
	//mysqlConn   *cdb.Client
	redisConn *redis.Client
	asConn    *as.Client
	vpcConn   *vpc.Client
	cbsConn   *cbs.Client
	cvmConn   *cvm.Client
	clbConn   *clb.Client
	//dayuConn    *dayu.Client
	dcConn *dc.Client
	//tagConn     *tag.Client
	//mongodbConn *mongodb.Client
	tkeConn *tke.Client
	//stsConn *sts.Client
	//camConn     *cam.Client
	//gaapConn    *gaap.Client
	//sslConn     *ssl.Client
	cfsConn *cfs.Client
	//scfConn     *scf.Client
	//tcaplusConn *tcaplusdb.Client
	//cdnConn     *cdn.Client
	monitorConn *monitor.Client
	//esConn      *es.Client
}

// NewClientProfile returns a new ClientProfile
func (me *TencentCloudClient) NewClientProfile(timeout int) *profile.ClientProfile {
	cpf := profile.NewClientProfile()

	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = timeout
	// request protocol
	cpf.HttpProfile.Scheme = me.Protocol
	// request domain
	cpf.HttpProfile.RootDomain = me.Domain
	// default language
	//cpf.Language = "en-US"

	return cpf
}

// UseCosClient returns cos client for service
func (me *TencentCloudClient) UseCosClient() *s3.S3 {
	if me.cosConn != nil {
		return me.cosConn
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           fmt.Sprintf("https://cos.%s.myqcloud.com", region),
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	creds := credentials.NewStaticCredentials(me.Credential.SecretId, me.Credential.SecretKey, me.Credential.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return s3.New(sess)
}

/*
// UseMysqlClient returns mysql(cdb) client for service
func (me *TencentCloudClient) UseMysqlClient() *cdb.Client {
	if me.mysqlConn != nil {
		return me.mysqlConn
	}

	cpf := me.NewClientProfile(300)
	me.mysqlConn, _ = cdb.NewClient(me.Credential, me.Region, cpf)
	me.mysqlConn.WithHttpTransport(&LogRoundTripper{})

	return me.mysqlConn
}
*/
// UseRedisClient returns redis client for service
func (me *TencentCloudClient) UseRedisClient() *redis.Client {
	if me.redisConn != nil {
		return me.redisConn
	}

	cpf := me.NewClientProfile(300)
	me.redisConn, _ = redis.NewClient(me.Credential, me.Region, cpf)
	me.redisConn.WithHttpTransport(&LogRoundTripper{})

	return me.redisConn
}

// UseAsClient returns as client for service
func (me *TencentCloudClient) UseAsClient() *as.Client {
	if me.asConn != nil {
		return me.asConn
	}

	cpf := me.NewClientProfile(300)
	me.asConn, _ = as.NewClient(me.Credential, me.Region, cpf)
	me.asConn.WithHttpTransport(&LogRoundTripper{})

	return me.asConn
}

// UseVpcClient returns vpc client for service
func (me *TencentCloudClient) UseVpcClient() *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}

	cpf := me.NewClientProfile(300)
	me.vpcConn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	me.vpcConn.WithHttpTransport(&LogRoundTripper{})

	return me.vpcConn
}

// UseCbsClient returns cbs client for service
func (me *TencentCloudClient) UseCbsClient() *cbs.Client {
	if me.cbsConn != nil {
		return me.cbsConn
	}

	cpf := me.NewClientProfile(300)
	me.cbsConn, _ = cbs.NewClient(me.Credential, me.Region, cpf)
	me.cbsConn.WithHttpTransport(&LogRoundTripper{})

	return me.cbsConn
}

// UseDcClient returns dc client for service
func (me *TencentCloudClient) UseDcClient() *dc.Client {
	if me.dcConn != nil {
		return me.dcConn
	}

	cpf := me.NewClientProfile(300)
	me.dcConn, _ = dc.NewClient(me.Credential, me.Region, cpf)
	me.dcConn.WithHttpTransport(&LogRoundTripper{})

	return me.dcConn
}

/*
// UseMongodbClient returns mongodb client for service
func (me *TencentCloudClient) UseMongodbClient() *mongodb.Client {
	if me.mongodbConn != nil {
		return me.mongodbConn
	}

	cpf := me.NewClientProfile(300)
	me.mongodbConn, _ = mongodb.NewClient(me.Credential, me.Region, cpf)
	me.mongodbConn.WithHttpTransport(&LogRoundTripper{})

	return me.mongodbConn
}
*/
// UseClbClient returns clb client for service
func (me *TencentCloudClient) UseClbClient() *clb.Client {
	if me.clbConn != nil {
		return me.clbConn
	}

	cpf := me.NewClientProfile(300)
	me.clbConn, _ = clb.NewClient(me.Credential, me.Region, cpf)
	me.clbConn.WithHttpTransport(&LogRoundTripper{})

	return me.clbConn
}

// UseCvmClient returns cvm client for service
func (me *TencentCloudClient) UseCvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	cpf := me.NewClientProfile(300)
	me.cvmConn, _ = cvm.NewClient(me.Credential, me.Region, cpf)
	me.cvmConn.WithHttpTransport(&LogRoundTripper{})

	return me.cvmConn
}

/*
// UseTagClient returns tag client for service
func (me *TencentCloudClient) UseTagClient() *tag.Client {
	if me.tagConn != nil {
		return me.tagConn
	}

	cpf := me.NewClientProfile(300)
	me.tagConn, _ = tag.NewClient(me.Credential, me.Region, cpf)
	me.tagConn.WithHttpTransport(&LogRoundTripper{})

	return me.tagConn
}
*/
// UseTkeClient returns tke client for service
func (me *TencentCloudClient) UseTkeClient() *tke.Client {
	if me.tkeConn != nil {
		return me.tkeConn
	}

	cpf := me.NewClientProfile(300)
	me.tkeConn, _ = tke.NewClient(me.Credential, me.Region, cpf)
	me.tkeConn.WithHttpTransport(&LogRoundTripper{})

	return me.tkeConn
}

/*
// UseGaapClient returns gaap client for service
func (me *TencentCloudClient) UseGaapClient() *gaap.Client {
	if me.gaapConn != nil {
		return me.gaapConn
	}

	cpf := me.NewClientProfile(300)
	me.gaapConn, _ = gaap.NewClient(me.Credential, me.Region, cpf)
	me.gaapConn.WithHttpTransport(&LogRoundTripper{})

	return me.gaapConn
}
*/
/*
// UseSslClient returns ssl client for service
func (me *TencentCloudClient) UseSslClient() *ssl.Client {
	if me.sslConn != nil {
		return me.sslConn
	}

	cpf := me.NewClientProfile(300)
	me.sslConn, _ = ssl.NewClient(me.Credential, me.Region, cpf)
	me.sslConn.WithHttpTransport(&LogRoundTripper{})

	return me.sslConn
}
*/
/*
// UseCamClient returns cam client for service
func (me *TencentCloudClient) UseCamClient() *cam.Client {
	if me.camConn != nil {
		return me.camConn
	}

	cpf := me.NewClientProfile(300)
	me.camConn, _ = cam.NewClient(me.Credential, me.Region, cpf)
	me.camConn.WithHttpTransport(&LogRoundTripper{})

	return me.camConn
}
*/
/*
// UseStsClient returns sts client for service
func (me *TencentCloudClient) UseStsClient() *sts.Client {
	/-*
		me.Credential will changed, don't cache it
		if me.stsConn != nil {
			return me.stsConn
		}
	*-/

	cpf := me.NewClientProfile(300)
	me.stsConn, _ = sts.NewClient(me.Credential, me.Region, cpf)
	me.stsConn.WithHttpTransport(&LogRoundTripper{})

	return me.stsConn
}
*/
// UseCfsClient returns cfs client for service
func (me *TencentCloudClient) UseCfsClient() *cfs.Client {
	if me.cfsConn != nil {
		return me.cfsConn
	}

	cpf := me.NewClientProfile(300)
	me.cfsConn, _ = cfs.NewClient(me.Credential, me.Region, cpf)
	me.cfsConn.WithHttpTransport(&LogRoundTripper{})

	return me.cfsConn
}

/*
// UseScfClient returns scf client for service
func (me *TencentCloudClient) UseScfClient() *scf.Client {
	if me.scfConn != nil {
		return me.scfConn
	}

	cpf := me.NewClientProfile(300)
	me.scfConn, _ = scf.NewClient(me.Credential, me.Region, cpf)
	me.scfConn.WithHttpTransport(&LogRoundTripper{})

	return me.scfConn
}
*/
/*
// UseTcaplusClient returns tcaplush client for service
func (me *TencentCloudClient) UseTcaplusClient() *tcaplusdb.Client {
	if me.tcaplusConn != nil {
		return me.tcaplusConn
	}

	cpf := me.NewClientProfile(300)
	me.tcaplusConn, _ = tcaplusdb.NewClient(me.Credential, me.Region, cpf)
	me.tcaplusConn.WithHttpTransport(&LogRoundTripper{})

	return me.tcaplusConn
}
*/
/*
// UseDayuClient returns dayu client for service
func (me *TencentCloudClient) UseDayuClient() *dayu.Client {
	if me.dayuConn != nil {
		return me.dayuConn
	}

	cpf := me.NewClientProfile(300)
	me.dayuConn, _ = dayu.NewClient(me.Credential, me.Region, cpf)
	me.dayuConn.WithHttpTransport(&LogRoundTripper{})

	return me.dayuConn
}
*/
/*
// UseCdnClient returns cdn client for service
func (me *TencentCloudClient) UseCdnClient() *cdn.Client {
	if me.cdnConn != nil {
		return me.cdnConn
	}

	cpf := me.NewClientProfile(300)
	me.cdnConn, _ = cdn.NewClient(me.Credential, me.Region, cpf)
	me.cdnConn.WithHttpTransport(&LogRoundTripper{})

	return me.cdnConn
}
*/
// UseMonitorClient returns monitor client for service
func (me *TencentCloudClient) UseMonitorClient() *monitor.Client {
	if me.monitorConn != nil {
		return me.monitorConn
	}

	cpf := me.NewClientProfile(300)
	me.monitorConn, _ = monitor.NewClient(me.Credential, me.Region, cpf)
	me.monitorConn.WithHttpTransport(&LogRoundTripper{})

	return me.monitorConn
}

/*
// UseEsClient returns es client for service
func (me *TencentCloudClient) UseEsClient() *es.Client {
	if me.esConn != nil {
		return me.esConn
	}

	cpf := me.NewClientProfile(300)
	me.esConn, _ = es.NewClient(me.Credential, me.Region, cpf)
	me.esConn.WithHttpTransport(&LogRoundTripper{})

	return me.esConn
}
*/

package adapters

import (
    "context"
    "encoding/base64"
    "fmt"
    "github.com/cihub/seelog"
    "github.com/tencentyun/cos-go-sdk-v5"
    "go_another_chatgpt/repo/censor"
    "net/http"
    "net/url"
)

type TencentCi struct {
    Client *cos.Client
}

func NewTencentCi(secretId, secretKey string, bucketUrl, serviceUrl, ciUrl string) (r *TencentCi) {
    r = new(TencentCi)
    // 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
    // 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
    // COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
    u, _ := url.Parse(bucketUrl)
    // 用于Get Service 查询，默认全地域 service.cos.myqcloud.com
    var su, cu *url.URL
    if serviceUrl != "" {
        su, _ = url.Parse(serviceUrl)
    }
    if ciUrl != "" {
        cu, _ = url.Parse(ciUrl)
    }

    seelog.Infof("检查参数: %+v, %+v", u, su)

    b := &cos.BaseURL{BucketURL: u, ServiceURL: su, CIURL: cu}
    // 1.永久密钥
    r.Client = cos.NewClient(b, &http.Client{
        Transport: &cos.AuthorizationTransport{
            SecretID:  secretId,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
            SecretKey: secretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
        },
    })
    return
}

func (r *TencentCi) MakeTextAuditing(id, text string) (result *censor.TextAuditingResult, err error) {
    result = new(censor.TextAuditingResult)
    //name := fmt.Sprintf("censor/%v.txt", id)
    //// 1.通过字符串上传对象
    //f := strings.NewReader(text)
    //
    //_, err = r.Client.Object.Put(context.Background(), name, f, nil)
    //if err != nil {
    //    err = fmt.Errorf("r.Client.Object.Put failed, %w", err)
    //    return
    //}

    // 购买 @see https://buy.cloud.tencent.com/ci
    // 文档 @see https://cloud.tencent.com/document/product/460/72958
    //seelog.Infof("开始审核")

    encoded := base64.StdEncoding.EncodeToString([]byte(text))

    opt := &cos.PutTextAuditingJobOptions{
        InputContent: encoded,
        Conf: &cos.TextAuditingJobConf{
            BizType: "55ae869651cceebb53473c51d67fa4ba",
        },
        InputDataId: fmt.Sprintf("%v", id),
    }

    var res *cos.PutTextAuditingJobResult
    // 返回结果： https://cloud.tencent.com/document/product/436/56288
    res, _, err = r.Client.CI.PutTextAuditingJob(context.Background(), opt)
    if err != nil {
        err = fmt.Errorf("r.Client.CI.PutTextAuditingJob, err=%w", err)
        return
    }

    if res.JobsDetail == nil {
        err = fmt.Errorf("r.Client.CI.PutTextAuditingJob no JobsDetail")
        return
    }
    msg := ""
    if res.JobsDetail.PornInfo != nil {
        msg += fmt.Sprintf(" PornInfo=%+v", *res.JobsDetail.PornInfo)
    }
    if res.JobsDetail.AdsInfo != nil {
        msg += fmt.Sprintf(" AdsInfo=%+v", *res.JobsDetail.AdsInfo)
    }
    if res.JobsDetail.IllegalInfo != nil {
        msg += fmt.Sprintf(" IllegalInfo=%+v", *res.JobsDetail.IllegalInfo)
    }
    if res.JobsDetail.AbuseInfo != nil {
        msg += fmt.Sprintf(" AbuseInfo=%+v", *res.JobsDetail.AbuseInfo)
    }
    seelog.Infof("JobsDetail(id=%v): result=%v, label=%v", id, res.JobsDetail.Result, res.JobsDetail.Label)
    seelog.Debugf("JobsDetail(id=%v): result=%v, label=%v, %v, %#v", id, res.JobsDetail.Result, res.JobsDetail.Label,
        msg,
        *res.JobsDetail)

    if res.JobsDetail.Result == 1 || (res.JobsDetail.Result == 2 && res.JobsDetail.Label == "Politics") {
        return
    }

    result.Safe = true
    return
}

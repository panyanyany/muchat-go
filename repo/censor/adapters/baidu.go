package adapters

import (
    "encoding/json"
    "fmt"
    baidu_censor "github.com/Baidu-AIP/golang-sdk/aip/censor"
    "go_another_chatgpt/repo/censor"
)

type Baidu struct {
    Client *baidu_censor.ContentCensorClient
}

func NewBaidu(apiKey, secretKey string) (r *Baidu) {
    r = new(Baidu)
    //如果是百度云ak sk,使用下面的客户端
    r.Client = baidu_censor.NewClient(apiKey, secretKey)
    return
}

func (r *Baidu) MakeTextAuditing(id, text string) (result *censor.TextAuditingResult, err error) {
    result = new(censor.TextAuditingResult)
    t := r.Client.TextCensor(text)
    fmt.Printf("res: %v\n", t)
    failed := TextCensorResponseFailed{}
    err = json.Unmarshal([]byte(t), &failed)
    if err != nil {
        err = fmt.Errorf("json.Unmarshal result body failed-0: %w", err)
        return
    }
    if failed.ErrorMsg != "" {
        err = fmt.Errorf("baidu TextCensor failed: %v", t)
        return
    }
    success := TextCensorResponseSuccess{}
    err = json.Unmarshal([]byte(t), &success)
    if err != nil {
        err = fmt.Errorf("json.Unmarshal result body failed-1: %w", err)
        return
    }

    //msg := fmt.Sprintf("(id=%v)：%v", id, success.Conclusion)
    if success.ConclusionType == 1 || success.ConclusionType == 3 {
        result.Safe = true
        //seelog.Infof("审核通过：%v", msg)
        return
    }
    //seelog.Infof("审核不通过：%v", msg)
    return
}

// TextCensorResponseSuccess https://ai.baidu.com/ai-doc/ANTIPORN/2kvuvd2pr#%E8%BF%94%E5%9B%9E%E5%8F%82%E6%95%B0%E8%AF%A6%E6%83%85-1
type TextCensorResponseSuccess struct {
    LogID          int64  `json:"log_id"`
    Conclusion     string `json:"conclusion"`     // 审核结果，可取值：合规、不合规、疑似、审核失败
    ConclusionType int    `json:"conclusionType"` // 审核结果类型，可取值1.合规，2.不合规，3.疑似，4.审核失败
    Data           []struct {
        Type           int    `json:"type"`
        SubType        int    `json:"subType"`
        Conclusion     string `json:"conclusion"`
        ConclusionType int    `json:"conclusionType"`
        Msg            string `json:"msg"`
        Hits           []struct {
            DatasetName string   `json:"datasetName"`
            Words       []string `json:"words"`
        } `json:"hits"`
    } `json:"data"`
}

type TextCensorResponseFailed struct {
    LogID     int64  `json:"log_id"`
    ErrorCode int    `json:"error_code,omitempty"`
    ErrorMsg  string `json:"error_msg,omitempty"`
}

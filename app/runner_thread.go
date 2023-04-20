package app

import (
    "encoding/json"
    "fmt"
    "github.com/cihub/seelog"
    "muchat-go/chatgpt/api_base"
    "muchat-go/config"
    "muchat-go/repo/constants"
    "muchat-go/util/utf8_util"
    "strings"
    "time"
)

type RunnerThread struct {
    Jobs       chan *JobParam
    ResultJobs chan *JobParam
    ShouldStop chan bool
    Stopped    chan bool
    GptConfig  *config.GptConfig
    Config     *config.Configuration
    Running    bool
    Index      int
}

func NewRunnerThread(index int, config *config.Configuration, jobs chan *JobParam, resultJobs chan *JobParam) (r *RunnerThread) {
    r = new(RunnerThread)
    r.Index = index
    r.Jobs = jobs
    r.ResultJobs = resultJobs
    r.Config = config
    r.GptConfig = config.Gpt
    r.ShouldStop = make(chan bool)
    r.Stopped = make(chan bool)
    return
}

func (r *RunnerThread) Stop() {
    seelog.Infof("正在暂停线程 #%v......", r.Index)
    r.ShouldStop <- true
}
func (r *RunnerThread) WaitForStop() {
    <-r.Stopped
    seelog.Infof("已退出线程 #%v......", r.Index)
}

func (r *RunnerThread) Start() {
    for {
        select {
        case <-r.ShouldStop:
            r.Stopped <- true
            return
        case job := <-r.Jobs:
            r.DoJob(job)
        }
    }
}

func (r *RunnerThread) DoJob(job *JobParam) {
    var err error
    var msg string
    acc := job.Account
    title := strings.TrimSpace(job.Title) + "\n"
    idMsg := fmt.Sprintf("slug=%v, thread=%v, jobId=%v, acc=%v", job.Slug, r.Index, job.TitleIndex, acc.Email)

    defer func() {
        job.DoneTime = time.Now()
        r.ResultJobs <- job
    }()

    messages := job.Messages
    if messages != nil && len(messages) > 0 {
        // 大概不能超过这个数
        maxLen := 1500 // 3500最极限，但会导致回复很短，
        textCnt := 0
        messages2 := []api_base.ChatMessage{}
        for _, m := range messages {
            curTextCnt := textCnt + utf8_util.Len(m.Content)
            if curTextCnt > maxLen {
                if textCnt == 0 { // 唯一的一个问题太长
                    m.Content = utf8_util.Substr(m.Content, len(m.Content)-maxLen, maxLen)
                    messages2 = append(messages2, m)
                }
                break
            }
            textCnt = curTextCnt
            messages2 = append(messages2, m)
        }
        messages = messages2
        bs, err := json.Marshal(messages)
        if err != nil {
            seelog.Errorf("json.Marshall messages failed: messages=%#v, err=%v", messages, err)
        } else {
            seelog.Infof("正在查询题目（%v）, idMsg=%v", string(bs), idMsg)
        }
    } else {
        // 大概不能超过这个数
        maxLen := 1900
        title = utf8_util.Substr(title, len(title)-maxLen, maxLen)
        seelog.Infof("正在查询题目（%#v）, idMsg=%v", title, idMsg)
    }

    // 审核题目
    var textToCensor string
    if messages != nil && len(messages) > 0 {
        textToCensor = messages[len(messages)-1].Content
    } else {
        pieces := strings.Split(title, "\n\n")
        textToCensor = pieces[len(pieces)-1]
    }

    job.Censor.AddJob(job.TitleIndex, textToCensor)
    censorJob := job.Censor.WaitForJob(job.TitleIndex)
    if censorJob.Err != nil {
        seelog.Errorf("审核出错-0: %v, idMsg=%v", censorJob.Err, idMsg)
    }
    if !censorJob.Safe {
        job.Error = fmt.Errorf("敏感问题，已拦截")
        seelog.Infof("审核不通过，问题：%+v, idMsg=%v", textToCensor, idMsg)
        job.ErrCode = constants.CodeErrSensitive
        job.UsedCap++
        job.BadCnt++
        return
    }

    // 查询 GPT
    beg := time.Now()

    if r.Config.Mock.Enabled {
        msg = r.Config.Mock.Response
    } else {
        if messages != nil && len(messages) > 0 {
            msg, err = api_base.ChatCompletions(messages, acc.ApiKey, r.Config, job.Slug)
        } else {
            msg, err = api_base.Completions(title, acc.ApiKey, r.GptConfig, job.Slug)
        }
        if err != nil {
            seelog.Errorf("ChatGPT接口异常，题目(%#v)，错误=%v, idMsg=%v", title, err, idMsg)
            job.Error = fmt.Errorf("ChatGPT接口异常")
            job.ErrCode = constants.CodeErrGptError
            job.ErrCnt++
            job.InnerError = err
            job.InnerErrorMsg = msg
            return
        }
    }

    end := time.Now()
    duration := end.Sub(beg).Seconds()
    seelog.Infof("查询成功，耗时: %v, idMsg=%v", duration, idMsg)

    msg = strings.TrimSpace(msg)

    job.Message = msg
    job.Error = nil

    if len(job.Message) == 0 {
        job.Message = "您的问题我无法回答"
    }
    if r.Config.Mock.IsFree() {
        // 不计费
        return
    }
    // 计费
    job.UsedCap++
    // 是否审核答案
    if r.Config.Mock.Enabled || r.Config.CensorEnabled == false || strings.Contains(job.Message, "已拦截") {
        // 不审核
        seelog.Infof("不审核，答案：%#v, idMsg=%v", job.Message, idMsg)
        return
    }
    // 审核
    job.Censor.AddJob(job.TitleIndex, job.Message)
    censorJob = job.Censor.WaitForJob(job.TitleIndex)
    if censorJob.Err != nil {
        seelog.Errorf("审核出错-1: %v, idMsg=%v", censorJob.Err, idMsg)
        job.Message = ""
        job.Error = fmt.Errorf("出错了，请稍后再试")
        job.ErrCode = constants.CodeErrUnknown
        job.InnerError = err
        job.UsedCap--
        return
    } else if !censorJob.Safe {
        //job.Error = fmt.Errorf("敏感信息，已拦截")
        seelog.Infof("发现敏感信息，答案：%#v, 原文：%#v, idMsg=%v", censorJob.AuditingResult.FilteredText, job.Message, idMsg)
        job.Error = nil
        job.Message = censorJob.AuditingResult.FilteredText
        job.BadCnt++
        return
    } else {
        seelog.Infof("审核通过，答案：%#v, idMsg=%v", job.Message, idMsg)
        return
    }
}

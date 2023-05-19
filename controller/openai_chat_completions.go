package controller

import (
    "encoding/json"
    "fmt"
    "github.com/cihub/seelog"
    "github.com/gin-gonic/gin"
    "io"
    "muchat-go/app"
    "muchat-go/chatgpt/api_base"
    "muchat-go/chatgpt/client"
    "muchat-go/repo/constants"
    "net/http"
    "strings"
    "time"
)

type OpenaiChatCompletionsParams struct {
    api_base.ChatCompletionsRequest
    PresetPrompt *api_base.PresetPrompt `json:"preset_prompt,omitempty"`
    Slug         string                 `json:"slug"`
}

func (r *Controller) HandleOpenaiChatCompletions(c *gin.Context) {
    acc := r.OpenAiClient.PickAccount()
    if acc == nil {
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrNoAccount))
        return
    }
    params := OpenaiChatCompletionsParams{}
    err := c.BindJSON(&params)
    if err != nil {
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrParams))
        return
    }
    cur, ok := r.CheckUserAndResponse(c, params.Slug, c.ClientIP())
    if !ok {
        return
    }

    id := time.Now().Nanosecond()
    job := &app.JobParam{
        Title:        "",
        TitleIndex:   id,
        Account:      acc,
        Censor:       r.Censor,
        Slug:         cur.User.Slug,
        EnqueueTime:  time.Now(),
        Messages:     params.Messages,
        PresetPrompt: params.PresetPrompt,
    }

    idMsg := fmt.Sprintf("slug=%v, jobId=%v, acc=%v", job.Slug, job.TitleIndex, acc.Email)
    messages := client.PrepareMessages(job.Messages, job.PresetPrompt)

    bs, err := json.Marshal(messages)
    if err != nil {
        seelog.Errorf("json.Marshall messages failed: messages=%#v, err=%v", messages, err)
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrUnknown))
        return
    } else {
        seelog.Infof("正在查询题目（%v）, idMsg=%v", string(bs), idMsg)
    }

    // 审核题目
    var textToCensor string
    textToCensor = messages[len(messages)-1].Content

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
        r.CheckJobAndResponse(c, cur, job, true)
        //c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrUnknown))
        return
    }

    // 发起请求
    chanResp := make(chan api_base.ChatCompletionsStreamResponseChunk)

    go api_base.ChatCompletionsStream(messages, acc.ApiKey, chanResp)
    //go func() {
    //    s := []rune("尊敬的领导、亲爱的同仁们：尊敬的领导、亲爱的同仁们：")
    //    for i := 0; i < 10; i++ {
    //        chunk := api_base.ChatCompletionsStreamResponseChunk{}
    //        chunk.Content.ID = fmt.Sprintf("%v", id)
    //        chunk.Content.Choices = append(chunk.Content.Choices, openai.ChatCompletionStreamChoice{
    //            Delta: openai.ChatCompletionStreamChoiceDelta{
    //                Content: fmt.Sprintf("%v", string(s[i])),
    //                Role:    "assistant",
    //            },
    //        })
    //        chanResp <- chunk
    //        time.Sleep(time.Second)
    //    }
    //    close(chanResp)
    //}()
    //seelog.Infof("stream started")
    text := ""
    entireText := ""
    lastResp := api_base.ChatCompletionsStreamResponseChunk{}
    beg := time.Now()
    contentId := 0
    c.Stream(func(w io.Writer) bool {
        resp, ok := <-chanResp
        if !ok {
            job.UsedCap++

            resp = lastResp
            resp.Content.Choices[0].Delta.Content = text

            // 审核
            //seelog.Infof("传入2:%#v", resp.Content.Choices[0].Delta.Content)
            resp.Content.Choices[0].Delta.Content = r.inspect(job, idMsg, resp.Content.Choices[0].Delta.Content, contentId)
            //seelog.Infof("传出2:%#v", resp.Content.Choices[0].Delta.Content)
            //entireText += resp.Content.Choices[0].Delta.Content

            c.SSEvent("message", resp.Content)
            return false
        }
        if resp.Error != nil {
            seelog.Errorf("api_base.ChatCompletionsStream failed: %v, idMsg=%v", resp.Error, idMsg)
            return false
        }
        //seelog.Infof("chunk: %v, idMsg=%v", resp, idMsg)
        //outputBytes := bytes.NewBufferString(resp.Content)
        //c.Writer.Write(outputBytes.Bytes())
        if strings.Contains(resp.Content.Choices[0].Delta.Content, "\n") ||
            strings.Contains(resp.Content.Choices[0].Delta.Content, "。") ||
            strings.Contains(resp.Content.Choices[0].Delta.Content, "，") {

            entireText += resp.Content.Choices[0].Delta.Content
            resp.Content.Choices[0].Delta.Content = text + resp.Content.Choices[0].Delta.Content

            // 审核
            //seelog.Infof("传入1:%#v", resp.Content.Choices[0].Delta.Content)
            resp.Content.Choices[0].Delta.Content = r.inspect(job, idMsg, resp.Content.Choices[0].Delta.Content, contentId)
            //seelog.Infof("传出1:%#v", resp.Content.Choices[0].Delta.Content)
            lastResp = resp
            contentId++

            c.SSEvent("message", resp.Content)
            text = ""
        } else {
            text += resp.Content.Choices[0].Delta.Content
            entireText += resp.Content.Choices[0].Delta.Content
            lastResp = resp
        }
        return true
    })
    //seelog.Infof("stream done")
    duration := time.Now().Sub(beg).Seconds()
    seelog.Infof("查询成功，耗时: %v, idMsg=%v, content=%#v", duration, idMsg, entireText)

    job.DoneTime = time.Now()
    r.CheckJobAndResponse(c, cur, job, false)
}

func (r *Controller) inspect(job *app.JobParam, jobIdMsg string, content string, contentId int) (result string) {
    jobIdMsg += fmt.Sprintf(", contentId=%v", contentId)
    // 是否审核答案
    if r.Config.Mock.Enabled || r.Config.CensorEnabled == false || strings.Contains(job.Message, "已拦截") {
        // 不审核
        //seelog.Infof("不审核，答案：%#v, jobIdMsg=%v", job.Message, jobIdMsg)
        result = content
        return
    }
    // 审核
    job.Censor.AddJob(contentId, content)
    censorJob := job.Censor.WaitForJob(contentId)
    if censorJob.Err != nil {
        seelog.Errorf("审核出错-1: %v, jobIdMsg=%v", censorJob.Err, jobIdMsg)
        return
    } else if !censorJob.Safe {
        //job.Error = fmt.Errorf("敏感信息，已拦截")
        seelog.Warnf("发现敏感信息，答案：%#v, 原文：%#v, jobIdMsg=%v", censorJob.AuditingResult.FilteredText, content, jobIdMsg)
        //job.Error = nil
        //job.Message = censorJob.AuditingResult.FilteredText
        job.BadCnt++
        result = censorJob.AuditingResult.FilteredText
        return
    } else {
        //seelog.Infof("审核通过，答案：%#v, jobIdMsg=%v", content, jobIdMsg)
        result = content
        return
    }
}

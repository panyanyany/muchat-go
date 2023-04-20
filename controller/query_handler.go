package controller

import (
    "errors"
    "github.com/cihub/seelog"
    "github.com/gin-gonic/gin"
    "muchat-go/app"
    "muchat-go/chatgpt/api_base"
    "muchat-go/models"
    "muchat-go/repo/constants"
    "net/http"
    "time"
)

func (r *Controller) HandleQuery(c *gin.Context) {
    acc := r.OpenAiClient.PickAccount()
    id := time.Now().Nanosecond()
    params := QueryParams{}
    err := c.BindJSON(&params)
    if err != nil || params.Question == "" {
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrParams))
        return
    }
    cur := r.CheckUser(c, params.Slug, c.ClientIP())
    if !cur.Auth {
        seelog.Warnf("无权限: %v", params.Slug)
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrUnAuth))
        return
    }
    if cur.ReachCap {
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrReachCap))
        return
    }
    if cur.Expired {
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrExpired))
        return
    }

    if acc == nil {
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrNoAccount))
        return
    }

    r.Runner.AddJob(id, params.Question, acc, r.Censor, cur.User.Slug, params.Messages)
    job := r.Runner.WaitForJob(id)
    if job == nil {
        seelog.Errorf("r.Runner.WaitForJob(id) returns nil, id=%v", id)
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrUnknown))
    } else {
        r.CheckInnerError(job)
        cur.User.Usage += job.UsedCap
        cur.User.BadCnt += job.BadCnt
        cur.User.Referer = c.GetHeader("referer")
        if cur.User.FirstTime.Year() == 1 {
            maxDays := cur.User.MaxDays
            if maxDays == 0 {
                maxDays = 30
            }
            cur.User.FirstTime = time.Now()
            cur.User.ExpiresAt = time.Now().Add(time.Duration(maxDays) * 24 * time.Hour)
            cur.User.FirstIp = c.ClientIP()
        }

        if cur.User.FirstIp == "" {
            cur.User.FirstIp = c.ClientIP()
        }

        err = r.Db.Model(cur.User).Updates(models.MuchatUser{
            FirstTime: cur.User.FirstTime,
            ExpiresAt: cur.User.ExpiresAt,
            Usage:     cur.User.Usage,
            BadCnt:    cur.User.BadCnt,
            Referer:   cur.User.Referer,
            FirstIp:   cur.User.FirstIp,
        }).Error
        if err != nil {
            seelog.Errorf("save user failed: %v, id=%v, userId=%v", err, id, cur.User.ID)
            err = nil
        }

        if job.Error != nil {
            respData := constants.GetResponseBody(job.ErrCode)
            respData["data"] = map[string]interface{}{
                "content": job.Error.Error(),
            }
            seelog.Errorf("r.Runner.WaitForJob(id) returns error, err=%v, id=%v", job.Error, id)
            c.JSON(http.StatusOK, respData)
        } else {
            respData := constants.GetResponseBody(constants.CodeOk)
            respData["data"] = map[string]interface{}{
                "answer": job.Message,
            }
            c.JSON(http.StatusOK, respData)
        }
    }
}

type QueryParams struct {
    Question string                 `json:"question"`
    Slug     string                 `json:"slug"`
    Messages []api_base.ChatMessage `json:"messages"`
}

func (r *Controller) CheckInnerError(job *app.JobParam) {
    if job.InnerError == nil {
        return
    }
    seelog.Errorf("InnerError: %v, %v, %v", job.InnerError, job.InnerErrorMsg, job.GetIdStr())
    var err error
    err = api_base.ExtractError([]byte(job.InnerErrorMsg))
    // 余额不足
    if errors.Is(err, api_base.ErrNoBalance) {
        err = r.Db.Where(models.OpenAiAccount{Email: job.Account.Email}).Updates(models.OpenAiAccount{Status: models.OpenAiAccountStatusPaused}).Error
        if err != nil {
            seelog.Errorf("update account status failed-0, err=%v", err)
            return
        }
    } else if errors.Is(err, api_base.ErrBaned) { // 被封禁
        err = r.Db.Where(models.OpenAiAccount{Email: job.Account.Email}).Updates(models.OpenAiAccount{Status: models.OpenAiAccountStatusDisabled}).Error
        if err != nil {
            seelog.Errorf("update account status failed-1, err=%v", err)
            return
        }
    } else { // 未知错误
        //err = r.Db.Where(models.OpenAiAccount{Email: job.Account.Email}).Updates(models.OpenAiAccount{Status: models.OpenAiAccountStatusUnknown}).Error
        //if err != nil {
        //	seelog.Errorf("update account status failed-2, err=%v", err)
        //	return
        //}
    }
    r.OpenAiClient.LoadAccounts()
}

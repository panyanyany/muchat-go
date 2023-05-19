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

func (r *Controller) CheckUser(c *gin.Context, slug string, ip string) (cur models.CheckUserResult) {
    var existed bool
    if slug != "" {
        cur, existed = models.CheckUser(r.Db, slug)
        return cur
    }

    guestCfg := r.Config.Guests.GetDomainConfig(c.GetHeader("referer"))
    if guestCfg == nil || !guestCfg.Enabled {
        return cur
    }

    slug = "ip-" + ip

    cur, existed = models.CheckUser(r.Db, slug)
    // 插入新用户
    if !existed {
        md := models.MuchatUser{
            Slug:      slug,
            FirstTime: time.Now(),
            Usage:     0,
            BadCnt:    0,
            MaxUsage:  guestCfg.MaxUsage,
            MaxDays:   1,
            ExpiresAt: time.Now().Add(time.Hour * 24),
            Name:      "",
            FirstIp:   ip,
        }
        err := r.Db.Create(&md).Error
        if err != nil {
            seelog.Errorf("insert new ip user failed, err=%v, ip=%v", err, ip)
            return
        }
        cur.User = &md
        cur.Auth = true
        cur.ReachCap = guestCfg.MaxUsage == 0 // 为了测试，会把 max_usage 设置为 0
    }
    return
}

func (r *Controller) CheckUserAndResponse(c *gin.Context, slug string, ip string) (cur models.CheckUserResult, ok bool) {
    cur = r.CheckUser(c, slug, ip)
    if !cur.Auth {
        seelog.Warnf("无权限: %v", slug)
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
    ok = true
    return
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

func (r *Controller) CheckJobAndResponse(c *gin.Context, cur models.CheckUserResult, job *app.JobParam, returnJson bool) {
    var err error
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
        seelog.Errorf("save user failed: %v, id=%v, userId=%v", err, job.TitleIndex, cur.User.ID)
        err = nil
    }

    if job.Error != nil {
        respData := constants.GetResponseBody(job.ErrCode)
        respData["data"] = map[string]interface{}{
            "content": job.Error.Error(),
        }
        seelog.Errorf("r.Runner.WaitForJob(id) returns error, err=%v, id=%v", job.Error, job.TitleIndex)
        if returnJson {
            c.JSON(http.StatusInternalServerError, respData)
        }
    } else {
        respData := constants.GetResponseBody(constants.CodeOk)
        respData["data"] = map[string]interface{}{
            "answer": job.Message,
        }
        if returnJson {
            c.JSON(http.StatusOK, respData)
        }
    }
}

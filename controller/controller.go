package controller

import (
    "github.com/cihub/seelog"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "muchat-go/app"
    "muchat-go/chatgpt/client"
    "muchat-go/config"
    "muchat-go/models"
    "muchat-go/repo/censor"
    "muchat-go/repo/constants"
    "net/http"
    "time"
)

type Controller struct {
    Db           *gorm.DB
    Runner       *app.Runner
    Censor       *censor.Censor
    OpenAiClient *client.Client
    Config       *config.Configuration
}

func (r *Controller) HandleVersion(c *gin.Context) {
    var mdVer models.Version
    err := r.Db.First(&mdVer).Error
    if err != nil {
        seelog.Errorf("get version failed, err=%v", err)
    }

    resp := constants.GetResponseBody(constants.CodeOk)
    resp["data"] = mdVer
    c.JSON(http.StatusOK, resp)
}

func (r *Controller) HandleCheck(c *gin.Context) {
    userSlug := c.Query("slug")
    resp := constants.GetResponseBody(constants.CodeOk)
    if userSlug == "" {
        resp["data"] = map[string]bool{
            "auth": false,
        }
    } else {
        cur := r.CheckUser(c, userSlug, c.ClientIP())
        resp["data"] = map[string]bool{
            "auth":      cur.Auth,
            "reach_cap": cur.ReachCap,
            "expired":   cur.Expired,
        }
    }
    c.JSON(http.StatusOK, resp)
}

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

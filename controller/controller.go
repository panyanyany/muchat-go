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

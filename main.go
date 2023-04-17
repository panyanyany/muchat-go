package main

import (
    "github.com/cihub/seelog"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "go_another_chatgpt/app"
    "go_another_chatgpt/chatgpt/client"
    "go_another_chatgpt/config"
    "go_another_chatgpt/controller"
    "go_another_chatgpt/models"
    "go_another_chatgpt/repo/censor"
    "go_another_chatgpt/repo/censor/adapters"
    "go_another_chatgpt/util/log_util"
    "go_another_chatgpt/util/thread_util"
    "net/http"
)

func main() {
    log_util.SetupSeelog()
    defer seelog.Flush()

    defer seelog.Info("退出程序")

    accounts := config.LoadAccounts()
    seelog.Infof("loaded accounts: %v", len(accounts))

    cfg := config.LoadConfig()
    db := models.InitDb(cfg.Db)

    // open ai
    openAiClient := client.NewClient(db)
    openAiClient.Concurrency = cfg.OpenAiAccountConfig.Concurrency
    openAiClient.LoadAccountIntoDb()
    openAiClient.LoadAccounts()
    go func() {
        openAiClient.RefreshUsage()
    }()

    funcTicker := thread_util.NewFuncTicker(cfg.OpenAiAccountConfig.QueryInterval)
    funcTicker.Handler = func() {
        openAiClient.LoadAccountIntoDb()
        openAiClient.LoadAccounts()
        openAiClient.RefreshUsage()
    }
    funcTicker.Start()

    runner := app.NewRunner(cfg)
    runner.Start()

    // 审核
    tencentCi := adapters.NewTencentCi(cfg.TencentCos.SecretId,
        cfg.TencentCos.SecretKey,
        cfg.TencentCos.BucketUrl,
        cfg.TencentCos.ServiceUrl,
        cfg.TencentCos.CiUrl)
    _ = tencentCi
    baiduAi := adapters.NewBaidu(cfg.BaiduAi.AppKey, cfg.BaiduAi.SecretKey)
    _ = baiduAi
    localFilter := adapters.NewLocalDirtyFilter()

    myCensor := censor.NewCensor(localFilter, 3)
    myCensor.Start()

    myController := controller.Controller{
        Db:           db,
        Runner:       runner,
        Censor:       myCensor,
        OpenAiClient: openAiClient,
        Config:       cfg,
    }

    router := gin.Default()
    router.Use(cors.Default())
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
    router.POST("/api/query", myController.HandleQuery)
    router.GET("/api/check", myController.HandleCheck)
    router.GET("/api/client/version", myController.HandleVersion)
    router.Run(cfg.Listen) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

    funcTicker.Stop()

    //runner.WaitForEmptyJob()
    //seelog.Infof("所有题目完成分发")
    //runner.Stop()
}

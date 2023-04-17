package constants

import "github.com/gin-gonic/gin"

const (
    CodeOk = 0
    // 服务端错误
    CodeErrUnknown   = 1001
    CodeErrNoAccount = 1002
    CodeErrSensitive = 1003
    CodeErrGptError  = 1004
    // 客户端错误
    CodeErrParams               = 2001
    CodeErrUnAuth               = 2002
    CodeErrReachCap             = 2003
    CodeErrInvalidRetailAccount = 2004
    CodeErrExpired              = 2005
)

var CodeMessages = map[int]string{
    CodeOk:                      "成功",
    CodeErrUnknown:              "服务端错误",
    CodeErrParams:               "参数错误",
    CodeErrUnAuth:               "无权限",
    CodeErrReachCap:             "额度用完",
    CodeErrInvalidRetailAccount: "无效账号",
    CodeErrNoAccount:            "系统临时维护中",
    CodeErrGptError:             "ChatGPT接口异常",
    CodeErrSensitive:            "敏感信息",
}

func GetResponseBody(code int) gin.H {
    return gin.H{
        "code":    code,
        "message": CodeMessages[code],
    }
}

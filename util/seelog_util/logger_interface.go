package seelog_util

import "github.com/cihub/seelog"

type SimpleLoggerInterface interface {
    Debug(context seelog.LogContextInterface, s string)
    Info(context seelog.LogContextInterface, s string)
    Warn(context seelog.LogContextInterface, s string)
    Error(context seelog.LogContextInterface, s string)
}

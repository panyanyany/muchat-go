package seelog_util

import "github.com/cihub/seelog"

type SomeCustomReceiver struct {
    Logger SimpleLoggerInterface
}

func (ar *SomeCustomReceiver) ReceiveMessage(message string, level seelog.LogLevel, context seelog.LogContextInterface) error {
    switch level {
    case seelog.TraceLvl:
        fallthrough
    case seelog.DebugLvl:
        ar.Logger.Debug(context, message)
    case seelog.InfoLvl:
        ar.Logger.Info(context, message)
    case seelog.WarnLvl:
        ar.Logger.Warn(context, message)
    case seelog.ErrorLvl:
        fallthrough
    case seelog.CriticalLvl:
        ar.Logger.Error(context, message)
    }
    return nil
}

/* NOTE: NOT called when LoggerFromCustomReceiver is used */
func (ar *SomeCustomReceiver) AfterParse(initArgs seelog.CustomReceiverInitArgs) error {
    return nil
}
func (ar *SomeCustomReceiver) Flush() {

}
func (ar *SomeCustomReceiver) Close() error {
    return nil
}

package goreq_util

import (
	"github.com/cihub/seelog"
)

type Logger struct {
	Prefix string
	Logger seelog.LoggerInterface
}

func (r *Logger) SetPrefix(prefix string) {
	r.Prefix = prefix
}
func (r *Logger) Printf(format string, v ...interface{}) {
	r.Logger.Debugf(format, v...)
}
func (r *Logger) Println(v ...interface{}) {
	r.Logger.Debug(v...)
}

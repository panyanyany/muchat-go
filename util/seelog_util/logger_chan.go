package seelog_util

import (
    "fmt"
    "go_another_chatgpt/util/time_util"
    "io"
    "time"

    "github.com/cihub/seelog"
)

type ChanLogger struct {
    ChText   chan string
    Filename string
    File     io.Writer
}

func (r *ChanLogger) AppendFile(text string) (err error) {
    r.File.Write([]byte(text))
    //err = file_util.AppendFile(r.Filename, text, nil)
    //var bs []byte
    //bs, err = ioutil.ReadFile(r.Filename)
    //lines := strings.Split(string(bs), "\n")
    //maxLine := 20000
    //if len(lines) >= maxLine {
    //	lines = lines[len(lines)-maxLine-1:]
    //	ioutil.WriteFile(r.Filename, []byte(strings.Join(lines, "\n")), 0644)
    //}
    return
}

func (r *ChanLogger) Debug(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [调试] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
}
func (r *ChanLogger) Info(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [信息] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
    text = fmt.Sprintf("%s [信息] %s\n", time.Now().Format(time_util.DefaultFormat), s)
    r.ChText <- text
}
func (r *ChanLogger) Error(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [错误] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
    text = fmt.Sprintf("%s [!错误!] %s\n", time.Now().Format(time_util.DefaultFormat), s)
    r.ChText <- text
}
func (r *ChanLogger) Raw(s string) {
    r.AppendFile(s)
    r.ChText <- s
}

package seelog_util

import (
    "fmt"
    "go_another_chatgpt/util/time_util"
    "io"
    "time"

    "github.com/cihub/seelog"
)

type SimpleLogger struct {
    Filename string
    File     io.Writer
}

func (r *SimpleLogger) AppendFile(text string) (err error) {
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

func (r *SimpleLogger) Debug(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [调试] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
}
func (r *SimpleLogger) Info(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [信息] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
    text = fmt.Sprintf("%s [信息] %s\n", time.Now().Format(time_util.DefaultFormat), s)
}
func (r *SimpleLogger) Warn(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [警告] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
    text = fmt.Sprintf("%s [!警告!] %s\n", time.Now().Format(time_util.DefaultFormat), s)
}
func (r *SimpleLogger) Error(context seelog.LogContextInterface, s string) {
    var text string
    text = fmt.Sprintf("%s %v|> %s [错误] %s\n",
        context.FullPath(),
        context.Line(),
        time.Now().Format(time_util.DefaultFormat), s)
    r.AppendFile(text)
    text = fmt.Sprintf("%s [!错误!] %s\n", time.Now().Format(time_util.DefaultFormat), s)
}
func (r *SimpleLogger) Raw(s string) {
    r.AppendFile(s)
}

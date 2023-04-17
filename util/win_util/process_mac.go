//go:build darwin
// +build darwin

package win_util

import (
    "github.com/cihub/seelog"
    "os/exec"
    "path/filepath"
    "strings"
)

// error is nil on success
func Reboot2() error {
    return nil
}
func Reboot() error {
    args := strings.Split("shutdown /r /f /t 0", " ")
    err := exec.Command(args[0], args[1:]...).Run()
    if err != nil {
        seelog.Errorf("reboot: %v", err)
    }
    return err
}
func TaskKill(exeName string) {
    args := strings.Split("taskkill /f /im "+exeName, " ")
    err := exec.Command(args[0], args[1:]...).Run()
    if err != nil {
        seelog.Error(err)
    }
}

// TaskKillPath 根据绝对路径来杀死进程, 路径最好不要有空格
func TaskKillPath(exeFullPath string) {
    args := []string{}
    args = append(args, strings.Split("wmic process where", " ")...)
    args = append(args, "ExecutablePath='"+strings.ReplaceAll(filepath.FromSlash(exeFullPath), "\\", "\\\\")+"'")
    args = append(args, "delete")
    seelog.Debugf("taskkill2: %#v", args)
    err := exec.Command(args[0], args[1:]...).Run()
    if err != nil {
        seelog.Error(err)
    }
}

// TaskKillCmdLine 根据命令行匹配来杀死进程
func TaskKillCmdLine(exeName string, partOfCmdLine string) {
    args := []string{}
    args = append(args, strings.Split("wmic Path win32_process Where", " ")...)
    args = append(args, "Name='"+exeName+"' AND CommandLine Like '%"+strings.ReplaceAll(filepath.FromSlash(partOfCmdLine), "\\", "\\\\")+"%'")
    //args = append(args, "Name='"+exeName+"' AND CommandLine Like '%"+partOfCmdLine+"%'")
    args = append(args, strings.Split("Call Terminate", " ")...)
    seelog.Debugf("TaskKillCmdLine: %#v", args)
    err := exec.Command(args[0], args[1:]...).Run()
    if err != nil {
        seelog.Error(err)
    }
}

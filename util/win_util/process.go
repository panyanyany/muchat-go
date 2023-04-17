// +build windows

package win_util

import (
    "github.com/cihub/seelog"
    "os/exec"
    "path/filepath"
    "strings"
    "syscall"
    "unsafe"
)

// error is nil on success
func Reboot2() error {

    user32 := syscall.MustLoadDLL("user32")
    defer user32.Release()

    kernel32 := syscall.MustLoadDLL("kernel32")
    defer user32.Release()

    advapi32 := syscall.MustLoadDLL("advapi32")
    defer advapi32.Release()

    ExitWindowsEx := user32.MustFindProc("ExitWindowsEx")
    GetCurrentProcess := kernel32.MustFindProc("GetCurrentProcess")
    GetLastError := kernel32.MustFindProc("GetLastError")
    OpenProdcessToken := advapi32.MustFindProc("OpenProcessToken")
    LookupPrivilegeValue := advapi32.MustFindProc("LookupPrivilegeValueW")
    AdjustTokenPrivileges := advapi32.MustFindProc("AdjustTokenPrivileges")

    currentProcess, _, _ := GetCurrentProcess.Call()

    const tokenAdjustPrivileges = 0x0020
    const tokenQuery = 0x0008
    var hToken uintptr

    result, _, err := OpenProdcessToken.Call(currentProcess, tokenAdjustPrivileges|tokenQuery, uintptr(unsafe.Pointer(&hToken)))
    if result != 1 {
        seelog.Error("OpenProcessToken(): ", result, " err: ", err)
        return err
    }
    //fmt.Println("hToken: ", hToken)

    const SeShutdownName = "SeShutdownPrivilege"

    type Luid struct {
        lowPart  uint32 // DWORD
        highPart int32  // long
    }
    type LuidAndAttributes struct {
        luid       Luid   // LUID
        attributes uint32 // DWORD
    }

    type TokenPrivileges struct {
        privilegeCount uint32 // DWORD
        privileges     [1]LuidAndAttributes
    }

    var tkp TokenPrivileges

    result, _, err = LookupPrivilegeValue.Call(uintptr(0), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(SeShutdownName))), uintptr(unsafe.Pointer(&(tkp.privileges[0].luid))))
    if result != 1 {
        seelog.Error("LookupPrivilegeValue(): ", result, " err: ", err)
        return err
    }
    //fmt.Println("LookupPrivilegeValue luid: ", tkp.privileges[0].luid)

    const SePrivilegeEnabled uint32 = 0x00000002

    tkp.privilegeCount = 1
    tkp.privileges[0].attributes = SePrivilegeEnabled

    result, _, err = AdjustTokenPrivileges.Call(hToken, 0, uintptr(unsafe.Pointer(&tkp)), 0, uintptr(0), 0)
    if result != 1 {
        seelog.Error("AdjustTokenPrivileges() ", result, " err: ", err)
        return err
    }

    result, _, _ = GetLastError.Call()
    if result != 0 {
        seelog.Error("GetLastError() ", result)
        return err
    }

    const ewxForceIfHung = 0x00000010
    const ewxReboot = 0x00000002
    const ewxShutdown = 0x00000001
    const shutdownReasonMajorSoftware = 0x00030000

    result, _, err = ExitWindowsEx.Call(ewxReboot|ewxForceIfHung, shutdownReasonMajorSoftware)
    if result != 1 {
        seelog.Error("Failed to initiate reboot:", err)
        return err
    }

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

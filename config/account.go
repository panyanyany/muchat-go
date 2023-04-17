package config

import (
    "fmt"
    "github.com/cihub/seelog"
    "io/ioutil"
    "os"
    "strings"
)

func LoadAccounts() (accounts []Account) {
    accounts = make([]Account, 0, 10)

    accFile := "账户.txt"
    _, err := os.Stat(accFile)
    if err != nil {
        if !os.IsExist(err) {
            return
        }
        panic(fmt.Errorf("get account file info failed, err=%w", err))
    }

    bs, err := ioutil.ReadFile(accFile)
    if err != nil {
        panic(fmt.Errorf("账户列表读取出错：%w", err))
    }
    for _, line := range strings.Split(string(bs), "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        pieces := strings.Split(line, ",")
        acc := Account{
            Email:    pieces[0],
            Password: pieces[1],
            ApiKey:   pieces[2],
        }
        if len(acc.ApiKey) < 5 {
            seelog.Warnf("没有 API_KEY: %v", acc.Email)
            continue
        }

        accounts = append(accounts, acc)
    }

    return
}

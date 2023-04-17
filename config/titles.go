package config

import (
    "fmt"
    "io/ioutil"
    "strings"
)

func LoadTitles() (titles []string) {
    titles = make([]string, 0, 10)

    accFile := "题目.txt"
    bs, err := ioutil.ReadFile(accFile)
    if err != nil {
        panic(fmt.Errorf("题目读取出错：%w", err))
    }
    for _, line := range strings.Split(string(bs), "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        titles = append(titles, line)
    }

    return
}

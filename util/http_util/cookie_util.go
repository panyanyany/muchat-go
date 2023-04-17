package http_util

import (
    "fmt"
    "net/http"
    "strings"
)

func CookiesToString(cookies []*http.Cookie) (str string) {
    list := []string{}
    for _, cookie := range cookies {
        list = append(list, fmt.Sprintf("%v=%+v", cookie.Name, cookie.Value))
    }
    str = strings.Join(list, "; ")
    return
}

func PickCookies(cookies []*http.Cookie, names []string) (newCookies []*http.Cookie) {
    newCookies = []*http.Cookie{}
    nameMap := map[string]bool{}
    for _, name := range names {
        nameMap[name] = true
    }
    for _, cookie := range cookies {
        _, found := nameMap[cookie.Name]
        if !found {
            continue
        }
        newCookies = append(newCookies, cookie)
    }
    return
}

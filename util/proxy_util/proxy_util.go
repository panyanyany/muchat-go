package proxy_util

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

func IsValidProxy(proxy string) bool {
	parts := strings.Split(proxy, ":")
	if len(parts) != 2 {
		return false
	}
	if net.ParseIP(parts[0]) == nil {
		return false
	}
	return true
}

func FilterProxies(proxies []string) (newProxies []string, invalidProxies []string) {
	newProxies = make([]string, 0)
	invalidProxies = make([]string, 0)

	proxyChan := make(chan string)
	invalidProxyChan := make(chan string)

	wg := sync.WaitGroup{}
	running := true

	go func() {
		for running {
			select {
			case proxy := <-proxyChan:
				newProxies = append(newProxies, proxy)
			case proxy := <-invalidProxyChan:
				invalidProxies = append(invalidProxies, proxy)
			}
		}
	}()

	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			sa := gorequest.New()
			sa.Get("https://www.baidu.com")
			sa.Proxy("http://" + proxy)
			sa.Header = http.Header{
				"User-Agent": []string{"Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.0)"},
			}
			sa.Timeout(10 * time.Second)
			resp, _, errs := sa.End()
			if len(errs) > 0 {
				invalidProxyChan <- proxy
				return
			}
			if resp.StatusCode != 200 {
				invalidProxyChan <- proxy
				return
			}
			proxyChan <- proxy
		}(proxy)
	}

	wg.Wait()
	running = false
	return
}


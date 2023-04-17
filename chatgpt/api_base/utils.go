package api_base

import (
    "errors"
    "fmt"
    "github.com/parnurzeal/gorequest"
    "os"
    "time"
)

func HttpPost(action string, apiKey string, requestData []byte) (body []byte, err error) {
    request := gorequest.New()
    request.Post(BASEURL + action).
        Timeout(60 * time.Second).
        //Proxy("http://127.0.0.1:7890").
        Send(string(requestData))

    request.AppendHeader("Content-Type", "application/json")
    request.AppendHeader("Authorization", "Bearer "+apiKey)

    resp, body, errs := request.
        EndBytes()

    if len(errs) > 0 {
        if os.IsTimeout(errs[0]) {
            return body, fmt.Errorf("completions timeout, err=%w", errs[0])
        }
        return body, errors.New(fmt.Sprintf("gtp api has errors: %v", errs))
    }
    if resp.StatusCode != 200 {
        //err = ExtractError(resp, body)
        err = fmt.Errorf("status code = %v", resp.StatusCode)
        return body, err
    }
    return body, nil
}

func HttpGet(url string, apiKey string) (body []byte, err error) {
    request := gorequest.New()
    request.Get(url).
        //Proxy("http://127.0.0.1:7890").
        Timeout(60 * time.Second)

    request.AppendHeader("Content-Type", "application/json")
    request.AppendHeader("Authorization", "Bearer "+apiKey)

    resp, body, errs := request.
        EndBytes()

    if len(errs) > 0 {
        return body, errors.New(fmt.Sprintf("gtp api has errors: %v", errs))
    }
    if resp.StatusCode != 200 {
        return body, errors.New(fmt.Sprintf("gtp api status code not equals 200,code is %d, %v", resp.StatusCode, string(body)))
    }
    return body, nil
}

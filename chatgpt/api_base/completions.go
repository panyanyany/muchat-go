package api_base

import (
    "encoding/json"
    "github.com/cihub/seelog"
    "muchat-go/config"
)

const BASEURL = "https://api.openai.com/v1/"

// CompletionsResponseBody 请求体
type CompletionsResponseBody struct {
    ID      string                 `json:"id"`
    Object  string                 `json:"object"`
    Created int                    `json:"created"`
    Model   string                 `json:"model"`
    Choices []ChoiceItem           `json:"choices"`
    Usage   map[string]interface{} `json:"usage"`
}

type ErrorRespBody struct {
    Error struct {
        Message string      `json:"message"`
        Type    string      `json:"type"`
        Param   interface{} `json:"param"`
        Code    interface{} `json:"code"`
    } `json:"error"`
}

type ChoiceItem struct {
    Text         string `json:"text"`
    Index        int    `json:"index"`
    Logprobs     int    `json:"logprobs"`
    FinishReason string `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
    Model            string  `json:"model"`
    Prompt           string  `json:"prompt"`
    MaxTokens        uint    `json:"max_tokens"`
    Temperature      float64 `json:"temperature"`
    TopP             int     `json:"top_p"`
    FrequencyPenalty int     `json:"frequency_penalty"`
    PresencePenalty  int     `json:"presence_penalty"`
    User             string  `json:"user"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string, apiKey string, cfg *config.GptConfig, user string) (string, error) {
    //cfg := config.LoadConfig()
    requestBody := ChatGPTRequestBody{
        Model:            cfg.Model,
        Prompt:           msg,
        MaxTokens:        cfg.MaxTokens,
        Temperature:      cfg.Temperature,
        TopP:             1,
        FrequencyPenalty: 0,
        PresencePenalty:  0,
        User:             user,
    }
    requestData, err := json.Marshal(requestBody)

    if err != nil {
        return "", err
    }

    seelog.Debugf("reqData: %v", string(requestData))

    var body []byte
    body, err = HttpPost("completions", apiKey, requestData)
    if err != nil {
        return string(body), err
    }

    seelog.Debugf("response: %v", string(body))

    gptResponseBody := &CompletionsResponseBody{}
    //log.Println(string(body))
    err = json.Unmarshal(body, gptResponseBody)
    if err != nil {
        return "", err
    }

    var reply string
    if len(gptResponseBody.Choices) > 0 {
        reply = gptResponseBody.Choices[0].Text
    }
    //seelog.Debugf("gpt response text: %s ", reply)
    return reply, nil
}

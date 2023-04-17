package api_base

import (
    "encoding/json"
    "github.com/cihub/seelog"
    "go_another_chatgpt/config"
)

func ChatCompletions(ms []ChatMessage, apiKey string, cfg *config.Configuration, user string) (string, error) {
    //cfg := config.LoadConfig()
    requestBody := map[string]interface{}{
        "model":    "gpt-3.5-turbo",
        "messages": ms,
    }
    requestData, err := json.Marshal(requestBody)

    if err != nil {
        return "", err
    }

    seelog.Debugf("reqData: %v", string(requestData))

    var body []byte
    body, err = HttpPost("chat/completions", apiKey, requestData)
    if err != nil {
        return string(body), err
    }

    seelog.Debugf("response: %v", string(body))

    gptResponseBody := &ChatMessageResponseBody{}
    //log.Println(string(body))
    err = json.Unmarshal(body, gptResponseBody)
    if err != nil {
        return "", err
    }
    var reply string
    if len(gptResponseBody.Choices) > 0 {
        reply = gptResponseBody.Choices[0].Message.Content
    }
    //seelog.Debugf("gpt response text: %s ", reply)
    return reply, nil
}

type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatMessageResponseBody struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int    `json:"created"`
    Choices []struct {
        Index   int `json:"index"`
        Message struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
        FinishReason string `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

package api_base

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/cihub/seelog"
    "github.com/sashabaranov/go-openai"
    "io"
    "muchat-go/config"
    "time"
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

type ChatCompletionsStreamResponseChunk struct {
    Content openai.ChatCompletionStreamResponse
    Error   error
}

type RecvResult struct {
    Response openai.ChatCompletionStreamResponse
    Error    error
}

func ChatCompletionsStream(ms []ChatMessage, apiKey string, chanResp chan ChatCompletionsStreamResponseChunk) {
    c := openai.NewClient(apiKey)
    ctx := context.Background()

    var ms2 []openai.ChatCompletionMessage
    for _, m := range ms {
        ms2 = append(ms2, openai.ChatCompletionMessage{
            Role:    m.Role,
            Content: m.Content,
        })
    }

    req := openai.ChatCompletionRequest{
        Model:    openai.GPT3Dot5Turbo,
        Messages: ms2,
        Stream:   true,
    }
    stream, err := c.CreateChatCompletionStream(ctx, req)
    if err != nil {
        err = fmt.Errorf("ChatCompletionStream error: %v", err)
        chanResp <- ChatCompletionsStreamResponseChunk{Error: err}
        return
    }
    defer stream.Close()

    for {
        recvCtx, _ := context.WithTimeout(ctx, time.Second*10)
        chRes := make(chan *RecvResult)
        go func() {
            var res RecvResult
            var err error
            var resp openai.ChatCompletionStreamResponse
            resp, err = stream.Recv()
            res.Error = err
            res.Response = resp
            chRes <- &res
        }()

        var response openai.ChatCompletionStreamResponse

        select {
        case res := <-chRes:
            if res.Error != nil {
                err = res.Error
            } else {
                response = res.Response
            }
        case <-recvCtx.Done():
            err = fmt.Errorf("stream timeout: %v", recvCtx.Err())
        }
        if errors.Is(err, io.EOF) {
            //seelog.Info("Stream finished")
            close(chanResp)
            return
        }

        if err != nil {
            err = fmt.Errorf("stream error: %v", err)
            chanResp <- ChatCompletionsStreamResponseChunk{Error: err}
            return
        }

        chanResp <- ChatCompletionsStreamResponseChunk{Content: response}
    }
}

type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type PresetPrompt struct {
    Act    string `json:"act"`
    Prompt string `json:"prompt"`
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

type ChatCompletionsRequest struct {
    Model    string        `json:"model"`
    Messages []ChatMessage `json:"messages"`

    Temperature      float32        `json:"temperature,omitempty"`
    TopP             float32        `json:"top_p,omitempty"`
    N                int            `json:"n,omitempty"`
    Stream           bool           `json:"stream,omitempty"`
    Stop             []string       `json:"stop,omitempty"`
    MaxTokens        int            `json:"max_tokens,omitempty"`
    PresencePenalty  float32        `json:"presence_penalty,omitempty"`
    FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
    LogitBias        map[string]int `json:"logit_bias,omitempty"`
    User             string         `json:"user,omitempty"`
}

package client

import (
    "github.com/panyanyany/go-util/utf8_util"
    "muchat-go/chatgpt/api_base"
)

func PrepareMessages(messages []api_base.ChatMessage, presetPrompt *api_base.PresetPrompt) []api_base.ChatMessage {

    // 大概不能超过这个数
    maxLen := 1500 // 3500最极限，但会导致回复很短，
    textCnt := 0
    messages2 := []api_base.ChatMessage{}

    // 插入预设
    if presetPrompt != nil {
        maxLen -= utf8_util.Len(presetPrompt.Prompt)
    }
    for _, m := range messages {
        curTextCnt := textCnt + utf8_util.Len(m.Content)
        if curTextCnt > maxLen {
            if textCnt == 0 { // 唯一的一个问题太长
                //m.Content = utf8_util.Substr(m.Content, len(m.Content)-maxLen, maxLen)
                messages2 = append(messages2, m)
            }
            break
        }
        textCnt = curTextCnt
        messages2 = append(messages2, m)
    }
    messages = messages2

    // 插入预设
    if presetPrompt != nil {
        messages = append([]api_base.ChatMessage{
            api_base.ChatMessage{
                Role:    "system",
                Content: presetPrompt.Prompt,
            },
        }, messages...)
    }
    return messages
}

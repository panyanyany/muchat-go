package controller

import (
    "github.com/cihub/seelog"
    "github.com/gin-gonic/gin"
    "muchat-go/chatgpt/api_base"
    "muchat-go/repo/constants"
    "net/http"
    "time"
)

func (r *Controller) HandleQuery(c *gin.Context) {
    acc := r.OpenAiClient.PickAccount()
    if acc == nil {
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrNoAccount))
        return
    }
    id := time.Now().Nanosecond()
    params := QueryParams{}
    err := c.BindJSON(&params)
    if err != nil || params.Question == "" {
        c.JSON(http.StatusBadRequest, constants.GetResponseBody(constants.CodeErrParams))
        return
    }
    cur, ok := r.CheckUserAndResponse(c, params.Slug, c.ClientIP())
    if !ok {
        return
    }

    r.Runner.AddJob(id, params.Question, acc, r.Censor, cur.User.Slug, params.Messages, params.PresetPrompt)
    job := r.Runner.WaitForJob(id)
    if job == nil {
        seelog.Errorf("r.Runner.WaitForJob(id) returns nil, id=%v", id)
        c.JSON(http.StatusInternalServerError, constants.GetResponseBody(constants.CodeErrUnknown))
    } else {
        r.CheckJobAndResponse(c, cur, job, true)
    }
}

type QueryParams struct {
    Question     string                 `json:"question"`
    Slug         string                 `json:"slug"`
    Messages     []api_base.ChatMessage `json:"messages"`
    PresetPrompt *api_base.PresetPrompt `json:"preset_prompt,omitempty"`
}

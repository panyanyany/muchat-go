package api_base

import (
    "encoding/json"
    "github.com/cihub/seelog"
)

type CreditGrantsResponseBody struct {
    Object         string  `json:"object"`
    TotalGranted   float64 `json:"total_granted"`
    TotalUsed      float64 `json:"total_used"`
    TotalAvailable float64 `json:"total_available"`
    Grants         struct {
        Object string `json:"object"`
        Data   []struct {
            Object      string  `json:"object"`
            ID          string  `json:"id"`
            GrantAmount float64 `json:"grant_amount"`
            UsedAmount  float64 `json:"used_amount"`
            EffectiveAt float64 `json:"effective_at"`
            ExpiresAt   float64 `json:"expires_at"`
        } `json:"data"`
    } `json:"grants"`
}

func CreditGrants(apiKey string) (jd *CreditGrantsResponseBody, err error) {
    var body []byte
    body, err = HttpGet("https://api.openai.com/dashboard/billing/credit_grants", apiKey)
    if err != nil {
        return nil, err
    }

    seelog.Debugf("response: %v", string(body))

    jd = &CreditGrantsResponseBody{}
    err = json.Unmarshal(body, jd)
    if err != nil {
        return nil, err
    }

    return jd, nil
}

package api_base

import (
	"encoding/json"
	"fmt"
	"time"
)

type UsageResponseBody struct {
	Object string `json:"object"`
	Data   []struct {
		AggregationTimestamp  int    `json:"aggregation_timestamp"`
		NRequests             int    `json:"n_requests"`
		Operation             string `json:"operation"`
		SnapshotID            string `json:"snapshot_id"`
		NContext              int    `json:"n_context"`
		NContextTokensTotal   int    `json:"n_context_tokens_total"`
		NGenerated            int    `json:"n_generated"`
		NGeneratedTokensTotal int    `json:"n_generated_tokens_total"`
	} `json:"data"`
	FtData          []interface{} `json:"ft_data"`
	DalleAPIData    []interface{} `json:"dalle_api_data"`
	CurrentUsageUsd float64       `json:"current_usage_usd"`
}

func Usage(apiKey string) (jd *UsageResponseBody, err error) {
	var body []byte

	dt := time.Now().Format("2006-01-02")
	body, err = HttpGet("https://api.openai.com/v1/usage?date="+dt, apiKey)
	if err != nil {
		return nil, err
	}

	//seelog.Debugf("response: %v", string(body))

	jd = &UsageResponseBody{}
	err = json.Unmarshal(body, jd)
	if err != nil {
		return nil, err
	}
	return jd, nil
}

type BillingUsageResponseBody struct {
	Object     string `json:"object"`
	DailyCosts []struct {
		Timestamp float64 `json:"timestamp"`
		LineItems []struct {
			Name string  `json:"name"`
			Cost float64 `json:"cost"`
		} `json:"line_items"`
	} `json:"daily_costs"`
	TotalUsage float64 `json:"total_usage"`
}

func BillingUsage(apiKey string) (jd *BillingUsageResponseBody, err error) {
	var body []byte

	dtEnd := time.Now().Add(time.Hour * 24)
	dtStart := dtEnd.Add(time.Hour * -24 * 15)

	strEnd := dtEnd.Format("2006-01-02")
	strStart := dtStart.Format("2006-01-02")
	link := fmt.Sprintf("https://api.openai.com/dashboard/billing/usage?end_date=%s&start_date=%s", strEnd, strStart)
	body, err = HttpGet(link, apiKey)
	if err != nil {
		return nil, err
	}

	//seelog.Debugf("GET %s, response: %v", link, string(body))

	jd = &BillingUsageResponseBody{}
	err = json.Unmarshal(body, jd)
	if err != nil {
		return nil, err
	}

	return jd, nil
}

package api_base

import (
    "fmt"
)

type SubscriptionResponseBody struct {
    Object             string      `json:"object"`
    HasPaymentMethod   bool        `json:"has_payment_method"`
    Canceled           bool        `json:"canceled"`
    CanceledAt         interface{} `json:"canceled_at"`
    Delinquent         bool        `json:"delinquent"`
    AccessUntil        int         `json:"access_until"`
    SoftLimit          int         `json:"soft_limit"`
    HardLimit          int         `json:"hard_limit"`
    SystemHardLimit    int         `json:"system_hard_limit"`
    SoftLimitUsd       float64     `json:"soft_limit_usd"`
    HardLimitUsd       float64     `json:"hard_limit_usd"`
    SystemHardLimitUsd float64     `json:"system_hard_limit_usd"`
    Plan               struct {
        Title string `json:"title"`
        ID    string `json:"id"`
    } `json:"plan"`
    AccountName    string      `json:"account_name"`
    PoNumber       interface{} `json:"po_number"`
    BillingEmail   interface{} `json:"billing_email"`
    TaxIds         interface{} `json:"tax_ids"`
    BillingAddress struct {
        City       string      `json:"city"`
        Line1      string      `json:"line1"`
        Line2      interface{} `json:"line2"`
        State      string      `json:"state"`
        Country    string      `json:"country"`
        PostalCode string      `json:"postal_code"`
    } `json:"billing_address"`
    BusinessAddress interface{} `json:"business_address"`
}

func Subscription(apiKey string) (body string, err error) {
    var bs []byte

    link := fmt.Sprintf("https://api.openai.com/dashboard/billing/subscription")
    bs, err = HttpGet(link, apiKey)
    return string(bs), err
}

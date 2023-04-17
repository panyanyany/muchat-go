package models

import (
    "gorm.io/gorm"
    "time"
)

const (
    OpenAiAccountStatusActive   = iota // 正常使用
    OpenAiAccountStatusPaused          // (额度耗尽)暂停
    OpenAiAccountStatusDisabled        // 被官方封禁
    OpenAiAccountStatusUnknown  = 99   // 未知
)

type OpenAiAccount struct {
    Email           string `gorm:"uniqueIndex;not null;type:varchar(64)"`
    FirstTime       time.Time
    QueryCnt        int     `gorm:"not null;default:0"`
    UsdSpent        float64 `gorm:"not null;default:0"`
    UsdSpentLimit   float64 `gorm:"not null;default:0"`
    CreditUsed      float64 `gorm:"not null;default:0"`
    CreditAvailable float64 `gorm:"not null;default:0"`
    ExpiresAt       time.Time
    Status          int8   `gorm:"not null;default:0"`
    ApiKey          string `gorm:"uniqueIndex;not null;type:varchar(64)"`
    Name            string
    Password        string
    EmailPassword   string
    gorm.Model
}

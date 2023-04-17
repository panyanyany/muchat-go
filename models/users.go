package models

import (
    "gorm.io/gorm"
    "time"
)

type MuchatUser struct {
    Slug      string `gorm:"uniqueIndex;not null;type:varchar(64)"`
    FirstTime time.Time
    FirstIp   string `gorm:"type:varchar(128)"` // 第一次访问时的 IP
    Usage     int    `gorm:"not null;default:0"`
    MaxUsage  int    `gorm:"not null;default:0"`
    Referer   string `gorm:"type:varchar(128)"`
    BadCnt    int    `gorm:"not null;default:0"`
    MaxDays   int    `gorm:"not null;default:0"`
    ExpiresAt time.Time
    Name      string
    gorm.Model
}

func CheckUser(db *gorm.DB, slug string) (cur CheckUserResult, existed bool) {
    var mdUser MuchatUser
    res := db.Where(MuchatUser{Slug: slug}).First(&mdUser)
    if res.Error == gorm.ErrRecordNotFound {
        return
    }

    existed = true
    cur.User = &mdUser
    cur.Auth = true

    if mdUser.FirstTime.Year() == 1 {
        return
    }
    if mdUser.ExpiresAt.Year() > 1 && mdUser.ExpiresAt.Before(time.Now()) {
        cur.Expired = true
        return
    }
    if mdUser.Usage >= mdUser.MaxUsage {
        cur.ReachCap = true
        return
    }
    return
}

type CheckUserResult struct {
    Auth     bool
    ReachCap bool
    Expired  bool
    User     *MuchatUser
}

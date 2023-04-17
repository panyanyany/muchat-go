package models

import (
    "fmt"
    "go_another_chatgpt/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func InitDb(cfg config.DbConfig) *gorm.DB {
    host := "127.0.0.1"
    if cfg.Host != "" {
        host = cfg.Host
    }
    port := "3306"
    if cfg.Port != "" {
        port = cfg.Port
    }
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User,
        cfg.Pass,
        host,
        port,
        cfg.Name,
    )
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    err = db.AutoMigrate(&MuchatUser{}, &OpenAiAccount{}, &Version{})
    if err != nil {
        panic(err)
    }

    return db
}

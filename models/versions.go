package models

import "gorm.io/gorm"

type Version struct {
	SerialNumber  int `gorm:"uniqueIndex;auto increment"`
	VersionNumber string
	UpdateLog     string
	DownloadLink  string
	gorm.Model
}

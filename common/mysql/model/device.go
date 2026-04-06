package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"type:varchar(50)" json:"username"`
	DeviceIp      string         `gorm:"type:varchar(50)" json:"device_ip"`
	DeviceBrowser string         `gorm:"type:varchar(50)" json:"device_browser"`
	DeviceInfo    string         `gorm:"type:varchar(255)" json:"device_info"`
	AccessToken   string         `gorm:"type:varchar(1024)" json:"access_token"`
	RefreshToken  string         `gorm:"type:varchar(1024)" json:"refresh_token"`
	IssuedTime    time.Time      `gorm:"type:datetime;" json:"issued_time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type DeviceInfo struct {
	Username      string    `json:"username"`
	DeviceIp      string    `json:"deviceIp"`
	DeviceBrowser string    `json:"deviceBrowser"`
	LoginTime     time.Time `json:"loginTime"`
}

package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Password   string
	Phone      string
	Utime      time.Time
	CTime      time.Time
	WechatInfo WechatInfo
}

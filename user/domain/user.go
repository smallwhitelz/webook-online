package domain

import (
	"time"
	domain2 "webook/oauth2/domain"
)

type User struct {
	Id          int64
	Email       string
	Password    string
	Nickname    string
	Birthday    time.Time
	Description string

	Phone string

	// UTC 0 的时区
	Ctime time.Time

	WechatInfo domain2.WechatInfo

	//Addr Address
}

//type Address struct {
//	Province string
//	Region   string
//}

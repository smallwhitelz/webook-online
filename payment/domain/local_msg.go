package domain

import "time"

type LocalMsg struct {
	Id      int64
	Content string
	Ctime   time.Time
	Utime   time.Time
}

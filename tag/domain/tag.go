package domain

import "time"

type Tag struct {
	Id    int64
	Name  string
	Uid   int64
	Ctime time.Time
	Utime time.Time
}

package utils

import (
	"time"
)

func ToThaiTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return t.In(loc)
}
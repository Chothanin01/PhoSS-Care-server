package utils

import (
	"time"
)

func ToThaiTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return t.In(loc)
}

func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

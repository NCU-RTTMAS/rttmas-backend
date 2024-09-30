package timeutils

import "time"

func GetDate() string {
	return time.Now().Format("2006-01-02")
}

func GetDatetime() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

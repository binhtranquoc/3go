package common

import "time"

const (
	YYYYMMDDFormat = "2006-01-02"
)

func FormatDateToYYYYMMDD(t time.Time) string {
	return t.Format(YYYYMMDDFormat)
}

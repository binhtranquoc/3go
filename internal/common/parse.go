package common

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseYYYYMMDDToPgDate(s *string) pgtype.Date {
	if s == nil || *s == "" {
		return pgtype.Date{}
	}
	t, err := ParseYYYYMMDDToTime(*s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func ParseYYYYMMDDToTime(s string) (time.Time, error) {
	return time.Parse(YYYYMMDDFormat, s)
}

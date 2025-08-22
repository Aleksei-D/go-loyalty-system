package models

import (
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	date := t.Time.Format("2006-01-02T15:04:05-07:00")
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse("2006-01-02T15:04:05-07:00", s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

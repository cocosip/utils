package time

import (
	"fmt"
	"time"
)

func Combine(date, tm string) (string, error) {
	dt := fmt.Sprintf("%s %s", date, tm)
	if _, err := time.Parse(time.DateTime, dt); err != nil {
		return "", err
	}
	return dt, nil
}

func Separate(dt string) (string, string, error) {
	datetime, err := time.Parse(time.DateTime, dt)
	if err != nil {
		return "", "", err
	}

	return datetime.Format(time.DateOnly), datetime.Format(time.TimeOnly), nil
}

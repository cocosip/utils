package util

import "time"

func Combine(d string, t string) string {
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		return d
	}

	tm, err := time.Parse(time.TimeOnly, t)
	if err != nil {
		return date.Format(time.DateTime)
	}

	date.Add(time.Hour * time.Duration(tm.Hour()))
	date.Add(time.Minute * time.Duration(tm.Minute()))
	date.Add(time.Second * time.Duration(tm.Second()))

	return date.Format(time.DateTime)
}

func Split(dt string) (string, string) {
	dateTime, err := time.Parse(time.DateTime, dt)
	if err != nil {
		dateTime, err = time.Parse(time.DateOnly, dt)
		if err != nil {
			return time.DateOnly, time.TimeOnly
		}
	}

	date := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 15, 04, 05, 0, time.Local)
	tm := time.Date(2006, 01, 02, dateTime.Hour(), dateTime.Minute(), dateTime.Second(), 0, time.Local)
	return date.Format(time.TimeOnly), tm.Format(time.TimeOnly)
}

package internal

import "time"

func UpDonwDate(now time.Time, before, after int) (time.Time, error) {
	if before != 0 {
		return now.AddDate(0, 0, -1*before), nil
	}
	if after != 0 {
		return now.AddDate(0, 0, after), nil
	}
	return now, nil
}

func GetTargetTime(date string, before, after int) (time.Time, error) {
	if date != "" {
		return time.Parse("2006-01-02", date)
	}
	now := time.Now()
	return UpDonwDate(now, before, after)
}

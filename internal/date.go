package internal

import "time"

func TargetTime(date string, before, after int) (time.Time, error) {
	now := time.Now()
	if date != "" {
		now, err := time.Parse("2006-01-02", date)
		if err != nil {
			return now, err
		}
	}
	if before != 0 {
		return now.AddDate(0, 0, -1*before), nil
	}
	if after != 0 {
		return now.AddDate(0, 0, after), nil
	}
	return now, nil
}

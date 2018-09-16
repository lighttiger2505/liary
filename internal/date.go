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

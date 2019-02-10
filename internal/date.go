package internal

import (
	"errors"
	"time"
)

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

func GetWeakDays(date time.Time) []time.Time {
	wd := date.Weekday()
	weekdays := []time.Time{}
	for i := time.Sunday; i <= time.Saturday; i++ {
		diff := i - wd
		tmpDate := date.AddDate(0, 0, int(diff))
		weekdays = append(weekdays, tmpDate)
	}
	return weekdays
}

func GetDateRange(start time.Time, end time.Time) []time.Time {
	dayDiff := end.Sub(start) / (time.Hour * 24)

	dateRange := []time.Time{start}
	for i := 1; i < int(dayDiff)+1; i++ {
		tmpDate := start.AddDate(0, 0, i)
		dateRange = append(dateRange, tmpDate)
	}

	return dateRange
}

var unitMap = map[string]string{
	"D": "Day",
	"M": "Month",
	"Y": "Year",
	"d": "Day",
	"m": "Month",
	"y": "Year",
}

func ParseDate(s string) (int, int, int, error) {
	// ([0-9]+[a-z]+)+
	orig := s
	var timeDeltaMap = map[string]int{}

	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, 0, 0, nil
	}
	if s == "" {
		return 0, 0, 0, errors.New("invalid date " + orig)
	}
	for s != "" {
		var v int64
		var err error

		// The next character must be [0-9]
		if !('0' <= s[0] && s[0] <= '9') {
			return 0, 0, 0, errors.New("invalid date " + orig)
		}

		// Consume [0-9]*
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, 0, 0, errors.New("invalid duration " + orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, 0, 0, errors.New("missing unit in date " + orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, 0, 0, errors.New("unknown unit " + u + " in duration " + orig)
		}
		timeDeltaMap[unit] = int(v)
	}

	return timeDeltaMap["Year"], timeDeltaMap["Month"], timeDeltaMap["Day"], nil
}

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

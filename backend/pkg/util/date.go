package util

import (
	"fmt"
	"time"
)

var shamsiMonthDays = []int{0, 31, 31, 31, 31, 31, 31, 30, 30, 30, 30, 30, 29}

func IsShamsiLeap(year int) bool {
	remainder := year % 2820
	switch {
	case remainder < 38:
		return (remainder-1)%4 == 0 && remainder != 0
	case remainder < 413:
		return (remainder-38)%4 == 0 && (remainder-38) != 0
	case remainder < 1993:
		return (remainder-413)%4 == 0 && (remainder-413) != 0
	default:
		return (remainder-1993)%4 == 0 && (remainder-1993) != 0
	}
}

func GregorianToShamsi(t time.Time) string {
	year, month, day := t.Date()

	daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if IsGregorianLeap(year) {
		daysInMonth[2] = 29
	}

	dayOfYear := 0
	for m := 1; m < int(month); m++ {
		dayOfYear += daysInMonth[m]
	}
	dayOfYear += day

	var shamsiYear, shamsiMonth, shamsiDay int

	if dayOfYear <= 79 {
		shamsiYear = year - 622
		shamsiMonth = 10
		shamsiDay = dayOfYear + 1
	} else if IsGregorianLeap(year-1) && dayOfYear <= 80 {
		shamsiYear = year - 622
		shamsiMonth = 10
		shamsiDay = dayOfYear
	} else {
		shamsiYear = year - 621
		dayOfYear -= 79
		if dayOfYear <= 186 {
			shamsiMonth = (dayOfYear-1)/31 + 1
			shamsiDay = (dayOfYear-1)%31 + 1
		} else {
			dayOfYear -= 186
			shamsiMonth = (dayOfYear-1)/30 + 7
			shamsiDay = (dayOfYear-1)%30 + 1
		}
	}

	return fmt.Sprintf("%04d-%02d-%02d", shamsiYear, shamsiMonth, shamsiDay)
}

func ShamsiToGregorian(dateStr string) (time.Time, error) {
	var year, month, day int
	_, err := fmt.Sscanf(dateStr, "%d-%d-%d", &year, &month, &day)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid Shamsi date format, expected YYYY-MM-DD: %w", err)
	}

	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("invalid Shamsi month: %d", month)
	}

	maxDays := shamsiMonthDays[month]
	if month == 12 && IsShamsiLeap(year) {
		maxDays = 30
	}

	if day < 1 || day > maxDays {
		return time.Time{}, fmt.Errorf("invalid Shamsi day: %d for month %d", day, month)
	}

	dayOfYear := 0
	for m := 1; m < month; m++ {
		dayOfYear += shamsiMonthDays[m]
		if m == 12 && IsShamsiLeap(year) {
			dayOfYear++
		}
	}
	dayOfYear += day

	var gregYear, gregMonth, gregDay int

	if dayOfYear <= 286 {
		gregYear = year + 621
		dayOfYear += 79

		daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		if IsGregorianLeap(gregYear) {
			daysInMonth[2] = 29
		}

		if dayOfYear <= 31 {
			gregMonth = 1
			gregDay = dayOfYear
		} else {
			dayOfYear -= 31
			for m := 2; m <= 12; m++ {
				if dayOfYear <= daysInMonth[m] {
					gregMonth = m
					gregDay = dayOfYear
					break
				}
				dayOfYear -= daysInMonth[m]
			}
		}
	} else {
		gregYear = year + 622
		dayOfYear -= 286

		daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		if IsGregorianLeap(gregYear) {
			daysInMonth[2] = 29
		}

		if dayOfYear <= 31 {
			gregMonth = 1
			gregDay = dayOfYear
		} else {
			dayOfYear -= 31
			for m := 2; m <= 12; m++ {
				if dayOfYear <= daysInMonth[m] {
					gregMonth = m
					gregDay = dayOfYear
					break
				}
				dayOfYear -= daysInMonth[m]
			}
		}
	}

	return time.Date(gregYear, time.Month(gregMonth), gregDay, 0, 0, 0, 0, time.UTC), nil
}

func FormatShamsi(t time.Time) string {
	return GregorianToShamsi(t)
}

func ParseShamsi(dateStr string) (time.Time, error) {
	return ShamsiToGregorian(dateStr)
}

func IsGregorianLeap(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func FormatShamsiWithTime(t time.Time) string {
	date := GregorianToShamsi(t)
	return fmt.Sprintf("%s %02d:%02d:%02d", date, t.Hour(), t.Minute(), t.Second())
}

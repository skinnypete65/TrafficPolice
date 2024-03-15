package domain

import "time"

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

func (d *Date) Before(another Date) bool {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).
		Before(time.Date(another.Year, another.Month, another.Day, 0, 0, 0, 0, time.UTC))
}

package services

import (
	"fmt"
	"time"
)

func ConvertToDateString(year int, month int, day int) string {
	t_day := time.Date(year, time.Month(month), day, 1, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%d.%s.%d", t_day.Day(), t_day.Month(), t_day.Year())
}

func GetMonthGerman(m time.Month) string {
	switch m {
	case time.January:
		return "Januar"
	case time.February:
		return "Februar"
	case time.March:
		return "März"
	case time.April:
		return "April"
	case time.May:
		return "Mai"
	case time.June:
		return "Juni"
	case time.July:
		return "Juli"
	case time.August:
		return "August"
	case time.September:
		return "September"
	case time.October:
		return "Oktober"
	case time.November:
		return "November"
	case time.December:
		return "Dezember"
	default:
		return "undefined"
	}
}

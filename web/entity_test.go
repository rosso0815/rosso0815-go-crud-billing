package web

import (
	"fmt"
	"testing"
	"time"
)

// Iterate though the months of a year with loop
func ListAllMonths() []time.Month {
	var months []time.Month
	for i := time.January; i <= time.December; i++ {
		fmt.Println("i:", i, int(i))
		months = append(months, i)
	}
	return months
}

func ListYearWithMonths() []time.Month {
	var months []time.Month
	months = append(months, time.January)
	months = append(months, time.February)
	months = append(months, time.March)
	return months
}

// func ListDaysInMonth(t time.Time) []int {
//     days := make([]int, DaysInMonth(t))
//     fmt.Println("days:", days)
//     for i := range days {
//         days[i] = i + 1
//     }
//     return days
// }

// func ListDaysInMonth2(t time.Time) []Day {
// 	days := make([]Day, DaysInMonth(t))
// 	fmt.Println("days:", days)
// 	for i := range days {
// 		days[i].Name = fmt.Sprintf("%d", i)
// 		// days[i].Workday = t.Day()
// 	}
// 	return days
// }

func Test_ShowMonths(t *testing.T) {
	// t.Log("ShowMonths:", ListYearWithMonths())
	t.Log("ShowMonths:", ListAllMonths())

}
func Test_ShowMonthdays(t *testing.T) {
	t1 := time.Date(2025, time.May+1, 0, 0, 0, 0, 0, time.Now().Local().Location())
	t.Log("Year(): ", t1.Year())
	t.Log("Month(): ", t1.Month())
	t.Log("Day(): ", t1.Day())
	t.Log("Hour(): ", t1.Hour())
	t.Log("Minute(): ", t1.Minute())
	t.Log("Weekday(): ", t1.Weekday())

	// t.Log("Day(): ", time.Date(2025, 14, 0, 0, 0, 0, 0, time.UTC).Day())

	// year := 2025
	// month := time.May
	// m := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// days := ListDaysInMonth2(m)

	// t.Log(m, days)
}

// func Test_ShowDays(t *testing.T) {
//     t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//     days := ListDaysInMonth(t1)
//     if len(days) != 31 {
//         t.Errorf("Expected 31 days in January, got %d", len(days))
//     }
//
//     t2 := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
//     days = ListDaysInMonth(t2)
//     if len(days) != 28 {
//         t.Errorf("Expected 28 days in February (non-leap year), got %d", len(days))
//     }
//
//     t3 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // Leap year
//     days = ListDaysInMonth(t3)
//     if len(days) != 29 {
//         t.Errorf("Expected 29 days in February (leap year), got %d", len(days))
//     }
//     fmt.Println("t3:", t3.Weekday())
//     fmt.Println("days:", days)
//     fmt.Println("now:", time.Now().Year(), time.Now().Month(), time.Now().Day())
//
// }

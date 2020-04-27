package utils

import (
	"fmt"
	"time"
)

// Min function returns the minimum uint64 value from x and y values
func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// Max function returns the maximum uint64 value from x and y values
func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

// MaxUint function returns the maximum uint value from x and y values
func MaxUint(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// MinUint function returns the minimum uint value from x and y values
func MinUint(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// NumberOfWeeksSinceEpoch returns the number of weeks since epoch date (01/01/1970) for time.Now()
func NumberOfWeeksSinceEpoch() uint {
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	today := time.Now()

	for epoch.Weekday() != time.Monday { // iterate back to Monday
		epoch = epoch.AddDate(0, 0, -1)
	}

	for today.Weekday() != time.Monday { // iterate back to Monday
		today = today.AddDate(0, 0, -1)
	}

	diff := today.Sub(epoch)

	return uint(diff.Hours() / 24 / 7)
}

// NumberOfWeeksSinceEpochFromDate returns the number of weeks since epoch date (01/01/1970) for date
func NumberOfWeeksSinceEpochFromDate(date time.Time) uint {
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	for epoch.Weekday() != time.Monday { // iterate back to Monday
		epoch = epoch.AddDate(0, 0, -1)
	}

	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
	}

	diff := date.Sub(epoch)

	return uint(diff.Hours() / 24 / 7)
}

// WeekDate represents a date with a week number from epoch and a number of seconds since start of week (monday midnight)
type WeekDate struct {
	WeekNumber uint
	Seconds    float64
}

// WeekTimeRange is a period of time between two dates represented as WeekDate
type WeekTimeRange struct {
	MinWeekDate WeekDate
	MaxWeekDate WeekDate
}

// GetWeekDate calculates a WeekDate from a specific date
func GetWeekDate(date time.Time) WeekDate {
	monday := date
	for monday.Weekday() != time.Monday { // iterate back to Monday
		monday = monday.AddDate(0, 0, -1)
	}

	mondayMidnight := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)

	weekNum := NumberOfWeeksSinceEpochFromDate(mondayMidnight)

	duration := date.Sub(mondayMidnight)

	return WeekDate{
		WeekNumber: weekNum,
		Seconds:    duration.Seconds(),
	}
}

// GetWeekTimeRangeFromDateRange cuts in weeks from dateMin to dateMax and returns an array of WeekTimeRange
func GetWeekTimeRangeFromDateRange(dateMin, dateMax time.Time) ([]*WeekTimeRange, error) {
	if dateMin.After(dateMax) {
		return nil, fmt.Errorf("dateMin superiror to dateMax")
	}

	result := make([]*WeekTimeRange, 0)

	weekDateMin := GetWeekDate(dateMin)
	weekDateMax := GetWeekDate(dateMax)

	numberOfWeeks := weekDateMax.WeekNumber - weekDateMin.WeekNumber
	if numberOfWeeks == 0 {
		result = append(result, &WeekTimeRange{MinWeekDate: weekDateMin, MaxWeekDate: weekDateMax})
		return result, nil
	}

	currentWeekTimeRange := WeekTimeRange{MinWeekDate: weekDateMin}
	currentWeekNumber := weekDateMin.WeekNumber
	for currentWeekNumber <= weekDateMax.WeekNumber {
		if currentWeekNumber == weekDateMax.WeekNumber {
			currentWeekTimeRange.MaxWeekDate = weekDateMax
		} else {
			currentWeekTimeRange.MaxWeekDate = WeekDate{WeekNumber: currentWeekNumber, Seconds: 604800}
		}

		result = append(result, &WeekTimeRange{
			MinWeekDate: WeekDate{WeekNumber: currentWeekTimeRange.MinWeekDate.WeekNumber, Seconds: currentWeekTimeRange.MinWeekDate.Seconds},
			MaxWeekDate: WeekDate{WeekNumber: currentWeekTimeRange.MaxWeekDate.WeekNumber, Seconds: currentWeekTimeRange.MaxWeekDate.Seconds},
		})

		currentWeekNumber++

		currentWeekTimeRange = WeekTimeRange{MinWeekDate: WeekDate{WeekNumber: currentWeekNumber, Seconds: 0}}
	}

	return result, nil
}

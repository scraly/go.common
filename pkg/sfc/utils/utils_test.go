package utils_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/scraly/go.common/pkg/sfc/utils"
)

func TestGetWeekDate(t *testing.T) {
	tests := []struct {
		date       time.Time
		WeekNumber uint
		Seconds    float64
	}{
		{
			date:       time.Date(2018, 10, 8, 12, 0, 0, 0, time.UTC),
			WeekNumber: 2545,
			Seconds:    43200,
		},
		{
			date:       time.Date(2018, 10, 10, 15, 42, 3, 30, time.UTC),
			WeekNumber: 2545,
			Seconds:    229323.00000003,
		},
		{
			date:       time.Date(2018, 10, 17, 15, 42, 3, 30, time.UTC),
			WeekNumber: 2546,
			Seconds:    229323.00000003,
		},
	}

	for _, tt := range tests {
		t.Run("GetWeekDate", func(t *testing.T) {
			RegisterTestingT(t)

			result := utils.GetWeekDate(tt.date)

			Expect(result.WeekNumber).To(Equal(tt.WeekNumber))
			Expect(result.Seconds).To(Equal(tt.Seconds))
		})
	}
}

func TestGetWeekTimeRangeFromDateRange(t *testing.T) {
	tests := []struct {
		dateMin        time.Time
		dateMax        time.Time
		WeekTimeRanges []utils.WeekTimeRange
		count          int
		Error          bool
	}{
		{
			dateMin: time.Date(2018, 10, 8, 12, 0, 0, 0, time.UTC),
			dateMax: time.Date(2018, 10, 10, 15, 42, 3, 30, time.UTC),
			WeekTimeRanges: []utils.WeekTimeRange{{
				MinWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 43200},
				MaxWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 229323.00000003},
			}},
			count: 1,
			Error: false,
		},
		{
			dateMin: time.Date(2018, 10, 14, 16, 30, 0, 0, time.UTC),
			dateMax: time.Date(2018, 10, 15, 16, 30, 0, 30, time.UTC),
			WeekTimeRanges: []utils.WeekTimeRange{
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 577800},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 604800},
				},
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 0},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 59400.00000003},
				},
			},
			count: 2,
			Error: false,
		},
		{
			dateMin: time.Date(2018, 10, 8, 12, 0, 0, 0, time.UTC),
			dateMax: time.Date(2018, 10, 17, 15, 42, 3, 30, time.UTC),
			WeekTimeRanges: []utils.WeekTimeRange{
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 43200},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 604800},
				},
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 0},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 229323.00000003},
				},
			},
			count: 2,
			Error: false,
		},
		{
			dateMin: time.Date(2018, 10, 8, 12, 0, 0, 0, time.UTC),
			dateMax: time.Date(2018, 10, 24, 15, 42, 3, 30, time.UTC),
			WeekTimeRanges: []utils.WeekTimeRange{
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 43200},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2545, Seconds: 604800},
				},
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 0},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2546, Seconds: 604800},
				},
				{
					MinWeekDate: utils.WeekDate{WeekNumber: 2547, Seconds: 0},
					MaxWeekDate: utils.WeekDate{WeekNumber: 2547, Seconds: 229323.00000003},
				},
			},
			count: 3,
			Error: false,
		},
	}

	for _, tt := range tests {
		t.Run("GetWeekDate", func(t *testing.T) {
			RegisterTestingT(t)

			result, err := utils.GetWeekTimeRangeFromDateRange(tt.dateMin, tt.dateMax)
			if tt.Error {
				Expect(err).To(Not(BeNil()))
			} else {
				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(tt.count))

				for i := 0; i < tt.count; i++ {
					Expect(result[i].MinWeekDate.WeekNumber).To(Equal(tt.WeekTimeRanges[i].MinWeekDate.WeekNumber))
					Expect(result[i].MinWeekDate.Seconds).To(Equal(tt.WeekTimeRanges[i].MinWeekDate.Seconds))
					Expect(result[i].MaxWeekDate.WeekNumber).To(Equal(tt.WeekTimeRanges[i].MaxWeekDate.WeekNumber))
					Expect(result[i].MaxWeekDate.Seconds).To(Equal(tt.WeekTimeRanges[i].MaxWeekDate.Seconds))
				}
			}

		})
	}
}

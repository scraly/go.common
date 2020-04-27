package z3

import (
	"math"
	"time"

	"github.com/scraly/go.common/pkg/sfc"
	"github.com/scraly/go.common/pkg/sfc/utils"
	"github.com/scraly/go.common/pkg/sfc/zorder"
	zranges "github.com/scraly/go.common/pkg/sfc/zorder/zrange"
)

// Z3Search struct
type Z3Search struct { // nolint: golint
	curve zorder.SpaceTimeFillingCurve
}

// NewSearch constructs a Z3Search and returns it as sfc.Z3Search
func NewSearch() (sfc.Z3Search, error) {
	z3Sfc, errZ3Sfc := NewZ3Sfc(NormalizerMaxPrecision)
	if errZ3Sfc != nil {
		return nil, errZ3Sfc
	}

	return &Z3Search{
		curve: z3Sfc,
	}, nil
}

// GetZ3Ranges returns all Z3 Index Ranges corresponding to bounding box bbox and time frame (dateMin, dateMax)
func (z3search *Z3Search) GetZ3Ranges(bbox sfc.BoundingBox, dateMin, dateMax time.Time) ([]*sfc.IndexRange, []*utils.WeekTimeRange, error) {
	z3Sfc := z3search.curve

	weekTimeRanges, errWeekTimeRanges := utils.GetWeekTimeRangeFromDateRange(dateMin, dateMax)
	if errWeekTimeRanges != nil {
		return nil, nil, errWeekTimeRanges
	}

	zbounds := make([]zorder.ZRange, 0)
	for _, weekTimeRange := range weekTimeRanges {
		z3nMin, errZ3nMin := z3Sfc.Index(bbox.SouthWest.Longitude, bbox.SouthWest.Latitude, uint64(math.Floor(weekTimeRange.MinWeekDate.Seconds)))
		if errZ3nMin != nil {
			return nil, nil, errZ3nMin
		}

		z3nMax, errZ3nMax := z3Sfc.Index(bbox.NorthEast.Longitude, bbox.NorthEast.Latitude, uint64(math.Floor(weekTimeRange.MaxWeekDate.Seconds)))
		if errZ3nMax != nil {
			return nil, nil, errZ3nMax
		}

		zrange, errZRange := zorder.NewZRange(z3nMin.GetZValue(), z3nMax.GetZValue())
		if errZRange != nil {
			return nil, nil, errZRange
		}
		zbounds = append(zbounds, *zrange)
	}

	result, errRanges := zranges.CalculateRanges(NewZ3(), zbounds, 64, 0, 7)
	if errRanges != nil {
		return nil, nil, errRanges
	}

	return result, weekTimeRanges, nil

}

// GetSpaceTimeFillingCurve returns the z3 spce filling curve
func (z3search *Z3Search) GetSpaceTimeFillingCurve() zorder.SpaceTimeFillingCurve {
	return z3search.curve
}

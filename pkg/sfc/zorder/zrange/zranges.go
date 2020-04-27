package zrange

import (
	"fmt"
	"math"
	"sort"

	"github.com/scraly/go.common/pkg/sfc"
	"github.com/scraly/go.common/pkg/sfc/utils"
	api "github.com/scraly/go.common/pkg/sfc/zorder"
)

type workItem struct {
	min  uint64
	max  uint64
	stop bool
}

// CalculateRanges method returns
func CalculateRanges(zn api.ZN, zbounds []api.ZRange, precision, maxRanges, maxRecurse int) ([]*sfc.IndexRange, error) {
	ranges := make([]*sfc.IndexRange, 0)

	zBoundsValues := make([]uint64, len(zbounds)*2)
	index := 0
	for _, zBoundValue := range zbounds {
		zBoundsValues[index] = zBoundValue.Min
		zBoundsValues[index+1] = zBoundValue.Max
		index = index + 2
	}
	commonPrefix, commonBits, errPrefix := LongestCommonPrefix(zn.GetTotalBits(), zn.GetDimensions(), zBoundsValues)
	if errPrefix != nil {
		return nil, errPrefix
	}

	offset := uint(64 - commonBits)

	remaining := utils.NewRingqueue()

	iRange := CheckValue(zn, zbounds, commonPrefix, 0, offset, precision, remaining)
	if iRange != nil {
		ranges = append(ranges, iRange)
	}
	remaining.Add(workItem{
		stop: true,
	})
	offset -= uint(zn.GetDimensions())

	// level of recursion
	level := 0

	if maxRanges == 0 {
		maxRanges = math.MaxInt32
	}

	if maxRecurse == 0 {
		maxRecurse = 20
	}

	for level < maxRecurse && offset >= 0 && remaining.Len() > 0 && len(ranges) < maxRanges {
		next, _ := remaining.Remove()

		if next.(workItem).stop {
			// we've fully processed a level, increment our state
			if remaining.Len() > 0 {
				level++
				offset -= uint(zn.GetDimensions())
				remaining.Add(workItem{
					stop: true,
				})
			}
		} else {
			prefix := next.(workItem).min
			quadrant := 0
			for quadrant < zn.GetQuadrants() {
				indexRange := CheckValue(zn, zbounds, prefix, uint64(quadrant), offset, precision, remaining)
				if indexRange != nil {
					ranges = append(ranges, indexRange)
				}
				quadrant++
			}
		}
	}

	// bottom out and get all the ranges that partially overlapped but we didn't fully process
	for remaining.Len() > 0 {
		minmax, _ := remaining.Remove()
		if !minmax.(workItem).stop {
			indexRange := &sfc.IndexRange{
				Lower:     minmax.(workItem).min,
				Upper:     minmax.(workItem).max,
				Contained: false,
			}
			ranges = append(ranges, indexRange)
		}
	}

	// we've got all our ranges - now reduce them down by merging overlapping values
	sort.SliceStable(ranges, func(i, j int) bool {

		if ranges[i].Lower == ranges[j].Lower {
			return ranges[i].Upper < ranges[j].Upper
		}

		return ranges[i].Lower < ranges[j].Lower

	})

	currentRange := ranges[0]
	result := make([]*sfc.IndexRange, 0)

	i := 1
	for i <= len(ranges)-1 {
		rangeZ := ranges[i]
		if rangeZ.Lower <= currentRange.Upper+1 {
			// merge the two ranges
			currentRange = &sfc.IndexRange{
				Lower:     currentRange.Lower,
				Upper:     utils.Max(currentRange.Upper, rangeZ.Upper),
				Contained: currentRange.Contained && rangeZ.Contained,
			}
		} else {
			// append the last range and set the current range for future merging
			result = append(result, currentRange)
			currentRange = rangeZ
		}
		i++
	}

	result = append(result, currentRange)

	return result, nil
}

// LongestCommonPrefix calculates the longes common prefix for all uint64 contained in values
func LongestCommonPrefix(totalBits, dimensions int, values []uint64) (prefix uint64, precision int, err error) {
	if len(values) < 2 {
		return 0, 0, fmt.Errorf("Wrong number of elements %d", len(values))
	}

	bitShift := totalBits - dimensions

	head := values[0] >> uint(bitShift)

	continueSearch := true

	for bitShift > -1 && continueSearch {
		for i, value := range values {
			if i != 0 {
				if (value>>uint(bitShift)) == head && bitShift > -1 {
					continueSearch = true
				} else {
					continueSearch = false
				}
			}
		}
		if continueSearch {
			bitShift -= dimensions
			head = values[0] >> uint(bitShift)
		}
	}

	bitShift += dimensions
	prefix = values[0] & (math.MaxInt64 << uint(bitShift))
	precision = 64 - bitShift

	return prefix, precision, nil
}

// IsContained returns true if rangeZ is fully contained in any ZRange from zbounds
func IsContained(zn api.ZN, zbounds []api.ZRange, rangeZ api.ZRange) bool {
	i := 0
	for i < len(zbounds) {
		if zn.Contains(zbounds[i], rangeZ.Min) && zn.Contains(zbounds[i], rangeZ.Max) {
			return true
		}
		i++
	}

	return false
}

// IsOverlapped returns true if rangeZ is overlapped any ZRange from zbounds
func IsOverlapped(zn api.ZN, zbounds []api.ZRange, rangeZ api.ZRange) bool {
	i := 0
	for i < len(zbounds) {
		if zn.Overlaps(zbounds[i], rangeZ) {
			return true
		}
		i++
	}

	return false
}

// CheckValue checks a single value and either: eliminates it as out of bounds or adds it to our results as fully matching, or queues up it's children for further processing
func CheckValue(zn api.ZN, zbounds []api.ZRange, prefix uint64, quadrant uint64, offset uint, precision int, remaining *utils.Ringqueue) *sfc.IndexRange {
	min := prefix | (quadrant << offset)     // QR + 000...
	max := min | ((uint64(1) << offset) - 1) // QR + 111...

	quadrantRange, errRange := api.NewZRange(min, max)
	if errRange != nil {
		// TODO use logger
		fmt.Println(errRange)
	}

	if IsContained(zn, zbounds, *quadrantRange) || (offset < uint(64-precision)) {
		// whole range matches, happy day
		indexRange := &sfc.IndexRange{
			Lower:     quadrantRange.Min,
			Upper:     quadrantRange.Max,
			Contained: true,
		}
		return indexRange
	} else if IsOverlapped(zn, zbounds, *quadrantRange) {
		// some portion of this range is excluded
		// queue up each sub-range for processing
		remaining.Add(workItem{
			min:  quadrantRange.Min,
			max:  quadrantRange.Max,
			stop: false,
		})
	}

	return nil

}

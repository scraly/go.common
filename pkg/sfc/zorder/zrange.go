package zorder

import (
	"errors"
)

// ZRange struct
type ZRange struct {
	Min, Max uint64
}

// NewZRange constructs a ZRange from two uint64. Represents a rectangle in defined by min and max as two opposing points
func NewZRange(min, max uint64) (*ZRange, error) {
	if min > max {
		return nil, errors.New("Range bounds must be ordered")
	}

	return &ZRange{
		Min: min,
		Max: max,
	}, nil

}

// Mid method returns mid value of bbbox
func (zRange *ZRange) Mid() uint64 {
	return (zRange.Max + zRange.Min) >> 1
}

// Length method returns length of ZRange
func (zRange *ZRange) Length() uint64 {
	return zRange.Max - zRange.Min + 1
}

// Contains indicates if ZRange contains the z value
func (zRange *ZRange) Contains(bits uint64) bool {
	return bits >= zRange.Min && bits <= zRange.Max
}

// ContainsZ indicates if ZRange contains other ZRange
func (zRange *ZRange) ContainsZ(otherZRange *ZRange) bool {
	return zRange.Contains(otherZRange.Min) && zRange.Contains(otherZRange.Max)
}

// Overlaps indicates if ZRange is overlapped by other ZRange
func (zRange *ZRange) Overlaps(otherZRange *ZRange) bool {
	return zRange.Contains(otherZRange.Min) || zRange.Contains(otherZRange.Max)
}

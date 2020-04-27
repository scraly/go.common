package sfc

import (
	"time"

	"github.com/scraly/go.common/pkg/sfc/utils"
	"github.com/scraly/go.common/pkg/sfc/zorder"
)

// Point structure has simply Laongitude and Latitude values
type Point struct {
	Longitude float64
	Latitude  float64
}

// BoundingBox defines a rectangle with two points (SouthWest and NorthEast)
type BoundingBox struct {
	SouthWest Point
	NorthEast Point
}

// IndexRange struct is the result unit for the ZN ranges
type IndexRange struct {
	Lower     uint64
	Upper     uint64
	Contained bool
}

// Z3Search interface
type Z3Search interface {
	GetSpaceTimeFillingCurve() zorder.SpaceTimeFillingCurve
	GetZ3Ranges(bbox BoundingBox, dateMin, dateMax time.Time) ([]*IndexRange, []*utils.WeekTimeRange, error)
}

// Z2Search interface
type Z2Search interface {
	GetSpaceFillingCurve() zorder.SpaceFillingCurve
	GetZ2Ranges(bbox BoundingBox) ([]*IndexRange, error)
}

package z2

import (
	"fmt"

	"github.com/scraly/go.common/pkg/sfc"
	"github.com/scraly/go.common/pkg/sfc/zorder"
	zranges "github.com/scraly/go.common/pkg/sfc/zorder/zrange"
)

// Z2Search struct
type Z2Search struct { // nolint: golint
	curve zorder.SpaceFillingCurve
}

// NewSearch constructs a Z2Search and returns it as sfc.Z2Search
func NewSearch(precision uint) (sfc.Z2Search, error) {
	if precision < NormalizerMinPrecision || precision > NormalizerMaxPrecision {
		return nil, fmt.Errorf("precision must be in range [%d, %d]", NormalizerMinPrecision, NormalizerMaxPrecision)
	}

	z2Sfc, errZ2Sfc := NewZ2Sfc(precision)
	if errZ2Sfc != nil {
		return nil, errZ2Sfc
	}

	return &Z2Search{
		curve: z2Sfc,
	}, nil
}

// GetZ2Ranges returns all Z2 Index Ranges corresponding to bounding box bbox
func (z2search *Z2Search) GetZ2Ranges(bbox sfc.BoundingBox) ([]*sfc.IndexRange, error) {
	z2Sfc := z2search.curve

	zbounds := make([]zorder.ZRange, 0)

	z2nMin, errZ2nMin := z2Sfc.Index(bbox.SouthWest.Longitude, bbox.SouthWest.Latitude)
	if errZ2nMin != nil {
		return nil, errZ2nMin
	}

	z2nMax, errZ2nMax := z2Sfc.Index(bbox.NorthEast.Longitude, bbox.NorthEast.Latitude)
	if errZ2nMax != nil {
		return nil, errZ2nMax
	}

	zrange, errZRange := zorder.NewZRange(z2nMin.GetZValue(), z2nMax.GetZValue())
	if errZRange != nil {
		return nil, errZRange
	}
	zbounds = append(zbounds, *zrange)

	result, errRanges := zranges.CalculateRanges(NewZ2(), zbounds, 64, 0, 7)
	if errRanges != nil {
		return nil, errRanges
	}

	return result, nil

}

// GetSpaceFillingCurve returns the z2 spce filling curve
func (z2search *Z2Search) GetSpaceFillingCurve() zorder.SpaceFillingCurve {
	return z2search.curve
}

package z2

import (
	"errors"
	"fmt"

	"github.com/scraly/go.common/pkg/sfc/zorder"
	"github.com/scraly/go.common/pkg/sfc/zorder/normalizer"
)

const (
	// NormalizerMaxPrecision is the maximum allowed precision value
	NormalizerMaxPrecision = 31 // Inclusive
	// NormalizerMinPrecision is the minimum allowed precision value
	NormalizerMinPrecision = 1 // Inclusive
)

// Z2Sfc struct
type Z2Sfc struct { // nolint: golint
	precision     uint
	lonNormalizer zorder.Normalizer
	latNormalizer zorder.Normalizer
}

// NewZ2Sfc creates a Z2 space filling curve
func NewZ2Sfc(precision uint) (zorder.SpaceFillingCurve, error) {

	if precision < NormalizerMinPrecision || precision > NormalizerMaxPrecision {
		return nil, errors.New("Invalid precision")
	}

	lonNormalizer, errLonNormalizer := normalizer.NewNormalizer(-180, 180, precision)
	if errLonNormalizer != nil {
		return nil, errLonNormalizer
	}

	latNormalizer, errLatNormalizer := normalizer.NewNormalizer(-90, 90, precision)
	if errLatNormalizer != nil {
		return nil, errLatNormalizer
	}

	return &Z2Sfc{
		precision:     precision,
		lonNormalizer: lonNormalizer,
		latNormalizer: latNormalizer,
	}, nil
}

// Index function takes longitude x latitude y and outputs a Z2 with normalized values interleaved
func (z2sfc *Z2Sfc) Index(x, y float64) (zorder.Z2N, error) {
	if x < z2sfc.lonNormalizer.GetMin() || x > z2sfc.lonNormalizer.GetMax() {
		return nil, fmt.Errorf("x value %f out of range (min %f, max %f)", x, z2sfc.lonNormalizer.GetMin(), z2sfc.lonNormalizer.GetMax())
	}

	if y < z2sfc.latNormalizer.GetMin() || y > z2sfc.latNormalizer.GetMax() {
		return nil, fmt.Errorf("y value %f out of range (min %f, max %f)", y, z2sfc.latNormalizer.GetMin(), z2sfc.latNormalizer.GetMax())
	}

	z2 := NewZ2()

	nx := z2sfc.lonNormalizer.Normalize(x)
	ny := z2sfc.latNormalizer.Normalize(y)

	z2.Apply(nx, ny)

	return z2, nil
}

// Invert function takes a Z2 and outputs denormalized deinterleaved lonvitude, latitude
func (z2sfc *Z2Sfc) Invert(z2 zorder.Z2N) (float64, float64) {
	x, y := z2.UnApply(z2.GetZValue())
	return z2sfc.lonNormalizer.DeNormalize(x), z2sfc.latNormalizer.DeNormalize(y)
}

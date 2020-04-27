package z3

import (
	"errors"
	"fmt"
	"math"

	"github.com/scraly/go.common/pkg/sfc/zorder"
	"github.com/scraly/go.common/pkg/sfc/zorder/normalizer"
)

const (
	// NormalizerMaxPrecision is the maximum allowed precision value
	NormalizerMaxPrecision = 21 // Inclusive
	// NormalizerMinPrecision is the minimum allowed precision value
	NormalizerMinPrecision = 1 // Inclusive
)

// Z3Sfc struct
type Z3Sfc struct { // nolint: golint
	lonNormalizer  zorder.Normalizer
	latNormalizer  zorder.Normalizer
	timeNormalizer zorder.Normalizer
}

// NewZ3Sfc creates a Z3 space filling curve
func NewZ3Sfc(precision uint) (zorder.SpaceTimeFillingCurve, error) {
	return NewZ3SfcWithPrecision(precision, precision, precision)
}

// NewZ3SfcWithPrecision creates a Z3 space filling curve
func NewZ3SfcWithPrecision(lonPrecision, latPrecision, timePrecision uint) (zorder.SpaceTimeFillingCurve, error) {
	if lonPrecision < NormalizerMinPrecision || lonPrecision > NormalizerMaxPrecision {
		return nil, errors.New("Longitude precision (bits) per dimension must be in range [1,21]")
	}

	if latPrecision < NormalizerMinPrecision || latPrecision > NormalizerMaxPrecision {
		return nil, errors.New("Latitude precision (bits) per dimension must be in range [1,21]")
	}

	if timePrecision < NormalizerMinPrecision || timePrecision > NormalizerMaxPrecision {
		return nil, errors.New("Time precision (bits) per dimension must be in range [1,21]")
	}

	lonNormalizer, errLonNormalizer := normalizer.NewNormalizer(-180, 180, lonPrecision)
	if errLonNormalizer != nil {
		return nil, errLonNormalizer
	}

	latNormalizer, errLatNormalizer := normalizer.NewNormalizer(-90, 90, latPrecision)
	if errLatNormalizer != nil {
		return nil, errLatNormalizer
	}

	// Seconds from beginning of the week
	timeNormalizer, errTimeNormalizer := normalizer.NewNormalizer(0, 604800, timePrecision)
	if errTimeNormalizer != nil {
		return nil, errTimeNormalizer
	}

	return &Z3Sfc{
		lonNormalizer:  lonNormalizer,
		latNormalizer:  latNormalizer,
		timeNormalizer: timeNormalizer,
	}, nil
}

// Index function takes longitude x latitude y time t and outputs a Z3 with normalized values interleaved
func (z3sfc *Z3Sfc) Index(x, y float64, t uint64) (zorder.Z3N, error) {
	if x < z3sfc.lonNormalizer.GetMin() || x > z3sfc.lonNormalizer.GetMax() {
		return nil, fmt.Errorf("x value %f out of range (min %f, max %f)", x, z3sfc.lonNormalizer.GetMin(), z3sfc.lonNormalizer.GetMax())
	}

	if y < z3sfc.latNormalizer.GetMin() || y > z3sfc.latNormalizer.GetMax() {
		return nil, fmt.Errorf("y value %f out of range (min %f, max %f)", y, z3sfc.latNormalizer.GetMin(), z3sfc.latNormalizer.GetMax())
	}

	if t < uint64(z3sfc.timeNormalizer.GetMin()) || t > uint64(z3sfc.timeNormalizer.GetMax()) {
		return nil, fmt.Errorf("t value %d out of range (min %f, max %f)", t, z3sfc.timeNormalizer.GetMin(), z3sfc.timeNormalizer.GetMax())
	}

	z3 := NewZ3()

	nx := z3sfc.lonNormalizer.Normalize(x)
	ny := z3sfc.latNormalizer.Normalize(y)
	nt := z3sfc.timeNormalizer.Normalize(float64(t))

	z3.Apply(nx, ny, nt)

	return z3, nil
}

// Invert function takes a Z3 and outputs denormalized deinterleaved lonvitude, latitude and time values
func (z3sfc *Z3Sfc) Invert(z3 zorder.Z3N) (float64, float64, uint64) {
	x, y, t := z3.UnApply(z3.GetZValue())
	return z3sfc.lonNormalizer.DeNormalize(x), z3sfc.latNormalizer.DeNormalize(y), uint64(math.Round(z3sfc.timeNormalizer.DeNormalize(t)))
}

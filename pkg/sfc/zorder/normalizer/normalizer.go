package normalizer

import (
	"errors"
	"math"

	api "github.com/scraly/go.common/pkg/sfc/zorder"
)

// BitNormalizedDimension struct
type BitNormalizedDimension struct {
	precision    uint
	min          float64
	max          float64
	bins         uint64
	normalizer   float64
	deNormalizer float64
	maxIndex     uint
}

// NewNormalizer method constructs a BitNormalizedDimension struct and returns it as Normalizer interface
func NewNormalizer(min, max float64, precision uint) (api.Normalizer, error) {
	if precision == 0 || precision > 31 {
		return nil, errors.New("Precision (bits) must be in [1,31]")
	}

	bins := 1 << precision

	bitNormalizer := &BitNormalizedDimension{
		precision:    precision,
		min:          min,
		max:          max,
		bins:         uint64(bins),
		normalizer:   float64(bins) / (max - min),
		deNormalizer: (max - min) / float64(bins),
		maxIndex:     uint(bins - 1),
	}

	return bitNormalizer, nil
}

// Normalize method
func (b *BitNormalizedDimension) Normalize(value float64) uint {
	if value >= b.max {
		return b.maxIndex
	}

	return uint(math.Floor((value - b.min) * b.normalizer))
}

// DeNormalize method
func (b *BitNormalizedDimension) DeNormalize(n uint) float64 {
	if n >= b.maxIndex {
		return b.min + (float64(b.maxIndex)+0.5)*b.deNormalizer
	}

	return b.min + (float64(n)+0.5)*b.deNormalizer
}

// GetMin returns min value of the normalizer
func (b *BitNormalizedDimension) GetMin() float64 {
	return b.min
}

// GetMax returns max value of the normalizer
func (b *BitNormalizedDimension) GetMax() float64 {
	return b.max
}

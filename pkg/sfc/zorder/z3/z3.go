package z3

import (
	"math"
	"sync"

	"github.com/scraly/go.common/pkg/sfc/utils"
	api "github.com/scraly/go.common/pkg/sfc/zorder"
)

// Z3 struct
type Z3 struct {
	*sync.RWMutex

	BitsPerDimension int
	Dimensions       int
	MaxMask          uint64
	TotalBits        int
	Quadrants        int
	ZValue           uint64
}

// NewZ3 method constructs a Z3 struct and returns it as Z3N interface
func NewZ3() api.Z3N {
	return newInternalZ3(0)
}

// NewZ3WithZValue constructs a Z3 struct with initial z3 value and returns it as Z3N interface
func NewZ3WithZValue(z uint64) api.Z3N {
	return newInternalZ3(z)
}

func newInternalZ3(z uint64) api.Z3N {
	dimensions := 3
	return &Z3{
		RWMutex:          &sync.RWMutex{},
		BitsPerDimension: 21,
		Dimensions:       dimensions,
		MaxMask:          0x1fffff,
		TotalBits:        63,
		Quadrants:        int(math.Pow(2, float64(dimensions))),
		ZValue:           z,
	}
}

// Apply method calculates a z3 from x,y,z
func (z3 *Z3) Apply(x, y, z uint) uint64 {
	z3.setZValue(z3.split(uint64(x)) | z3.split(uint64(y))<<1 | z3.split(uint64(z))<<2)
	return z3.GetZValue()
}

// UnApply method calculates x,y,z from z3
func (z3 *Z3) UnApply(z uint64) (uint, uint, uint) {
	z3.setZValue(z)
	return z3.combine(z), z3.combine(z >> 1), z3.combine(z >> 2)
}

// GetDimensions returns the number of dimensions (3 for Z3)
func (z3 *Z3) GetDimensions() int {
	return z3.Dimensions
}

// GetQuadrants returns the number of quadrants (8 for Z3)
func (z3 *Z3) GetQuadrants() int {
	return z3.Quadrants
}

// GetTotalBits returns the total bits used to encode Z3
func (z3 *Z3) GetTotalBits() int {
	return z3.TotalBits
}

// GetZValue returns the computed Z3 value
func (z3 *Z3) GetZValue() uint64 {
	z3.RLock()
	defer z3.RUnlock()

	return z3.ZValue
}

func (z3 *Z3) setZValue(zValue uint64) {
	z3.Lock()
	defer z3.Unlock()

	z3.ZValue = zValue
}

// Contains indicates if z3 value is contained in rangeZ
func (z3 *Z3) Contains(rangeZ api.ZRange, value uint64) bool {
	x, y, z := z3.UnApply(value)
	d0Min := z3.dim(0, rangeZ.Min)
	d1Min := z3.dim(1, rangeZ.Min)
	d2Min := z3.dim(2, rangeZ.Min)
	d0Max := z3.dim(0, rangeZ.Max)
	d1Max := z3.dim(1, rangeZ.Max)
	d2Max := z3.dim(2, rangeZ.Max)

	return x >= d0Min && x <= d0Max && y >= d1Min && y <= d1Max && z >= d2Min && z <= d2Max
}

// Overlaps indicates if range1 and range2 are overlapped
func (z3 *Z3) Overlaps(range1 api.ZRange, range2 api.ZRange) bool {
	range1d0Min := z3.dim(0, range1.Min)
	range1d1Min := z3.dim(1, range1.Min)
	range1d2Min := z3.dim(2, range1.Min)
	range1d0Max := z3.dim(0, range1.Max)
	range1d1Max := z3.dim(1, range1.Max)
	range1d2Max := z3.dim(2, range1.Max)

	range2d0Min := z3.dim(0, range2.Min)
	range2d1Min := z3.dim(1, range2.Min)
	range2d2Min := z3.dim(2, range2.Min)
	range2d0Max := z3.dim(0, range2.Max)
	range2d1Max := z3.dim(1, range2.Max)
	range2d2Max := z3.dim(2, range2.Max)

	return (utils.MaxUint(range1d0Min, range2d0Min) <= utils.MinUint(range1d0Max, range2d0Max)) && (utils.MaxUint(range1d1Min, range2d1Min) <= utils.MinUint(range1d1Max, range2d1Max)) && (utils.MaxUint(range1d2Min, range2d2Min) <= utils.MinUint(range1d2Max, range2d2Max))
}

func (z3 *Z3) split(value uint64) uint64 {
	x := value & z3.MaxMask
	x = (x | x<<32) & 0x1f00000000ffff
	x = (x | x<<16) & 0x1f0000ff0000ff
	x = (x | x<<8) & 0x100f00f00f00f00f
	x = (x | x<<4) & 0x10c30c30c30c30c3
	x = (x | x<<2) & 0x1249249249249249
	return x
}

func (z3 *Z3) combine(z uint64) uint {
	x := z & 0x1249249249249249
	x = (x ^ (x >> 2)) & 0x10c30c30c30c30c3
	x = (x ^ (x >> 4)) & 0x100f00f00f00f00f
	x = (x ^ (x >> 8)) & 0x1f0000ff0000ff
	x = (x ^ (x >> 16)) & 0x1f00000000ffff
	x = (x ^ (x >> 32)) & z3.MaxMask
	return uint(x)
}

func (z3 *Z3) dim(i uint, zValue uint64) uint {
	return z3.combine(zValue >> i)
}

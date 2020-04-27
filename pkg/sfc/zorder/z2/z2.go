package z2

import (
	"math"
	"sync"

	"github.com/scraly/go.common/pkg/sfc/utils"
	api "github.com/scraly/go.common/pkg/sfc/zorder"
)

// Z2 struct
type Z2 struct {
	*sync.RWMutex

	BitsPerDimension int
	Dimensions       int
	MaxMask          uint64
	TotalBits        int
	Quadrants        int
	ZValue           uint64
}

// NewZ2 method constructs a Z2 struct and returns it as Z2N interface
func NewZ2() api.Z2N {
	return newInternalZ2(0)
}

// NewZ2WithZValue constructs a Z3 struct with initial z3 value and returns it as Z3N interface
func NewZ2WithZValue(z uint64) api.Z2N {
	return newInternalZ2(z)
}

func newInternalZ2(z uint64) api.Z2N {
	dimensions := 2
	return &Z2{
		RWMutex:          &sync.RWMutex{},
		BitsPerDimension: 31,
		Dimensions:       dimensions,
		MaxMask:          0x7fffffff,
		TotalBits:        62,
		Quadrants:        int(math.Pow(2, float64(dimensions))),
		ZValue:           z,
	}
}

// Apply method calculates a z2 from x,y
func (z2 *Z2) Apply(x, y uint) uint64 {
	z2.setZValue(z2.split(uint64(x)) | z2.split(uint64(y))<<1)
	return z2.GetZValue()
}

// UnApply method calculates x,y from z2
func (z2 *Z2) UnApply(z uint64) (uint, uint) {
	z2.setZValue(z)
	return z2.combine(z), z2.combine(z >> 1)
}

// GetDimensions returns the number of dimensions (2 for Z2)
func (z2 *Z2) GetDimensions() int {
	return z2.Dimensions
}

// GetQuadrants returns the number of quadrants (4 for Z2)
func (z2 *Z2) GetQuadrants() int {
	return z2.Quadrants
}

// GetTotalBits returns the total bits used to encode Z2
func (z2 *Z2) GetTotalBits() int {
	return z2.TotalBits
}

// GetZValue returns the computed Z2 value
func (z2 *Z2) GetZValue() uint64 {
	return z2.ZValue
}

func (z2 *Z2) setZValue(zValue uint64) {
	z2.Lock()
	defer z2.Unlock()

	z2.ZValue = zValue
}

// Contains indicates if z2 value is contained in rangeZ
func (z2 *Z2) Contains(rangeZ api.ZRange, value uint64) bool {
	x, y := z2.UnApply(value)
	d0Min := z2.dim(0, rangeZ.Min)
	d1Min := z2.dim(1, rangeZ.Min)
	d0Max := z2.dim(0, rangeZ.Max)
	d1Max := z2.dim(1, rangeZ.Max)

	return x >= d0Min && x <= d0Max && y >= d1Min && y <= d1Max
}

// Overlaps indicates if range1 and range2 are overlapped
func (z2 *Z2) Overlaps(range1 api.ZRange, range2 api.ZRange) bool {
	range1d0Min := z2.dim(0, range1.Min)
	range1d1Min := z2.dim(1, range1.Min)
	range1d0Max := z2.dim(0, range1.Max)
	range1d1Max := z2.dim(1, range1.Max)

	range2d0Min := z2.dim(0, range2.Min)
	range2d1Min := z2.dim(1, range2.Min)
	range2d0Max := z2.dim(0, range2.Max)
	range2d1Max := z2.dim(1, range2.Max)

	return (utils.MaxUint(range1d0Min, range2d0Min) <= utils.MinUint(range1d0Max, range2d0Max)) && (utils.MaxUint(range1d1Min, range2d1Min) <= utils.MinUint(range1d1Max, range2d1Max))
}

func (z2 *Z2) split(value uint64) uint64 {
	x := value & z2.MaxMask
	x = (x ^ (x << 32)) & 0x00000000ffffffff
	x = (x ^ (x << 16)) & 0x0000ffff0000ffff
	x = (x ^ (x << 8)) & 0x00ff00ff00ff00ff // 11111111000000001111111100000000..
	x = (x ^ (x << 4)) & 0x0f0f0f0f0f0f0f0f // 1111000011110000
	x = (x ^ (x << 2)) & 0x3333333333333333 // 11001100..
	x = (x ^ (x << 1)) & 0x5555555555555555 // 1010...
	return x
}

func (z2 *Z2) combine(z uint64) uint {
	x := z & 0x5555555555555555
	x = (x ^ (x >> 1)) & 0x3333333333333333
	x = (x ^ (x >> 2)) & 0x0f0f0f0f0f0f0f0f
	x = (x ^ (x >> 4)) & 0x00ff00ff00ff00ff
	x = (x ^ (x >> 8)) & 0x0000ffff0000ffff
	x = (x ^ (x >> 16)) & 0x00000000ffffffff
	return uint(x)
}

func (z2 *Z2) dim(i uint, zValue uint64) uint {
	return z2.combine(zValue >> i)
}

package zorder

// SpaceFillingCurve interface define methods for Z2 stuff
type SpaceFillingCurve interface {
	Index(x, y float64) (Z2N, error)
	Invert(z2 Z2N) (float64, float64)
}

// SpaceTimeFillingCurve interface define methods for Z3 stuff
type SpaceTimeFillingCurve interface {
	Index(x, y float64, t uint64) (Z3N, error)
	Invert(z3 Z3N) (float64, float64, uint64)
}

// Z2N interface defines methods applicable to Z2
type Z2N interface {
	ZN
	Apply(x, y uint) uint64
	UnApply(z uint64) (uint, uint)
	// ZDivide(p int64, rmin int64, rmax int64) (int64, int64)
}

// Z3N interface defines methods applicable to Z3
type Z3N interface {
	ZN
	Apply(x, y, z uint) uint64
	UnApply(z uint64) (uint, uint, uint)
	// ZDivide(p int64, rmin int64, rmax int64) (int64, int64)
}

// ZN interface defines common methods to Z2 and Z3
type ZN interface {
	GetDimensions() int
	GetQuadrants() int
	GetTotalBits() int
	GetZValue() uint64
	Contains(rangeZ ZRange, value uint64) bool
	Overlaps(range1 ZRange, range2 ZRange) bool
}

// Normalizer interface defines Normalize and DeNormalize methods
type Normalizer interface {
	Normalize(value float64) uint
	DeNormalize(n uint) float64
	GetMin() float64
	GetMax() float64
}

package zrange_test

import (
	"testing"

	"github.com/scraly/go.common/pkg/sfc"
	api "github.com/scraly/go.common/pkg/sfc/zorder"
	"github.com/scraly/go.common/pkg/sfc/zorder/z2"
	"github.com/scraly/go.common/pkg/sfc/zorder/z3"
	zranges "github.com/scraly/go.common/pkg/sfc/zorder/zrange"

	. "github.com/onsi/gomega"
)

func TestLongestCommonPrefix(t *testing.T) {

	tests := []struct {
		totalBits  int
		dimensions int
		values     []uint64
		prefix     uint64
		precision  int
		err        bool
	}{
		{
			totalBits:  63,
			dimensions: 3,
			values:     []uint64{8112745685912240673},
			prefix:     0,
			precision:  0,
			err:        true,
		},
		{
			totalBits:  63,
			dimensions: 3,
			values:     []uint64{8112745685912240673, 8770437881794063205},
			prefix:     8070450532247928832,
			precision:  4,
		},
		{
			totalBits:  63,
			dimensions: 3,
			values:     []uint64{4086666153209713566, 4086668626887099574, 4086668390556888644},
			prefix:     4086664818117836800,
			precision:  22,
		},
	}

	for _, tt := range tests {
		t.Run("LongestCommonPrefix", func(t *testing.T) {
			RegisterTestingT(t)

			prefix, precision, err := zranges.LongestCommonPrefix(tt.totalBits, tt.dimensions, tt.values)
			if tt.err {
				Expect(err).To(Not(BeNil()))
			} else {
				Expect(err).To(BeNil())
				Expect(prefix).To(Equal(tt.prefix))
				Expect(precision).To(Equal(tt.precision))
			}

		})
	}

}

func TestZ2IsContained(t *testing.T) {
	tests := []struct {
		x1min     uint
		y1min     uint
		x1max     uint
		y1max     uint
		x2min     uint
		y2min     uint
		x2max     uint
		y2max     uint
		contained bool
	}{
		{
			x1min:     0,
			y1min:     0,
			x1max:     2,
			y1max:     2,
			x2min:     1,
			y2min:     1,
			x2max:     3,
			y2max:     3,
			contained: false,
		},
		{
			x1min:     0,
			y1min:     0,
			x1max:     2,
			y1max:     2,
			x2min:     0,
			y2min:     0,
			x2max:     1,
			y2max:     1,
			contained: true,
		},
		{
			x1min:     0,
			y1min:     0,
			x1max:     1,
			y1max:     1,
			x2min:     2,
			y2min:     2,
			x2max:     3,
			y2max:     3,
			contained: false,
		},
	}

	for _, tt := range tests {
		t.Run("IsContained", func(t *testing.T) {
			RegisterTestingT(t)

			z2 := z2.NewZ2()
			zrange, _ := api.NewZRange(z2.Apply(tt.x1min, tt.y1min), z2.Apply(tt.x1max, tt.y1max))
			zrange1, _ := api.NewZRange(z2.Apply(tt.x2min, tt.y2min), z2.Apply(tt.x2max, tt.y2max))

			zbounds := []api.ZRange{*zrange}
			isContained := zranges.IsContained(z2, zbounds, *zrange1)

			Expect(isContained).To(Equal(tt.contained))

		})
	}
}

func TestZ2IsOverlapped(t *testing.T) {
	tests := []struct {
		x1min      uint
		y1min      uint
		x1max      uint
		y1max      uint
		x2min      uint
		y2min      uint
		x2max      uint
		y2max      uint
		overlapped bool
	}{
		{
			x1min:      0,
			y1min:      0,
			x1max:      2,
			y1max:      2,
			x2min:      1,
			y2min:      1,
			x2max:      3,
			y2max:      3,
			overlapped: true,
		},
		{
			x1min:      0,
			y1min:      0,
			x1max:      2,
			y1max:      2,
			x2min:      0,
			y2min:      0,
			x2max:      1,
			y2max:      1,
			overlapped: true,
		},
		{
			x1min:      0,
			y1min:      0,
			x1max:      1,
			y1max:      1,
			x2min:      2,
			y2min:      2,
			x2max:      3,
			y2max:      3,
			overlapped: false,
		},
	}

	for _, tt := range tests {
		t.Run("IsOverlapped", func(t *testing.T) {
			RegisterTestingT(t)

			z2 := z2.NewZ2()
			zrange, _ := api.NewZRange(z2.Apply(tt.x1min, tt.y1min), z2.Apply(tt.x1max, tt.y1max))
			zrange1, _ := api.NewZRange(z2.Apply(tt.x2min, tt.y2min), z2.Apply(tt.x2max, tt.y2max))

			zbounds := []api.ZRange{*zrange}
			isOverlapped := zranges.IsOverlapped(z2, zbounds, *zrange1)

			Expect(isOverlapped).To(Equal(tt.overlapped))

		})
	}
}

type Z3Values struct {
	x uint
	y uint
	z uint
}

type Z3IndexRangeValues struct {
	lower Z3Values
	upper Z3Values
}

func TestCalculateRangesZ3(t *testing.T) {
	tests := []struct {
		x1           uint
		y1           uint
		z1           uint
		x2           uint
		y2           uint
		z2           uint
		count        int
		ranges       []sfc.IndexRange
		rangesValues []Z3IndexRangeValues
	}{
		{
			x1:    2,
			y1:    2,
			z1:    0,
			x2:    3,
			y2:    6,
			z2:    0,
			count: 3,
			ranges: []sfc.IndexRange{
				{Lower: 24, Upper: 27, Contained: true},
				{Lower: 136, Upper: 139, Contained: true},
				{Lower: 152, Upper: 153, Contained: true},
			},
			rangesValues: []Z3IndexRangeValues{
				{lower: Z3Values{x: 2, y: 2, z: 0}, upper: Z3Values{x: 3, y: 3, z: 0}},
				{lower: Z3Values{x: 2, y: 4, z: 0}, upper: Z3Values{x: 3, y: 5, z: 0}},
				{lower: Z3Values{x: 2, y: 6, z: 0}, upper: Z3Values{x: 3, y: 6, z: 0}},
			},
		},
	}

	for _, tt := range tests {
		t.Run("CalculateRanges", func(t *testing.T) {
			RegisterTestingT(t)

			z3 := z3.NewZ3()
			zrange, _ := api.NewZRange(z3.Apply(tt.x1, tt.y1, tt.z1), z3.Apply(tt.x2, tt.y2, tt.z2))

			zbounds := []api.ZRange{*zrange}

			result, err := zranges.CalculateRanges(z3, zbounds, 64, 0, 7)

			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(tt.count))

			Expect(result[0].Lower).To(Equal(tt.ranges[0].Lower))
			Expect(result[0].Upper).To(Equal(tt.ranges[0].Upper))
			Expect(result[1].Lower).To(Equal(tt.ranges[1].Lower))
			Expect(result[1].Upper).To(Equal(tt.ranges[1].Upper))
			Expect(result[2].Lower).To(Equal(tt.ranges[2].Lower))
			Expect(result[2].Upper).To(Equal(tt.ranges[2].Upper))

			x, y, z := z3.UnApply(result[0].Lower)
			Expect(x).To(Equal(tt.rangesValues[0].lower.x))
			Expect(y).To(Equal(tt.rangesValues[0].lower.y))
			Expect(z).To(Equal(tt.rangesValues[0].lower.z))

			x, y, z = z3.UnApply(result[0].Upper)
			Expect(x).To(Equal(tt.rangesValues[0].upper.x))
			Expect(y).To(Equal(tt.rangesValues[0].upper.y))
			Expect(z).To(Equal(tt.rangesValues[0].upper.z))

			x, y, z = z3.UnApply(result[1].Lower)
			Expect(x).To(Equal(tt.rangesValues[1].lower.x))
			Expect(y).To(Equal(tt.rangesValues[1].lower.y))
			Expect(z).To(Equal(tt.rangesValues[1].lower.z))

			x, y, z = z3.UnApply(result[1].Upper)
			Expect(x).To(Equal(tt.rangesValues[1].upper.x))
			Expect(y).To(Equal(tt.rangesValues[1].upper.y))
			Expect(z).To(Equal(tt.rangesValues[1].upper.z))

			x, y, z = z3.UnApply(result[2].Lower)
			Expect(x).To(Equal(tt.rangesValues[2].lower.x))
			Expect(y).To(Equal(tt.rangesValues[2].lower.y))
			Expect(z).To(Equal(tt.rangesValues[2].lower.z))

			x, y, z = z3.UnApply(result[2].Upper)
			Expect(x).To(Equal(tt.rangesValues[2].upper.x))
			Expect(y).To(Equal(tt.rangesValues[2].upper.y))
			Expect(z).To(Equal(tt.rangesValues[2].upper.z))

		})
	}
}

func BenchmarkCalculateRanges(b *testing.B) {
	z3 := z3.NewZ3()
	zrange, _ := api.NewZRange(z3.Apply(2, 2, 0), z3.Apply(3, 6, 0))

	zbounds := []api.ZRange{*zrange}

	for n := 0; n < b.N; n++ {
		zranges.CalculateRanges(z3, zbounds, 64, 0, 15)
	}
}

type Z2Values struct {
	x uint
	y uint
}

type Z2IndexRangeValues struct {
	lower Z2Values
	upper Z2Values
}

func TestCalculateRangesZ2(t *testing.T) {
	tests := []struct {
		x1           uint
		y1           uint
		x2           uint
		y2           uint
		count        int
		ranges       []sfc.IndexRange
		rangesValues []Z2IndexRangeValues
	}{
		{
			x1:    2,
			y1:    2,
			x2:    3,
			y2:    6,
			count: 3,
			ranges: []sfc.IndexRange{
				{Lower: 12, Upper: 15, Contained: true},
				{Lower: 36, Upper: 39, Contained: true},
				{Lower: 44, Upper: 45, Contained: true},
			},
			rangesValues: []Z2IndexRangeValues{
				{lower: Z2Values{x: 2, y: 2}, upper: Z2Values{x: 3, y: 3}},
				{lower: Z2Values{x: 2, y: 4}, upper: Z2Values{x: 3, y: 5}},
				{lower: Z2Values{x: 2, y: 6}, upper: Z2Values{x: 3, y: 6}},
			},
		},
		{
			x1:    0,
			y1:    2,
			x2:    3,
			y2:    5,
			count: 2,
			ranges: []sfc.IndexRange{
				{Lower: 8, Upper: 15, Contained: true},
				{Lower: 32, Upper: 39, Contained: true},
			},
			rangesValues: []Z2IndexRangeValues{
				{lower: Z2Values{x: 0, y: 2}, upper: Z2Values{x: 3, y: 3}},
				{lower: Z2Values{x: 0, y: 4}, upper: Z2Values{x: 3, y: 5}},
			},
		},
	}

	for _, tt := range tests {
		t.Run("CalculateRanges", func(t *testing.T) {
			RegisterTestingT(t)

			z2 := z2.NewZ2()
			zmin := z2.Apply(tt.x1, tt.y1)
			zmax := z2.Apply(tt.x2, tt.y2)

			zrange, _ := api.NewZRange(zmin, zmax)

			zbounds := []api.ZRange{*zrange}

			result, err := zranges.CalculateRanges(z2, zbounds, 64, 0, 7)

			Expect(err).To(BeNil())
			Expect(len(result)).To(Equal(tt.count))

			for i := 0; i < len(result); i++ {
				Expect(result[i].Lower).To(Equal(tt.ranges[i].Lower))
				Expect(result[i].Upper).To(Equal(tt.ranges[i].Upper))

				x, y := z2.UnApply(result[i].Lower)
				Expect(x).To(Equal(tt.rangesValues[i].lower.x))
				Expect(y).To(Equal(tt.rangesValues[i].lower.y))

				x, y = z2.UnApply(result[i].Upper)
				Expect(x).To(Equal(tt.rangesValues[i].upper.x))
				Expect(y).To(Equal(tt.rangesValues[i].upper.y))
			}

		})
	}
}

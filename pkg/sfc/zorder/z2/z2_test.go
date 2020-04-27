package z2_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/scraly/go.common/pkg/sfc"
	"github.com/scraly/go.common/pkg/sfc/zorder"
	"github.com/scraly/go.common/pkg/sfc/zorder/normalizer"
	"github.com/scraly/go.common/pkg/sfc/zorder/z2"
)

func TestZ2(t *testing.T) {

	tests := []struct {
		lat       float64
		lon       float64
		precision uint
		z2        uint64
		x         uint
		y         uint
	}{
		{
			lat:       43.56362116,
			lon:       1.391483667,
			precision: 31,
			z2:        3650378452887927909,
			x:         1082042347,
			y:         1593476068,
		},
		{
			lat:       43.563621169,
			lon:       1.3914836678,
			precision: 31,
			z2:        3650378452887927911,
			x:         1082042347,
			y:         1593476069,
		},
	}

	for _, tt := range tests {
		t.Run("z2", func(t *testing.T) {
			RegisterTestingT(t)

			bNormLon, _ := normalizer.NewNormalizer(-180, 180, tt.precision)
			bNormLat, _ := normalizer.NewNormalizer(-90, 90, tt.precision)

			columnLon := bNormLon.Normalize(tt.lon)
			lineLat := bNormLat.Normalize(tt.lat)

			z2Obj := z2.NewZ2()

			z2 := z2Obj.Apply(uint(columnLon), uint(lineLat))

			Expect(z2).To(Equal(tt.z2))

			x, y := z2Obj.UnApply(z2)

			Expect(x).To(Equal(tt.x))
			Expect(y).To(Equal(tt.y))
		})
	}

}

func TestZ2Sfc(t *testing.T) {

	tests := []struct {
		lat       float64
		lon       float64
		precision uint
		z2        uint64
		x         uint
		y         uint
	}{
		{
			lat:       43.56362116,
			lon:       1.391483667,
			precision: 31,
			z2:        3650378452887927909,
			x:         1082042347,
			y:         1593476068,
		},
		{
			lat:       43.563621169,
			lon:       1.3914836678,
			precision: 31,
			z2:        3650378452887927911,
			x:         1082042347,
			y:         1593476069,
		},
	}

	for _, tt := range tests {
		t.Run("z2Sfc", func(t *testing.T) {
			RegisterTestingT(t)

			z2Sfc, errZ2Sfc := z2.NewZ2Sfc(uint(31))

			Expect(errZ2Sfc).To(BeNil())

			z2N, errIndex := z2Sfc.Index(tt.lon, tt.lat)

			Expect(errIndex).To(BeNil())

			Expect(z2N.GetZValue()).To(Equal(tt.z2))

			x, y := z2Sfc.Invert(z2N)

			Expect(x).Should(BeNumerically("~", tt.lon, 0.0001))
			Expect(y).Should(BeNumerically("~", tt.lat, 0.0001))

		})
	}

}

func TestSearch(t *testing.T) {
	tests := []struct {
		southWestLat float64
		southWestLon float64
		northEastLat float64
		northEastLon float64
		count        int
		stepxy       float64
		err          bool
	}{
		{
			southWestLat: 43.5,
			southWestLon: 1.4,
			northEastLat: 44,
			northEastLon: 1.5,
			count:        50,
			stepxy:       0.1,
			err:          false,
		},
		{
			southWestLat: 43.52,
			southWestLon: -0.83,
			northEastLat: 47.44,
			northEastLon: 5.29,
			count:        6,
			stepxy:       0.1,
			err:          false,
		},
	}

	for _, tt := range tests {
		t.Run("search", func(t *testing.T) {
			RegisterTestingT(t)

			search, _ := z2.NewSearch(z2.NormalizerMaxPrecision)
			bbox := sfc.BoundingBox{
				SouthWest: sfc.Point{Longitude: tt.southWestLon, Latitude: tt.southWestLat},
				NorthEast: sfc.Point{Longitude: tt.northEastLon, Latitude: tt.northEastLat},
			}
			result, err := search.GetZ2Ranges(bbox)

			if tt.err {
				Expect(err).To(Not(BeNil()))
			} else {
				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(tt.count))
			}

			sfc := search.GetSpaceFillingCurve()
			for x := tt.southWestLon; x <= tt.northEastLon; x = x + tt.stepxy {
				for y := tt.southWestLat; y <= tt.northEastLat; y = y + tt.stepxy {
					z2n, _ := sfc.Index(x, y)
					Expect(isInRange(z2n, result)).To(BeTrue())
				}
			}

		})
	}
}

func isInRange(z2n zorder.Z2N, ranges []*sfc.IndexRange) bool {

	for _, indexRange := range ranges {
		if indexRange.Lower <= z2n.GetZValue() && z2n.GetZValue() <= indexRange.Upper {
			return true
		}
	}

	return false
}

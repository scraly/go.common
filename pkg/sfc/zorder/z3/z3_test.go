package z3_test

import (
	"math"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/scraly/go.common/pkg/sfc"
	"github.com/scraly/go.common/pkg/sfc/utils"
	"github.com/scraly/go.common/pkg/sfc/zorder"
	"github.com/scraly/go.common/pkg/sfc/zorder/normalizer"
	"github.com/scraly/go.common/pkg/sfc/zorder/z3"
)

func TestZ3(t *testing.T) {

	tests := []struct {
		lat       float64
		lon       float64
		time      int64
		precision uint
		z3        uint64
		x         uint
		y         uint
		zTime     uint
	}{
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      2457,
			precision: 21,
			z3:        3499933763342795557,
			x:         1056681,
			y:         1556128,
			zTime:     8519,
		},
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      546789,
			precision: 21,
			z3:        8760296496909503265,
			x:         1056681,
			y:         1556128,
			zTime:     1895998,
		},
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      324000,
			precision: 21,
			z3:        8112745685912240673,
			x:         1056681,
			y:         1556128,
			zTime:     1123474,
		},
		{
			lat:       43.606189,
			lon:       1.447046,
			time:      604800,
			precision: 21,
			z3:        8770437881794063205,
			x:         1057005,
			y:         1556624,
			zTime:     2097151,
		},
	}

	for _, tt := range tests {
		t.Run("z3", func(t *testing.T) {
			RegisterTestingT(t)

			bNormLon, _ := normalizer.NewNormalizer(-180, 180, tt.precision)
			bNormLat, _ := normalizer.NewNormalizer(-90, 90, tt.precision)
			bTime, _ := normalizer.NewNormalizer(0, 604800, tt.precision)

			columnLon := bNormLon.Normalize(tt.lon)
			lineLat := bNormLat.Normalize(tt.lat)
			nTime := bTime.Normalize(float64(tt.time))

			z3Obj := z3.NewZ3()

			z3 := z3Obj.Apply(uint(columnLon), uint(lineLat), uint(nTime))

			Expect(z3).To(Equal(tt.z3))

			x, y, time := z3Obj.UnApply(z3)

			Expect(x).To(Equal(tt.x))
			Expect(y).To(Equal(tt.y))
			Expect(time).To(Equal(tt.zTime))

			Expect(bNormLon.DeNormalize(uint(x))).Should(BeNumerically("~", tt.lon, 0.0001))
			Expect(bNormLat.DeNormalize(uint(y))).Should(BeNumerically("~", tt.lat, 0.0001))
			Expect(bTime.DeNormalize(uint(time))).Should(BeNumerically("~", tt.time, 0.2))
		})
	}
}

func TestZ3Sfc(t *testing.T) {

	tests := []struct {
		lat       float64
		lon       float64
		time      uint64
		precision uint
		z3        uint64
		x         uint
		y         uint
		zTime     uint
	}{
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      2457,
			precision: 21,
			z3:        3499933763342795557,
			x:         1056681,
			y:         1556128,
			zTime:     8519,
		},
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      546789,
			precision: 21,
			z3:        8760296496909503265,
			x:         1056681,
			y:         1556128,
			zTime:     1895998,
		},
		{
			lat:       43.56362116,
			lon:       1.391483667,
			time:      324000,
			precision: 21,
			z3:        8112745685912240673,
			x:         1056681,
			y:         1556128,
			zTime:     1123474,
		},
		{
			lat:       43.606189,
			lon:       1.447046,
			time:      604800,
			precision: 21,
			z3:        8770437881794063205,
			x:         1057005,
			y:         1556624,
			zTime:     2097151,
		},
	}

	for _, tt := range tests {
		t.Run("z3Sfc", func(t *testing.T) {
			RegisterTestingT(t)

			z3Sfc, errZ3Sfc := z3.NewZ3Sfc(z3.NormalizerMaxPrecision)

			Expect(errZ3Sfc).To(BeNil())

			z3N, errIndex := z3Sfc.Index(tt.lon, tt.lat, tt.time)

			Expect(errIndex).To(BeNil())

			Expect(z3N.GetZValue()).To(Equal(tt.z3))

			x, y, time := z3Sfc.Invert(z3N)

			Expect(x).Should(BeNumerically("~", tt.lon, 0.0001))
			Expect(y).Should(BeNumerically("~", tt.lat, 0.0001))
			Expect(time).Should(BeNumerically("~", tt.time, 0.2))
		})
	}

}

func TestSearch(t *testing.T) {
	tests := []struct {
		southWestLat float64
		southWestLon float64
		northEastLat float64
		northEastLon float64
		dateMin      time.Time
		dateMax      time.Time
		count        int
		stepxy       float64
		stepTime     float64
		err          bool
	}{
		{
			southWestLat: 43.5,
			southWestLon: 1.4,
			northEastLat: 44,
			northEastLon: 1.5,
			dateMin:      time.Date(2018, 9, 10, 7, 0, 0, 0, time.UTC),
			dateMax:      time.Date(2018, 9, 16, 23, 0, 0, 0, time.UTC),
			count:        246,
			stepxy:       0.1,
			stepTime:     30,
			err:          false,
		},
		// {
		// 	southWestLat: 43.52,
		// 	southWestLon: -0.83,
		// 	northEastLat: 47.44,
		// 	northEastLon: 5.29,
		// 	dateMin:      time.Date(2018, 6, 10, 7, 0, 0, 0, time.UTC),
		// 	dateMax:      time.Date(2018, 12, 16, 23, 0, 0, 0, time.UTC),
		// 	count:        640,
		// 	stepxy:       0.05,
		// 	stepTime:     30000,
		// 	err:          false,
		// },
		{
			southWestLat: 43.52,
			southWestLon: -0.83,
			northEastLat: 47.44,
			northEastLon: 5.29,
			dateMin:      time.Date(2018, 6, 10, 7, 0, 0, 0, time.UTC),
			dateMax:      time.Date(2018, 6, 11, 7, 0, 0, 0, time.UTC),
			count:        96,
			stepxy:       0.1,
			stepTime:     300,
			err:          false,
		},
	}

	for _, tt := range tests {
		t.Run("search", func(t *testing.T) {
			RegisterTestingT(t)

			search, _ := z3.NewSearch()
			bbox := sfc.BoundingBox{
				SouthWest: sfc.Point{Longitude: tt.southWestLon, Latitude: tt.southWestLat},
				NorthEast: sfc.Point{Longitude: tt.northEastLon, Latitude: tt.northEastLat},
			}
			result, _, err := search.GetZ3Ranges(bbox, tt.dateMin, tt.dateMax)

			if tt.err {
				Expect(err).To(Not(BeNil()))
			} else {
				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(tt.count))
			}

			curve := search.GetSpaceTimeFillingCurve()

			weekTimeRanges, _ := utils.GetWeekTimeRangeFromDateRange(tt.dateMin, tt.dateMax)

			stfc := search.GetSpaceTimeFillingCurve()
			for x := tt.southWestLon; x <= tt.northEastLon; x = x + tt.stepxy {
				for y := tt.southWestLat; y <= tt.northEastLat; y = y + tt.stepxy {
					for _, weekTimeRange := range weekTimeRanges {

						for s := weekTimeRange.MinWeekDate.Seconds; s <= weekTimeRange.MaxWeekDate.Seconds; s = s + tt.stepTime {
							z3n, _ := stfc.Index(x, y, uint64(math.Round(s)))
							Expect(isInRange(z3n, curve, result)).To(BeTrue())
						}
					}
				}
			}

		})
	}
}

func TestZ3Loose(t *testing.T) {
	RegisterTestingT(t)

	lat := 43.606084
	lon := 1.441322
	// Monday 0 second to 86400 every week
	seconds := uint64(86400)

	z3Sfc, errZ3Sfc := z3.NewZ3SfcWithPrecision(7, 7, 3)

	Expect(errZ3Sfc).To(BeNil())

	z3N, errIndex := z3Sfc.Index(lon, lat, seconds)

	Expect(errIndex).To(BeNil())
	Expect(z3N.GetZValue()).To(Equal(uint64(795798)))

	z3N, _ = z3Sfc.Index(0, 43.59375, 86400)
	Expect(z3N.GetZValue()).To(Equal(uint64(795798)))

	z3N, _ = z3Sfc.Index(2.8124, 44.99, 86400)
	Expect(z3N.GetZValue()).To(Equal(uint64(795798)))

}

func BenchmarkGetZ3Ranges(b *testing.B) {
	b.ReportAllocs()

	search, _ := z3.NewSearch()
	bbox := sfc.BoundingBox{
		SouthWest: sfc.Point{Longitude: 1.4, Latitude: 43.5},
		NorthEast: sfc.Point{Longitude: 1.5, Latitude: 44},
	}

	for n := 0; n < b.N; n++ {
		search.GetZ3Ranges(bbox, time.Date(2018, 9, 10, 7, 0, 0, 0, time.UTC), time.Date(2018, 9, 16, 23, 0, 0, 0, time.UTC))
	}
}

func isInRange(z3n zorder.Z3N, stfc zorder.SpaceTimeFillingCurve, ranges []*sfc.IndexRange) bool {

	for _, indexRange := range ranges {
		if indexRange.Lower <= z3n.GetZValue() && z3n.GetZValue() <= indexRange.Upper {
			return true
		}
	}

	return false
}

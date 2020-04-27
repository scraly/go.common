package normalizer_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/scraly/go.common/pkg/sfc/zorder/normalizer"
)

func TestNormalize(t *testing.T) {

	tests := []struct {
		lat       float64
		lon       float64
		lineLat   uint
		columnLon uint
		precision uint
	}{
		{
			lat:       43.56362116,
			lon:       1.391483667,
			lineLat:   1593476068,
			columnLon: 1082042347,
			precision: 31,
		},
		{
			lat:       -43.56362116,
			lon:       -1.391483667,
			lineLat:   554007579,
			columnLon: 1065441300,
			precision: 31,
		},
		{
			lat:       -80.56565345,
			lon:       -179.35345457,
			lineLat:   112556138,
			columnLon: 3856793,
			precision: 31,
		},
		{
			lat:       -80.56565345,
			lon:       -179.35345457,
			lineLat:   0,
			columnLon: 0,
			precision: 2,
		},
		{
			lat:       30.56565345,
			lon:       -179.35345457,
			lineLat:   2,
			columnLon: 0,
			precision: 2,
		},
		{
			lat:       30.56565345,
			lon:       -179.35345457,
			lineLat:   5,
			columnLon: 0,
			precision: 3,
		},
		{
			lat:       -80.56565345,
			lon:       -179.35345457,
			lineLat:   1,
			columnLon: 0,
			precision: 5,
		},
		{
			lat:       -80.56565345,
			lon:       -179.35345457,
			lineLat:   53,
			columnLon: 1,
			precision: 10,
		},
	}

	for _, tt := range tests {
		t.Run("Normalizer", func(t *testing.T) {
			RegisterTestingT(t)

			bNormLon, _ := normalizer.NewNormalizer(-180, 180, tt.precision)
			bNormLat, _ := normalizer.NewNormalizer(-90, 90, tt.precision)

			columnLon := bNormLon.Normalize(tt.lon)
			lineLat := bNormLat.Normalize(tt.lat)

			Expect(columnLon).To(Equal(tt.columnLon))
			Expect(lineLat).To(Equal(tt.lineLat))
		})
	}

	//Expect(lon1).Should((BeNumerically("~", lon, 0.0000001)))

}

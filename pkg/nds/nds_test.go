package nds_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/scraly/go.common/pkg/nds"
)

func TestCalculatePackedTileID(t *testing.T) {
	tests := []struct {
		label             string
		lon               float64
		lat               float64
		packedTileIDLvl13 int64
	}{
		{
			label:             "Eiffel Tower",
			lon:               2.2945,
			lat:               48.858222,
			packedTileIDLvl13: 545299690,
		},
		{
			label:             "Near Quito",
			lon:               -78.45,
			lat:               0.0,
			packedTileIDLvl13: 621019217,
		},
		{
			label:             "Statue of Liberty",
			lon:               -74.044444,
			lat:               40.689167,
			packedTileIDLvl13: 623795102,
		},
		{
			label:             "Near the Millennium Dome",
			lon:               0,
			lat:               51.503,
			packedTileIDLvl13: 545392682,
		},
		{
			label:             "Sugarloaf Mountain",
			lon:               -43.157444,
			lat:               -22.948658,
			packedTileIDLvl13: 667597199,
		},
	}

	for _, tt := range tests {
		t.Run("CalculatePackedTileID", func(t *testing.T) {
			RegisterTestingT(t)

			packedTileIDLvl13, errLvl13 := nds.CalculatePackedTileID(tt.lon, tt.lat, 13)
			Expect(errLvl13).To(BeNil())
			Expect(packedTileIDLvl13).To(Equal(tt.packedTileIDLvl13))

			calcLon, calcLat, calcLevel := nds.ReadPackedTileID(packedTileIDLvl13)
			Expect(calcLevel).To(Equal(uint(13)))
			Expect(calcLon).Should(BeNumerically("~", tt.lon, 0.3))
			Expect(calcLat).Should(BeNumerically("~", tt.lat, 0.3))
		})
	}
}

func BenchmarkCalculatePackedTileID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nds.CalculatePackedTileID(-43.157444, -22.948658, 13)
	}
}

func BenchmarkReadPackedTileID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nds.ReadPackedTileID(667597199)
	}
}

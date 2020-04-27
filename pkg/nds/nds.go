package nds

import (
	"errors"
	"math"
	"math/bits"
)

// CalculatePackedTileID returns NDS packed tile id for lat,lon wgs84 coordinates at level passed as input
func CalculatePackedTileID(lon, lat float64, level uint) (int64, error) {

	if lon < -180 || lon > 180 {
		return 0, errors.New("lon must be in the range [-180, 180]")
	}

	if lat < -90 || lat > 90 {
		return 0, errors.New("lat must be in the range [-90, 90]")
	}

	if level < 0 || level > 15 {
		return 0, errors.New("level must be in the range [0, 15]")
	}

	morton := wgs2morton(lon, lat)
	tileID := calculateTileID(morton, level)

	return int64(1<<(16+level)) | tileID, nil
}

// ReadPackedTileID estracts longitude latitude and tile level
func ReadPackedTileID(packedTileID int64) (lon, lat float64, level uint) {
	level = uint(15 - (bits.LeadingZeros64(uint64(packedTileID)) - 32))
	tileBits := int64(1<<(2*level+1)) - 1
	tileID := packedTileID & tileBits

	lon, lat = morton2wgs(tileID, level)

	return lon, lat, level
}

func calculateTileID(morton int64, level uint) int64 {
	shift := 64 - (2*level + 1) - 1
	return morton >> shift
}

func wgs2morton(lon, lat float64) int64 {
	return nds2Morton(normalize(lon, -180, 180, 32), normalize(lat, -90, 90, 31))
}

func morton2wgs(morton int64, level uint) (lon, lat float64) {
	shift := 64 - (2*level + 1) - 1
	realmorton := morton << shift

	x, y := morton2nds(realmorton)
	lon = denormalize(x, -180, 180, 32)
	lat = denormalize(y, -90, 90, 31)
	return lon, lat
}

func nds2Morton(x, y int) int64 {
	y = y & 0x000000007fffffff // Remove bit 31 for y
	return split(int64(x)) | split(int64(y))<<1
}

func morton2nds(morton int64) (x, y int) {
	x = combine(morton)
	y = combine(morton >> 1)

	return x, y
}

func normalize(value, min, max float64, precision uint) int {
	bins := 1 << precision
	normalizer := float64(bins) / (max - min)

	return int(math.Trunc((value) * normalizer))
}

func denormalize(value int, min, max float64, precision uint) float64 {
	bins := 1 << precision
	normalizer := float64(bins) / (max - min)

	denormalizeVal := float64(value) / normalizer
	if denormalizeVal > max {
		denormalizeVal -= (max - min)
	}

	if denormalizeVal < min {
		denormalizeVal += (max - min)
	}

	return denormalizeVal
}

func split(value int64) int64 {
	x := value & 0xffffffff
	x = (x ^ (x << 32)) & 0x00000000ffffffff
	x = (x ^ (x << 16)) & 0x0000ffff0000ffff
	x = (x ^ (x << 8)) & 0x00ff00ff00ff00ff // 11111111000000001111111100000000..
	x = (x ^ (x << 4)) & 0x0f0f0f0f0f0f0f0f // 1111000011110000
	x = (x ^ (x << 2)) & 0x3333333333333333 // 11001100..
	x = (x ^ (x << 1)) & 0x5555555555555555 // 1010...
	return x
}

func combine(z int64) int {
	x := z & 0x5555555555555555
	x = (x ^ (x >> 1)) & 0x3333333333333333
	x = (x ^ (x >> 2)) & 0x0f0f0f0f0f0f0f0f
	x = (x ^ (x >> 4)) & 0x00ff00ff00ff00ff
	x = (x ^ (x >> 8)) & 0x0000ffff0000ffff
	x = (x ^ (x >> 16)) & 0x00000000ffffffff
	return int(x)
}

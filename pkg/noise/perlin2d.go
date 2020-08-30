package noise

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
)

// Perlin2D creates a 2D image with the specified number of octaves and persistance
func Perlin2D(width, height, octaves int, persistance float32) []uint8 {
	// determine maximum value of highest frequency
	repeat := int(cgm.Pow32(2, float32(octaves-1)))

	// setup perlin util
	perlin := makeperlin(repeat)

	// calc random value for each pixel
	var data []uint8
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float32(x) / float32(width-1)
			fy := float32(y) / float32(height-1)
			rnd := perlin.octave(fx, fy, 0, octaves, persistance)
			val := cgm.Map(rnd, 0, 1, 0, 255)
			data = append(data, uint8(val))
		}
	}

	return data
}

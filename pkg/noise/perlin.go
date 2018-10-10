package noise

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
)

func Perlin2D(width, height, res int) []uint8 {
	// setup perlin util
	perlin := makeperlinutil(-1)

	// calc random value for each pixel
	var data []uint8
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rnd := perlin.octave(float32(x)/float32(width), float32(y)/float32(height), 1, res, 1)
			val := cgm.Map(rnd, 0, 1, 0, 255)
			data = append(data, uint8(val))
		}
	}

	return data
}

func Perlin3D(width, height, slices, res int) []uint8 {
	// setup perlin util
	perlin := makeperlinutil(-1)

	// calc random value for each pixel
	var data []uint8
	for z := 0; z < slices; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rnd := perlin.octave(float32(x)/float32(width), float32(y)/float32(height), float32(z)/float32(slices), res, 1)
				val := cgm.Map(rnd, 0, 1, 0, 255)
				data = append(data, uint8(val))
			}
		}
	}

	return data
}

package noise

import "github.com/adrianderstroff/realtime-clouds/pkg/cgm"

func Perlin3D(width, height, slices, res, persistance int) []uint8 {
	// setup perlin util
	perlin := makeperlin(4)

	// calc random value for each pixel
	var data []uint8
	for z := 0; z < slices; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rnd := perlin.octave(float32(x)/float32(width-1), float32(y)/float32(height-1), float32(z)/float32(slices-1), res, float32(persistance))
				val := cgm.Map(rnd, 0, 1, 0, 255)
				data = append(data, uint8(val))
			}
		}
	}

	return data
}

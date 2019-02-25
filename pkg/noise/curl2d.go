package noise

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

func Curl2D(width, height, res int) []uint8 {
	// create 2D perlin noise
	perlin := Perlin2D(width, height, res)

	// calculate curl
	var maxval float32 = 0
	var values []float32
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// calculate curl
			a := deriveX(perlin, x, y, width, height)
			b := deriveY(perlin, x, y, width, height)
			curl := mgl32.Vec2{a, -b}

			// add magnitude
			val := curl.Len()
			values = append(values, val)

			// save maximum value for following normalization
			maxval = cgm.Max32(maxval, val)
		}
	}

	// map distance to 0..255 and save in data slice
	var data []uint8
	for i := 0; i < len(values); i++ {
		val := cgm.Map(values[i], 0, maxval, 0, 255)
		data = append(data, uint8(val))
	}

	return data
}

func deriveX(data []uint8, x, y, width, height int) float32 {
	yp := loop(y+1, height)
	yn := loop(y-1, height)
	return (f(data, x, yp, width) - f(data, x, yn, width)) / 2.0
}
func deriveY(data []uint8, x, y, width, height int) float32 {
	xp := loop(y+1, width)
	xn := loop(y-1, width)
	return (f(data, xp, y, width) - f(data, xn, y, width)) / 2.0
}

func f(data []uint8, x, y, w int) float32 {
	return float32(data[y*w+x]) / 255.0
}

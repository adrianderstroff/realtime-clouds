package noise

import (
	"math/rand"
	"time"

	"github.com/adrianderstroff/realtime-clouds/pkg/math"
)

func linspread(x, min, max float32) float32 {
	return (x - min) / (max - min)
}

func clamp(x, min, max float32) float32 {
	return math.Min32(math.Max32(x, min), max)
}

func dist2(x1, y1, x2, y2 float32) float32 {
	dx := x1 - x2
	dy := y1 - y2
	return math.Sqrt32(dx*dx + dy*dy)
}

func distToPoint(x1, y1, x2, y2 float32) float32 {
	mindist := float32(1)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			dist := dist2(x1, y1, x2+float32(dx), y2+float32(dy))
			mindist = math.Min32(mindist, dist)
		}
	}
	return mindist
}

// MakeWorley creates worley noise of the size specified by width and height
// with the specified number of points.
func MakeWorley(width, height, points int) []uint8 {
	data := []uint8{}

	// randomly set points
	xs := []float32{}
	ys := []float32{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < points; i++ {
		xs = append(xs, r.Float32())
		ys = append(ys, r.Float32())
	}

	// calc shortest distance to point
	globmin, globmax := float32(1), float32(0)
	temp := []float32{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mindist := float32(1)
			for p := 0; p < points; p++ {
				x1, y1 := float32(x)/float32(width), float32(y)/float32(height)
				dist := distToPoint(x1, y1, xs[p], ys[p])
				mindist = math.Min32(mindist, dist)

				// keep track of global values
				globmin = math.Min32(globmin, mindist)
				globmax = math.Max32(globmax, mindist)
			}
			temp = append(temp, mindist)
		}
	}

	// update data
	for i := 0; i < len(temp); i++ {
		// local spread
		//val := linspread(temp[i], globmin, globmax)
		val := temp[i]
		val = clamp(1-val-0.8, 0, 1)
		col := uint8(val * 255)
		// rgb channel
		data = append(data, col)
		data = append(data, col)
		data = append(data, col)
	}

	return data
}

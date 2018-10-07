package noise

import (
	"math/rand"
	"time"

	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

func Perlin3D(length, width, height, res int) []uint8 {
	// divide volume into cells
	xstep := float32(length) / float32(res)
	ystep := float32(height) / float32(res)
	zstep := float32(width) / float32(res)

	// build random vector grid
	vectors := make([][][]mgl32.Vec3, res)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for y := 0; y < res; y++ {
		vectors[y] = make([][]mgl32.Vec3, res)
		for z := 0; z < res; z++ {
			vectors[y][z] = make([]mgl32.Vec3, res)
			for x := 0; x < res; x++ {
				randomvector := mgl32.Vec3{r.Float32(), r.Float32(), r.Float32()}
				vectors[y][z][x] = randomvector.Normalize()
			}
		}
	}

	// calc random value for each pixel
	var maxval float32 = 0
	var values []float32
	for y := 0; y < height; y++ {
		for z := 0; z < width; z++ {
			for x := 0; x < length; x++ {
				// get cell index of current voxel
				xcell := int(cgm.Floor32(float32(x) / xstep))
				ycell := int(cgm.Floor32(float32(y) / ystep))
				zcell := int(cgm.Floor32(float32(z) / zstep))

				// get relative position within the cell
				xrel := cgm.Mod32(float32(x), xstep) / xstep
				yrel := cgm.Mod32(float32(y), ystep) / ystep
				zrel := cgm.Mod32(float32(z), zstep) / zstep
				pixelposition := mgl32.Vec3{xrel, yrel, zrel}

				// define weight function
				weight := func(x, y, z int) float32 {
					// get gradient vector and its position
					xabs := loop(xcell+x, res)
					yabs := loop(ycell+y, res)
					zabs := loop(zcell+z, res)
					gradientvector := vectors[yabs][zabs][xabs]

					// calc distance from this point to the pixel position
					gradientposition := mgl32.Vec3{float32(x), float32(y), float32(z)}
					distancevector := pixelposition.Sub(gradientposition)

					// return dot product
					return distancevector.Dot(gradientvector)
				}

				// get all 8 random vectors on the grid points surrounding the current cell
				w000 := weight(0, 0, 0)
				w001 := weight(0, 0, 1)
				w100 := weight(1, 0, 0)
				w101 := weight(1, 0, 1)
				w010 := weight(0, 1, 0)
				w011 := weight(0, 1, 1)
				w110 := weight(1, 1, 0)
				w111 := weight(1, 1, 1)

				// interpolate values
				tx := fade(pixelposition.X())
				ty := fade(pixelposition.Y())
				tz := fade(pixelposition.Z())
				wz0 := lerp(lerp(w000, w100, tx), lerp(w001, w101, tx), tz)
				wz1 := lerp(lerp(w010, w110, tx), lerp(w011, w111, tx), tz)
				val := lerp(wz0, wz1, ty)

				// bookkeeping of the biggest smallest distance for the following
				// normalization step
				maxval = cgm.Max32(maxval, val)

				// append val to result data
				values = append(values, val)
			}
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

// loop loops the value val between 0 and res-1
func loop(val, res int) int {
	newval := val % res
	if newval < 0 {
		newval = res + newval
	}
	return newval
}

func lerp(a, b, t float32) float32 {
	return (1-t)*a + t*b
}

func fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

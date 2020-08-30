package noise

import (
	"math"
	"math/rand"
	"time"

	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

// Worley2D creates 3D worley noise of the size specified by length x width x height
// with the specified resolution.
// It returns a 1D slice of uint8 values between 0 and 255.
func Worley2D(width, height, res int, radius float32) []uint8 {
	data := []uint8{}

	// divide volume into cells
	xstep := float32(width) / float32(res)
	ystep := float32(height) / float32(res)

	// position randomly exactly one point per cell
	points := make([][]mgl32.Vec2, res)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for y := 0; y < res; y++ {
		points[y] = make([]mgl32.Vec2, res)
		for x := 0; x < res; x++ {
			xr := cgm.Map(r.Float32(), 0, 1, float32(x)*xstep, float32(x+1)*xstep)
			yr := cgm.Map(r.Float32(), 0, 1, float32(y)*ystep, float32(y+1)*ystep)
			points[y][x] = mgl32.Vec2{xr, yr}
		}
	}

	// for each voxel find shortest distance to point in 27-neighborhood
	// loop at the edges to have tileable noise
	//var maxdist float32 = 0
	voxels := make([][]float32, height)
	for y := 0; y < height; y++ {
		voxels[y] = make([]float32, width)
		for x := 0; x < width; x++ {
			// center of current voxel
			voxel := mgl32.Vec2{float32(x) + 0.5, float32(y) + 0.5}

			// get cell index of current voxel
			xcell := int(cgm.Floor32(float32(x) / xstep))
			ycell := int(cgm.Floor32(float32(y) / ystep))

			// calc distance to each point in 27-neighborhood
			var mindist float32 = math.MaxFloat32
			for yd := -1; yd <= 1; yd++ {
				for xd := -1; xd <= 1; xd++ {
					// get position of point in current neighborhood cell
					// make sure to loop at the edges
					xabs := loop(xcell+xd, res)
					yabs := loop(ycell+yd, res)
					point := points[yabs][xabs]

					// offset the point if its on the other side of
					// the volume, this has to be done else the
					// distance doesn't loop
					var (
						xoff float32 = 0
						yoff float32 = 0
					)
					if xabs < xcell+xd {
						xoff = xstep * float32(res)
					}
					if xcell+xd < 0 {
						xoff = -xstep * float32(res)
					}
					if yabs < ycell+yd {
						yoff = ystep * float32(res)
					}
					if ycell+yd < 0 {
						yoff = -ystep * float32(res)
					}
					point = point.Add(mgl32.Vec2{xoff, yoff})

					// calc distance to this point
					dist := point.Sub(voxel).Len()

					// keep the shortest distance
					mindist = cgm.Min32(mindist, dist)
				}
			}
			// each voxel stores the shortest distance to a point in any of
			// the neighboring cells
			voxels[y][x] = mindist

			// bookkeeping of the biggest smallest distance for the following
			// normalization step
			//maxdist = cgm.Max32(maxdist, mindist)
		}
	}

	//maxdist *= 0.9
	maxdist := radius

	// map distance to 0..255 and save in data slice
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := voxels[y][x]
			val = cgm.Clamp(float32(val), 0, maxdist)
			val = cgm.Map(val, 0, maxdist, 255, 0)
			data = append(data, uint8(val))
		}
	}

	return data
}

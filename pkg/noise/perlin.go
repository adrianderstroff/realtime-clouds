package noise

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

// taken from http://flafla2.github.io/2014/08/09/perlinnoise.html
type perlin struct {
	repeat int
	p      []int
}

var permutation = []int{
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

var gradients3d = []mgl32.Vec3{
	mgl32.Vec3{1, 1, 0}, mgl32.Vec3{-1, 1, 0}, mgl32.Vec3{1, -1, 0}, mgl32.Vec3{-1, -1, 0},
	mgl32.Vec3{1, 0, 1}, mgl32.Vec3{-1, 0, 1}, mgl32.Vec3{1, 0, -1}, mgl32.Vec3{-1, 0, -1},
	mgl32.Vec3{0, 1, 1}, mgl32.Vec3{0, -1, 1}, mgl32.Vec3{0, 1, -1}, mgl32.Vec3{0, -1, -1},
}

func makeperlin(repeat int) perlin {
	var p [512]int
	for i := 0; i < 512; i++ {
		p[i] = permutation[i%256]
	}

	return perlin{
		repeat: repeat,
		p:      p[:],
	}
}

func (pln *perlin) octave(x, y, z float32, octaves int, persistence float32) float32 {
	var (
		total     float32 = 0
		frequency float32 = 1
		amplitude float32 = 1
		maxValue  float32 = 0
	)
	for i := 0; i < octaves; i++ {
		total += pln.simple(x*frequency, y*frequency, z*frequency) * amplitude

		maxValue += amplitude

		amplitude *= persistence
		frequency *= 2
	}

	return total / maxValue
}

func (pln *perlin) simple(x, y, z float32) float32 {
	if pln.repeat > 0 {
		x = cgm.Mod32(x, float32(pln.repeat))
		y = cgm.Mod32(y, float32(pln.repeat))
		z = cgm.Mod32(z, float32(pln.repeat))
	}

	xi := int(x) & 255
	yi := int(y) & 255
	zi := int(z) & 255
	xf := x - cgm.Floor32(x)
	yf := y - cgm.Floor32(y)
	zf := z - cgm.Floor32(z)
	u := fade(xf)
	v := fade(yf)
	w := fade(zf)

	aaa := pln.p[pln.p[pln.p[xi]+yi]+zi]
	aba := pln.p[pln.p[pln.p[xi]+pln.inc(yi)]+zi]
	aab := pln.p[pln.p[pln.p[xi]+yi]+pln.inc(zi)]
	abb := pln.p[pln.p[pln.p[xi]+pln.inc(yi)]+pln.inc(zi)]
	baa := pln.p[pln.p[pln.p[pln.inc(xi)]+yi]+zi]
	bba := pln.p[pln.p[pln.p[pln.inc(xi)]+pln.inc(yi)]+zi]
	bab := pln.p[pln.p[pln.p[pln.inc(xi)]+yi]+pln.inc(zi)]
	bbb := pln.p[pln.p[pln.p[pln.inc(xi)]+pln.inc(yi)]+pln.inc(zi)]

	x1 := cgm.Lerp(grad(aaa, xf, yf, zf), grad(baa, xf-1, yf, zf), u)
	x2 := cgm.Lerp(grad(aba, xf, yf-1, zf), grad(bba, xf-1, yf-1, zf), u)
	y1 := cgm.Lerp(x1, x2, v)

	x1 = cgm.Lerp(grad(aab, xf, yf, zf-1), grad(bab, xf-1, yf, zf-1), u)
	x2 = cgm.Lerp(grad(abb, xf, yf-1, zf-1), grad(bbb, xf-1, yf-1, zf-1), u)
	y2 := cgm.Lerp(x1, x2, v)
	return (cgm.Lerp(y1, y2, w) + 1) / 2
}

func (pln *perlin) inc(num int) int {
	num++
	if pln.repeat > 0 {
		num %= pln.repeat
	}
	return num
}

func fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

func grad(hash int, x, y, z float32) float32 {
	h := hash & 15
	u := x
	if h >= 8 {
		u = y
	}
	if h&1 != 0 {
		u = -u
	}

	v := z
	if h < 4 {
		v = y
	} else if h == 12 {
		v = x
	}
	if h&2 != 0 {
		v = -v
	}

	return u + v
}

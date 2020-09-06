package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/ssbo"
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// Perlin is a gpu perlin noise generator
type Perlin struct {
	computeshader shader.Shader
	permuations   ssbo.SSBO
	// dimensions
	width   int32
	height  int32
	octaves int32
	// other
	z          int32
	repeat     int32
	resolution int32
	// fbm
	scale       float32
	persistance float32
	// post processing
	brightness float32
	contrast   float32
}

// MakePerlin creates a perlin noise generator
func MakePerlin(shaderpath string) Perlin {
	computeshader, err := shader.MakeCompute(shaderpath + "/noise/perlin.comp")
	if err != nil {
		panic(err)
	}

	// create permuations buffer
	var permutation = []int32{
		151, 160, 137, 91, 90, 15, 131, 13, 201, 95,
		96, 53, 194, 233, 7, 225, 140, 36, 103, 30,
		69, 142, 8, 99, 37, 240, 21, 10, 23, 190,
		6, 148, 247, 120, 234, 75, 0, 26, 197, 62,
		94, 252, 219, 203, 117, 35, 11, 32, 57, 177, // 50
		33, 88, 237, 149, 56, 87, 174, 20, 125, 136,
		171, 168, 68, 175, 74, 165, 71, 134, 139, 48,
		27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
		60, 211, 133, 230, 220, 105, 92, 41, 55, 46,
		245, 40, 244, 102, 143, 54, 65, 25, 63, 161, // 100
		1, 216, 80, 73, 209, 76, 132, 187, 208, 89,
		18, 169, 200, 196, 135, 130, 116, 188, 159, 86,
		164, 100, 109, 198, 173, 186, 3, 64, 52, 217,
		226, 250, 124, 123, 5, 202, 38, 147, 118, 126,
		255, 82, 85, 212, 207, 206, 59, 227, 47, 16, // 150
		58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
		119, 248, 152, 2, 44, 154, 163, 70, 221, 153,
		101, 155, 167, 43, 172, 9, 129, 22, 39, 253,
		19, 98, 108, 110, 79, 113, 224, 232, 178, 185,
		112, 104, 218, 246, 97, 228, 251, 34, 242, 193, // 200
		238, 210, 144, 12, 191, 179, 162, 241, 81, 51,
		145, 235, 249, 14, 239, 107, 49, 192, 214, 31,
		181, 199, 106, 157, 184, 84, 204, 176, 115, 121,
		50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
		222, 114, 67, 29, 24, 72, 243, 141, 128, 195, // 250
		78, 66, 215, 61, 156, 180, // 256
	}

	var p = make([]int32, 512)
	permlen := len(permutation)
	for i := 0; i < 512; i++ {
		p[i] = permutation[i%permlen]
	}

	permutationsbuffer := ssbo.Make(ssbo.Int32, len(p))
	permutationsbuffer.UploadArrayI32(p)

	return Perlin{
		computeshader: computeshader,

		permuations: permutationsbuffer,

		width:  1024,
		height: 1024,

		z:          0,
		repeat:     2,
		resolution: 10,

		octaves:     1,
		scale:       1,
		persistance: 1,

		brightness: 1.0,
		contrast:   1.0,
	}
}

// UpdateState updates the worley noise parameters
func (p *Perlin) UpdateState(state *State) {
	p.z = state.Pz
	p.octaves = state.Poctaves
	p.resolution = int32(cgm.Pow32(2, float32(state.Presolution)))
	p.brightness = state.Pbrightness
	p.contrast = state.Pcontrast
	p.scale = state.Pscale
	p.persistance = state.Ppersistance
}

// GenerateTexture populates the texture with a worley noise
func (p *Perlin) GenerateTexture(tex *texture.Texture) {
	gl.BindImageTexture(0, tex.GetHandle(), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	p.permuations.Bind(1)

	p.computeshader.Use()
	p.computeshader.UpdateInt32("uWidth", p.width)
	p.computeshader.UpdateInt32("uHeight", p.height)
	p.computeshader.UpdateInt32("uOctaves", p.octaves)
	p.computeshader.UpdateInt32("uResolution", p.resolution)
	p.computeshader.UpdateFloat32("uBrightness", p.brightness)
	p.computeshader.UpdateFloat32("uContrast", p.contrast)
	p.computeshader.UpdateFloat32("uScale", p.scale)
	p.computeshader.UpdateFloat32("uPersistance", p.persistance)
	p.computeshader.UpdateInt32("uZ", p.z)
	p.computeshader.Compute(uint32(p.width), uint32(p.height), 1)
	p.computeshader.Release()
	gl.BindImageTexture(0, 0, 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	p.permuations.Unbind()

	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// Worley is a gpu worley noise generator
type Worley struct {
	computeshader shader.Shader
	noisetexture  texture.Texture
	width         int32
	height        int32
	resolution    int32
	octaves       int32
	radius        float32
	// fbm
	scale       float32
	persistance float32
	// post processing
	brightness float32
	contrast   float32
}

// MakeWorley creates a worley noise generator
func MakeWorley(shaderpath string) Worley {
	computeshader, err := shader.MakeCompute(shaderpath + "/noise/worley.comp")
	if err != nil {
		panic(err)
	}

	//create random seed
	randomdata := createRandom(1024 * 1024 * 4)
	noisetexture, err := texture.MakeFromData(randomdata, 1024, 1024, gl.RGBA32F, gl.RGBA)
	if err != nil {
		panic(err)
	}

	return Worley{
		computeshader: computeshader,
		noisetexture:  noisetexture,

		width:      1024,
		height:     1024,
		resolution: 32,
		octaves:    1,
		radius:     40.0,

		brightness: 1.0,
		contrast:   1.0,
	}
}

// UpdateState updates the worley noise parameters
func (w *Worley) UpdateState(state *State) {
	w.resolution = state.resolution
	w.octaves = state.octaves
	w.radius = state.radius
	w.brightness = state.wbrightness
	w.contrast = state.wcontrast
	w.scale = state.wscale
	w.persistance = state.wpersistance
}

// GenerateTexture populates the texture with a worley noise
func (w *Worley) GenerateTexture(tex *texture.Texture) {
	gl.BindImageTexture(0, tex.GetHandle(), 0, false, 0, gl.READ_WRITE, gl.RGBA32F)
	gl.BindImageTexture(1, w.noisetexture.GetHandle(), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)

	w.computeshader.Use()
	w.computeshader.UpdateInt32("uWidth", w.width)
	w.computeshader.UpdateInt32("uHeight", w.height)
	w.computeshader.UpdateInt32("uResolution", w.resolution)
	w.computeshader.UpdateInt32("uOctaves", w.octaves)
	w.computeshader.UpdateFloat32("uRadius", w.radius)
	w.computeshader.UpdateFloat32("uBrightness", w.brightness)
	w.computeshader.UpdateFloat32("uContrast", w.contrast)
	w.computeshader.UpdateFloat32("uScale", w.scale)
	w.computeshader.UpdateFloat32("uPersistance", w.persistance)
	w.computeshader.Compute(uint32(w.width), uint32(w.height), 1)
	w.computeshader.Compute(1024, 1024, 1)
	w.computeshader.Release()

	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)

	gl.BindImageTexture(0, 0, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, 0, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
}

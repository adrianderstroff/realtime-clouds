package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// operations
const (
	OP_LERP int = 0
)

// Merge is a gpu texture merge processor
type Merge struct {
	computeshader shader.Shader
	width         int32
	height        int32
	operation     int
}

// MakeMerge creates a merge processor
func MakeMerge(shaderpath string) Merge {
	computeshader, err := shader.MakeCompute(shaderpath + "/noise/merge.comp")
	if err != nil {
		panic(err)
	}

	return Merge{
		computeshader: computeshader,
		width:         1024,
		height:        1024,
		operation:     OP_LERP,
	}
}

// UpdateState updates the worley noise parameters
func (m *Merge) UpdateState(operation int32) {
	m.operation = int(operation)
}

// MergeTextures merges two textures with a specified operation and writes the results to the out texture
func (m *Merge) MergeTextures(tex1 *texture.Texture, tex2 *texture.Texture, texout *texture.Texture) {
	gl.BindImageTexture(0, tex1.GetHandle(), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, tex2.GetHandle(), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(2, texout.GetHandle(), 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	m.computeshader.Use()
	m.computeshader.UpdateInt32("uWidth", m.width)
	m.computeshader.UpdateInt32("uHeight", m.height)
	m.computeshader.UpdateInt32("uOperation", int32(m.operation))
	m.computeshader.Compute(uint32(m.width), uint32(m.height), 1)
	m.computeshader.Release()

	gl.BindImageTexture(0, 0, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, 0, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(2, 0, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

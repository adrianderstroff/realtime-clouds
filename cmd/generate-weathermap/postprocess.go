package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// PostProcess is a gpu post process processor
type PostProcess struct {
	computeshader shader.Shader
	width         int32
	height        int32
	threshold     float32
}

// MakePostProcess creates a post process processor
func MakePostProcess(shaderpath string) PostProcess {
	computeshader, err := shader.MakeCompute(shaderpath + "/noise/postprocess.comp")
	if err != nil {
		panic(err)
	}

	return PostProcess{
		computeshader: computeshader,
		width:         1024,
		height:        1024,
		threshold:     0.0,
	}
}

// UpdateState updates the worley noise parameters
func (p *PostProcess) UpdateState(state *State) {
	p.threshold = state.Threshold
}

// Apply applies post processing to a texture and stores it into another texture
func (p *PostProcess) Apply(texin *texture.Texture, texout *texture.Texture) {
	gl.BindImageTexture(0, texin.GetHandle(), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, texout.GetHandle(), 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	p.computeshader.Use()
	p.computeshader.UpdateInt32("uWidth", p.width)
	p.computeshader.UpdateInt32("uHeight", p.height)
	p.computeshader.UpdateFloat32("uThreshold", p.threshold)
	p.computeshader.Compute(uint32(p.width), uint32(p.height), 1)
	p.computeshader.Release()

	gl.BindImageTexture(0, 0, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)
	gl.BindImageTexture(1, 0, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// Clear clears a texture to a specific color
type Clear struct {
	computeshader shader.Shader
	width         int32
	height        int32
	clearcolor    mgl32.Vec3
}

// MakeClear creates a clear texture processor
func MakeClear(shaderpath string) Clear {
	computeshader, err := shader.MakeCompute(shaderpath + "/noise/clear.comp")
	if err != nil {
		panic(err)
	}

	return Clear{
		computeshader: computeshader,
		width:         1024,
		height:        1024,
		clearcolor:    mgl32.Vec3{1, 1, 1},
	}
}

// ClearTexture clears a texture to the specified color
func (c *Clear) ClearTexture(tex *texture.Texture) {
	gl.BindImageTexture(0, tex.GetHandle(), 0, false, 0, gl.READ_ONLY, gl.RGBA32F)

	c.computeshader.Use()
	c.computeshader.UpdateInt32("uWidth", c.width)
	c.computeshader.UpdateInt32("uHeight", c.height)
	c.computeshader.UpdateVec3("uClearColor", c.clearcolor)
	c.computeshader.Compute(uint32(c.width), uint32(c.height), 1)
	c.computeshader.Release()

	gl.BindImageTexture(0, 0, 0, false, 0, gl.READ_ONLY, gl.RGBA32F)

	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

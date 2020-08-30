package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/plane"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// Renderpass encapulates rendering textured planes
type Renderpass struct {
	shader shader.Shader
}

// MakeRenderpass creates a render pass
func MakeRenderpass(shaderpath string) Renderpass {
	// make shader
	plane := plane.Make(2, 2, gl.TRIANGLES)
	weathermapshader, err := shader.Make(shaderpath+"/texture/texture.vert", shaderpath+"/texture/texture.frag")
	if err != nil {
		panic(err)
	}
	weathermapshader.AddRenderable(plane)

	return Renderpass{
		shader: weathermapshader,
	}
}

// Render renders the textures planes to the screen
func (r *Renderpass) Render(camera *trackball.Trackball, tex *texture.Texture) {
	var off float32 = 2.00

	tex.Bind(0)

	r.shader.Use()
	r.shader.UpdateMat4("M", mgl32.Ident4())
	r.shader.UpdateMat4("V", camera.GetView())
	r.shader.UpdateMat4("P", camera.GetPerspective())
	r.shader.Render()

	r.shader.UpdateMat4("M", mgl32.Translate3D(off, 0, 0))
	r.shader.Render()

	r.shader.UpdateMat4("M", mgl32.Translate3D(-off, 0, 0))
	r.shader.Render()

	r.shader.UpdateMat4("M", mgl32.Translate3D(0, 0, off))
	r.shader.Render()

	r.shader.UpdateMat4("M", mgl32.Translate3D(0, 0, -off))
	r.shader.Render()
	r.shader.Release()

	tex.Unbind()
}

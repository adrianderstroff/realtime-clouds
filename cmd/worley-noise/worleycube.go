package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

type WorleyCube struct {
	shader      shader.Shader
	mesh        mesh.Mesh
	modelmatrix mgl32.Mat4
}

func MakeWorleyCube(width, height, points int) (WorleyCube, error) {
	// make worley texture
	data := noise.MakeWorley(width, height, points)
	tex, err := texture.MakeFromData(int32(width), int32(height), gl.RGB, data)
	if err != nil {
		return WorleyCube{}, err
	}

	// make box and apply texture
	mesh := box.Make(2, 2, 2, false, gl.TRIANGLES)
	mesh.AddTexture(tex)

	// make shader and add mesh
	shader, _ := shader.Make(SHADER_PATH+"/texture/texture.vert", SHADER_PATH+"/texture/texture.frag")
	shader.AddRenderable(mesh)

	// initial model matrix
	m := mgl32.Ident4()

	// construct worley object
	return WorleyCube{
		shader:      shader,
		mesh:        mesh,
		modelmatrix: m,
	}, nil
}

func (w *WorleyCube) SetModelMatrix(m mgl32.Mat4) {
	w.modelmatrix = m
}

func (w *WorleyCube) Render(camera camera.Camera) {
	w.shader.Use()
	w.shader.UpdateMat4("M", w.modelmatrix)
	w.shader.UpdateMat4("V", camera.GetView())
	w.shader.UpdateMat4("P", camera.GetPerspective())
	w.shader.Render()
}

package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
	"github.com/go-gl/mathgl/mgl32"
)

type LandscapePass struct {
	landscapeshader shader.Shader
}

func MakeLandscapePass(shaderpath string) LandscapePass {
	// create shaders
	//plane := plane.Make(100, 100, gl.TRIANGLES)
	box := box.Make(4000, 1, 4000, false, gl.TRIANGLES)
	landscapeshader, err := shader.Make(shaderpath+"/flat/flat.vert", shaderpath+"/flat/flat.frag")
	if err != nil {
		panic(err)
	}
	landscapeshader.AddRenderable(box)

	return LandscapePass{
		landscapeshader: landscapeshader,
	}
}

func (lsp *LandscapePass) Render(camera camera.Camera) {
	lsp.landscapeshader.Use()
	lsp.landscapeshader.UpdateMat4("M", mgl32.Ident4())
	lsp.landscapeshader.UpdateMat4("V", camera.GetView())
	lsp.landscapeshader.UpdateMat4("P", camera.GetPerspective())
	lsp.landscapeshader.UpdateVec3("flatColor", mgl32.Vec3{0, 0.8, 0.2})
	lsp.landscapeshader.Render()
	lsp.landscapeshader.Release()
}

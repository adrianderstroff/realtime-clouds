package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/plane"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

type RaymarchingPass struct {
	cloudbasefbo   texture.Texture
	clouddetailfbo texture.Texture
	turbulencefbo  texture.Texture
	cloudmapfbo    texture.Texture
	raymarchshader shader.Shader
}

func MakeRaymarchingPass(width, height int, texpath, shaderpath string) RaymarchingPass {
	// create fbos
	cloudbasefbo, err := texture.Make3DFromPath(MakePathsFromDirectory(texpath+"cloud-base/", "base", "png", 0, 127), gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	clouddetailfbo, err := texture.Make3DFromPath(MakePathsFromDirectory(texpath+"cloud-detail/", "detail", "png", 0, 31), gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	turbulencefbo, err := texture.MakeFromPath(texpath+"cloud-turbulence/turbulence.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	//cloudmapfbo, err := texture.MakeFromPath(texpath+"cloud-map/cloud-map.png", gl.RGBA, gl.RGBA)
	cloudmapfbo, err := texture.MakeFromPath(texpath+"debug.jpg", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}

	// create shaders
	plane := plane.Make(2, 2, gl.TRIANGLES)
	raymarchshader, err := shader.Make(shaderpath+"/cloud/raymarch.vert", shaderpath+"/cloud/raymarch.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(plane)

	return RaymarchingPass{
		cloudbasefbo:   cloudbasefbo,
		clouddetailfbo: clouddetailfbo,
		turbulencefbo:  turbulencefbo,
		cloudmapfbo:    cloudmapfbo,
		raymarchshader: raymarchshader,
	}
}

func (rmp *RaymarchingPass) Render(camera camera.Camera) {
	rmp.cloudbasefbo.Bind(0)
	rmp.clouddetailfbo.Bind(1)
	rmp.turbulencefbo.Bind(2)
	rmp.cloudmapfbo.Bind(3)

	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateVec3("cameraPos", camera.GetPos())
	rmp.raymarchshader.UpdateFloat32("width", 800)
	rmp.raymarchshader.UpdateFloat32("height", 600)
	rmp.raymarchshader.UpdateMat4("M", mgl32.Ident4())
	rmp.raymarchshader.UpdateMat4("V", camera.GetView())
	rmp.raymarchshader.UpdateMat4("P", camera.GetPerspective())
	rmp.raymarchshader.Render()
	rmp.raymarchshader.Release()

	rmp.cloudbasefbo.Unbind()
	rmp.clouddetailfbo.Unbind()
	rmp.turbulencefbo.Unbind()
	rmp.cloudmapfbo.Unbind()
}

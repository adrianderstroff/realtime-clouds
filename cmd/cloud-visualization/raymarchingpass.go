package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/fbo"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

type RaymarchingPass struct {
	raystartfbo    fbo.FBO
	rayendfbo      fbo.FBO
	setupshader    shader.Shader
	raymarchshader shader.Shader
}

func MakeRaymarchingPass(width, height int, shaderpath string) RaymarchingPass {
	// create fbos
	raystartfbo := fbo.Make(width, height)
	raystartcolor := texture.MakeColor(width, height)
	raystartdepth := texture.MakeDepth(width, height)
	raystartfbo.AttachColorTexture(&raystartcolor, 0)
	raystartfbo.AttachDepthTexture(&raystartdepth)
	rayendfbo := fbo.Make(width, height)
	rayendcolor := texture.MakeColor(width, height)
	rayenddepth := texture.MakeDepth(width, height)
	rayendfbo.AttachColorTexture(&rayendcolor, 0)
	rayendfbo.AttachDepthTexture(&rayenddepth)

	// create shaders
	cube := box.Make(2, 2, 2, false, gl.TRIANGLES)
	setupshader, _ := shader.Make(shaderpath+"/setup/setup.vert", shaderpath+"/setup/setup.frag")
	setupshader.AddRenderable(cube)
	raymarchshader, err := shader.Make(shaderpath+"/raymarch/raymarch.vert", shaderpath+"/raymarch/cloud.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(cube)

	return RaymarchingPass{
		raystartfbo:    raystartfbo,
		rayendfbo:      rayendfbo,
		setupshader:    setupshader,
		raymarchshader: raymarchshader,
	}
}

func (rmp *RaymarchingPass) Render(fbo *fbo.FBO, tex3d *texture.Texture, M, V, P mgl32.Mat4, iterations int32) {
	// calculate ray start and end positions
	rmp.raystartfbo.Bind()
	rmp.raystartfbo.Clear()
	rmp.setupshader.Use()
	rmp.setupshader.UpdateMat4("M", M)
	rmp.setupshader.UpdateMat4("V", V)
	rmp.setupshader.UpdateMat4("P", P)
	rmp.setupshader.Render()
	rmp.raystartfbo.Unbind()

	gl.CullFace(gl.FRONT)
	rmp.rayendfbo.Bind()
	rmp.rayendfbo.Clear()
	rmp.setupshader.Use()
	rmp.setupshader.UpdateMat4("M", M)
	rmp.setupshader.UpdateMat4("V", V)
	rmp.setupshader.UpdateMat4("P", P)
	rmp.setupshader.Render()
	rmp.rayendfbo.Unbind()
	gl.CullFace(gl.BACK)

	// render box
	fbo.Bind()
	fbo.Clear()
	rmp.raystartfbo.GetColorTexture(0).Bind(0)
	rmp.rayendfbo.GetColorTexture(0).Bind(1)
	tex3d.Bind(2)
	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateInt32("iterations", iterations)
	rmp.raymarchshader.Render()
	rmp.raystartfbo.GetColorTexture(0).Unbind()
	rmp.rayendfbo.GetColorTexture(0).Unbind()
	tex3d.Unbind()
	fbo.Unbind()
}

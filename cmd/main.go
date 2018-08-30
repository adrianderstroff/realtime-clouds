package main

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"
	shader "github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	window "github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	trackball "github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
	box "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
	skybox "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/skybox"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/"
	CUBEMAP_PATH = "./assets/images/cubemap/"
)

var (
	width  int = 800
	height int = 600
)

func main() {
	runtime.LockOSThread()

	// setup window
	title := "Real-time clouds"
	window, _ := window.New(title, int(width), int(height))
	interaction := interaction.New(window)
	defer window.Close()

	// make textured cube
	texshader, _ := shader.Make(SHADER_PATH+"/texture/texture.vert", SHADER_PATH+"/texture/texture.frag")
	mesh := box.Make(1, 1, 1, false, gl.TRIANGLES)
	texture, _ := tex.MakeFromPath(TEX_PATH + "/profile.png")
	texture.GenMipmap()
	mesh.AddTexture(texture)
	texshader.AddRenderable(mesh)

	// make skybox
	skyboxshader, _ := shader.Make(SHADER_PATH+"/skybox/skybox.vert", SHADER_PATH+"/skybox/skybox.frag")
	sky, err := skybox.MakeFromDirectory(50.0, CUBEMAP_PATH+"/debug/", "png", gl.TRIANGLES)
	if err != nil {
		panic(err)
	}
	skyboxshader.AddRenderable(sky)

	// make camera
	camera := trackball.MakeDefault(width, height, 10.0)
	interaction.AddInteractable(&camera)

	// main loop
	render := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// get camera matrices
		M := mgl32.Ident4()
		V := camera.GetView()
		P := camera.GetPerspective()

		// render box
		texshader.Use()
		texshader.UpdateMat4("M", M)
		texshader.UpdateMat4("V", V)
		texshader.UpdateMat4("P", P)
		texshader.Render()

		// render skybox
		skyboxshader.Use()
		skyboxshader.UpdateMat4("M", M)
		skyboxshader.UpdateMat4("V", V)
		skyboxshader.UpdateMat4("P", P)
		skyboxshader.Render()
	}
	window.RunMainLoop(render)
}

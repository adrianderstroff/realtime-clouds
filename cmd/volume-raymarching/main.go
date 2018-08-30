package main

import (
	"runtime"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/"
	CUBEMAP_PATH = "./assets/images/cubemap/"

	WIDTH  int = 800
	HEIGHT int = 600
)

func main() {
	runtime.LockOSThread()

	// setup window
	title := "Volume raymarching"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 10.0)
	interaction.AddInteractable(&camera)

	// generate worley noise
	worleyw, worleyh := 128, 128
	worleydata := noise.MakeWorley(worleyw, worleyh, 20)
	worleytex, _ := texture.MakeFromData(int32(worleyw), int32(worleyh), gl.RGB, worleydata)
	worleytex.GenMipmap()

	// make textured cube
	texshader, _ := shader.Make(SHADER_PATH+"/texture/texture.vert", SHADER_PATH+"/texture/texture.frag")
	mesh := box.Make(2, 2, 2, false, gl.TRIANGLES)
	mesh.AddTexture(worleytex)
	texshader.AddRenderable(mesh)

	// render loop
	renderloop := func() {
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
	}
	window.RunMainLoop(renderloop)
}

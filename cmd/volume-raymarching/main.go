package main

import (
	"runtime"

	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/fbo"
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
	window.LockFPS(60)
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 10.0)
	interaction.AddInteractable(&camera)

	// generate 3D texture with worley noise
	noisedata := noise.Worley3D(128, 128, 128, 20)
	noisetex, err := texture.Make3DFromData(noisedata, 128, 128, 128, gl.RED, gl.RED)
	if err != nil {
		panic(err)
	}

	// create fbos
	raystartfbo := fbo.Make(WIDTH, HEIGHT)
	raystartcolor := texture.MakeColor(WIDTH, HEIGHT)
	raystartdepth := texture.MakeDepth(WIDTH, HEIGHT)
	raystartfbo.AttachColorTexture(&raystartcolor, 0)
	raystartfbo.AttachDepthTexture(&raystartdepth)

	rayendfbo := fbo.Make(WIDTH, HEIGHT)
	rayendcolor := texture.MakeColor(WIDTH, HEIGHT)
	rayenddepth := texture.MakeDepth(WIDTH, HEIGHT)
	rayendfbo.AttachColorTexture(&rayendcolor, 0)
	rayendfbo.AttachDepthTexture(&rayenddepth)

	// create cube
	cube := box.Make(2, 2, 2, false, gl.TRIANGLES)

	// prepare setup shader
	setupshader, _ := shader.Make(SHADER_PATH+"/setup/setup.vert", SHADER_PATH+"/setup/setup.frag")
	setupshader.AddRenderable(cube)

	// prepare raymarching shader
	raymarchshader, err := shader.Make(SHADER_PATH+"/raymarch/raymarch.vert", SHADER_PATH+"/raymarch/raymarch.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(cube)

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()
		M := mgl32.Ident4()
		V := camera.GetView()
		P := camera.GetPerspective()

		// calculate ray start and end positions
		raystartfbo.Bind()
		raystartfbo.Clear()
		setupshader.Use()
		setupshader.UpdateMat4("M", M)
		setupshader.UpdateMat4("V", V)
		setupshader.UpdateMat4("P", P)
		setupshader.Render()
		raystartfbo.Unbind()

		gl.CullFace(gl.FRONT)
		rayendfbo.Bind()
		rayendfbo.Clear()
		setupshader.Use()
		setupshader.UpdateMat4("M", M)
		setupshader.UpdateMat4("V", V)
		setupshader.UpdateMat4("P", P)
		setupshader.Render()
		rayendfbo.Unbind()
		gl.CullFace(gl.BACK)

		// render box
		raystartcolor.Bind(0)
		rayendcolor.Bind(1)
		noisetex.Bind(2)
		raymarchshader.Use()
		raymarchshader.Render()
		raystartcolor.Unbind()
		rayendcolor.Unbind()
		noisetex.Unbind()

		// copy textures to screen
		raystartfbo.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), 0, 0, 200, 200)
		rayendfbo.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), 200, 0, 200, 200)
	}
	window.RunMainLoop(renderloop)
}

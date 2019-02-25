package main

import (
	"runtime"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/fps"
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
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "Realtime Clouds"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := fps.MakeDefault(WIDTH, HEIGHT, mgl32.Vec3{1000, 0, 0})
	interaction.AddInteractable(&camera)

	// make passes
	raymarchingpass := MakeRaymarchingPass(WIDTH, HEIGHT, TEX_PATH, SHADER_PATH)

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// do raymarching passes
		raymarchingpass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}

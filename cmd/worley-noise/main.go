package main

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
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
	title := "Worley noise"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 10.0)
	interaction.AddInteractable(&camera)

	// generate worley cubes
	w1, _ := MakeWorleyCube(128, 128, 5)
	w2, _ := MakeWorleyCube(128, 128, 50)
	w3, _ := MakeWorleyCube(128, 128, 500)
	w1.SetModelMatrix(mgl32.Translate3D(-2, 0, 0))
	w3.SetModelMatrix(mgl32.Translate3D(2, 0, 0))

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// render worley cubes
		w1.Render(&camera)
		w2.Render(&camera)
		w3.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}

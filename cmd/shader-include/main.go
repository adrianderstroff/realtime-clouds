package main

import (
	"runtime"

	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/box"
)

const (
	SHADER_PATH = "./assets/shaders/"

	WIDTH  int = 800
	HEIGHT int = 600
)

func main() {
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "Include Test"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	defer window.Close()
	inter := interaction.New(window)

	// create camera
	cam := trackball.MakeDefault(WIDTH, HEIGHT, 10)
	inter.AddInteractable(&cam)

	// setup shader
	cube := box.Make(1, 1, 1, false, gl.TRIANGLES)
	shader, err := shader.Make(SHADER_PATH+"/shared/main.vert", SHADER_PATH+"/shared/main.frag")
	if err != nil {
		panic(err)
	}
	shader.AddRenderable(cube)

	// update uniforms
	M := mgl32.Ident4()

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		cam.Update()

		// draw
		shader.Use()
		shader.UpdateMat4("M", M)
		shader.UpdateMat4("V", cam.GetView())
		shader.UpdateMat4("P", cam.GetPerspective())
		shader.UpdateVec3("cameraPos", cam.GetPos())
		shader.Render()
		shader.Release()
	}
	window.RunMainLoop(renderloop)
}

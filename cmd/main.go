package main

import (
	"runtime"
	"strconv"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-clouds/pkg/core"
	"github.com/adrianderstroff/realtime-clouds/pkg/mesh"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
)

const (
	SHADER_PATH = "./assets/shaders/"
	TEX_PATH    = "./assets/images/textures/"
	SKY_PATH    = "./assets/images/skyboxes/"
)

var (
	width  int = 800
	height int = 600
)

func main() {
	runtime.LockOSThread()

	// setup opengl
	title := "Real-time clouds"
	windowManager, err := core.NewWindowManager(title, int(width), int(height))
	if err != nil {
		panic(err)
	}
	defer windowManager.Close()

	// make shader
	shader, err := core.MakeProgram(SHADER_PATH+"/flat/flat.vert", SHADER_PATH+"/flat/flat.frag")
	if err != nil {
		panic(err)
	}

	// make mesh
	mesh := mesh.MakeQuad(2, 2, 2, false, gl.TRIANGLES)
	shader.AddRenderable(mesh)

	// make camera
	camera := camera.MakeDefaultTrackballCamera(width, height, 10.0)
	windowManager.AddInteractable(&camera)

	// main loop
	render := func() {
		// update title
		windowManager.SetTitle(title + " " + strconv.FormatFloat(windowManager.GetFPS(), 'f', 0, 64) + "FPS")

		// update camera
		camera.Update()

		// get camera matrices
		M := mgl32.Ident4()
		V := camera.GetView()
		P := camera.GetPerspective()

		// render
		shader.Use()
		shader.UpdateMat4("M", M)
		shader.UpdateMat4("V", V)
		shader.UpdateMat4("P", P)
		shader.UpdateVec3("flatColor", mgl32.Vec3{0, 0, 1})
		shader.Render()
	}
	windowManager.RunMainLoop(render)
}

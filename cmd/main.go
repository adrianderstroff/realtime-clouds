package main

import (
	"runtime"
	"strconv"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-clouds/pkg/core"
	"github.com/adrianderstroff/realtime-clouds/pkg/mesh"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/texture"
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
	windowManager, _ := core.NewWindowManager(title, int(width), int(height))
	defer windowManager.Close()

	// make shader
	shader, _ := core.MakeProgram(SHADER_PATH+"/texture/texture.vert", SHADER_PATH+"/texture/texture.frag")

	// make mesh
	mesh := mesh.MakeQuad(2, 2, 2, false, gl.TRIANGLES)
	texture, _ := tex.MakeTextureFromPath(TEX_PATH + "/profile.png")
	texture.GenMipmap()
	mesh.AddTexture(texture)
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
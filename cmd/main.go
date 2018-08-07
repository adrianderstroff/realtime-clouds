package main

import (
	"runtime"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
	"github.com/adrianderstroff/realtime-grass/pkg/scene"
)

const (
	SHADER_PATH = "./assets/shaders/"
	TEX_PATH    = "./assets/images/textures/"
	SKY_PATH    = "./assets/images/skyboxes/"
)

var (
	width         int32   = 800
	height        int32   = 600
	terrainheight float32 = 300.0
	viewdist      float32 = 5000.0
	windradius    int32   = 30
	windinfluence float32 = 4.0
	bladecount    int     = 100
	grassHeight   float32 = 50.0
)

func main() {
	runtime.LockOSThread()

	// setup opengl
	windowManager, err := engine.NewWindowManager("Grass", int(width), int(height))
	if err != nil {
		panic(err)
	}
	windowManager.LockFPS(30)
	windowManager.SetClearColor(0, 0, 0)
	defer windowManager.Close()

	// make terrain
	terrain, err := scene.MakeTerrain(SHADER_PATH, TEX_PATH, 5000.0, 10, 10, terrainheight, bladecount, grassHeight, viewdist, windradius, windinfluence)
	if err != nil {
		panic(err)
	}

	// make skybox
	sky, err := scene.MakeSky(SHADER_PATH, SKY_PATH)
	if err != nil {
		panic(err)
	}

	// set camera
	camera := engine.MakeCameraFPS(int(width), int(height), mgl32.Vec3{0.0, 100.0, 0.0}, 6.0, 45.0, 0.1, viewdist)
	windowManager.AddInteractable(&camera)
	oldpos := camera.Pos

	// fbo
	fbo := engine.MakeFBO(width, height)
	if !fbo.IsComplete() {
		panic("Fbo not complete")
	}

	// postprocessing
	pp, err := scene.MakePostprocessing(SHADER_PATH, width, height)
	if err != nil {
		panic(err)
	}

	// main loop
	render := func() {
		// update title
		windowManager.SetTitle("Grass " + strconv.FormatFloat(windowManager.GetFPS(), 'f', 0, 64) + "FPS")

		// update camera
		camera.Update()
		// collision check with terrain
		y := terrain.GetHeight(camera.Pos) + grassHeight
		if camera.Pos.Y() < y {
			newy := mathutils.Interpolate(camera.Pos.Y(), y, 0.5)
			camera.SetPos(mgl32.Vec3{camera.Pos.X(), newy, camera.Pos.Z()})
		}

		// get camera matrices
		M := mgl32.Ident4()
		V := camera.GetView()
		P := camera.GetPerspective()

		// render everything into an fbo
		fbo.Bind()
		fbo.Clear()

		// render terrain
		mvp := P.Mul4(V)
		cameradelta := camera.Pos.Sub(oldpos)
		cameradelta = mgl32.Vec3{cameradelta.X(), 0.0, cameradelta.Z()}
		terrain.Update(camera.Pos, cameradelta, mvp)
		terrain.Render(M, V, P, camera.Pos)

		// done rendering into fbo
		fbo.Unbind()

		// apply dof and bloom
		pp.Bloom(&fbo)
		pp.DOF(&fbo)

		// render skybox
		fbo.Bind()
		Vc := mathutils.ExtractRotation(&V)
		sky.Render(&Vc, &P)
		fbo.Unbind()

		// add fog
		pp.Fog(&fbo, &camera)

		// render fbo to screen
		fbo.CopyToScreen(0, 0, 0, width, height)

		// update old camera pos
		oldpos = camera.Pos
	}
	windowManager.RunMainLoop(render)
}

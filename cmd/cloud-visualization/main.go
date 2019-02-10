package main

import (
	"math"
	"runtime"
	"time"

	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/fbo"
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
	title := "Cloud Visualization"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 10.0)
	interaction.AddInteractable(&camera)

	// setup fbos
	fbo1 := fbo.Make(WIDTH, HEIGHT)

	// generate cloud base texture
	cloudbasetex, err := Make3DCloudTexture(TEX_PATH+"cloud-base/base", 128)
	if err != nil {
		panic(err)
	}

	// generate 3D texture with worley noise
	clouddetailtex, err := Make3DCloudTexture(TEX_PATH+"cloud-detail/detail", 32)
	if err != nil {
		panic(err)
	}
	_ = clouddetailtex

	// setup raymarching pass
	raymarchingpass := MakeRaymarchingPass(WIDTH, HEIGHT, SHADER_PATH)

	start := time.Now()
	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// get delta time
		delta := float32(time.Now().Sub(start).Seconds())

		// update camera
		camera.Update()
		M := mgl32.HomogRotate3DY(delta * math.Pi * 0.25)
		V := camera.GetView()
		P := camera.GetPerspective()

		// do raymarching passes
		//raymarchingpass.Render(&fbo1, &clouddetailtex, M, V, P, 10)
		raymarchingpass.Render(&fbo1, &cloudbasetex, M, V, P, 10)

		// copy textures to screen
		fbo1.CopyToScreen(0, 0, 0, int32(WIDTH), int32(HEIGHT))
	}
	window.RunMainLoop(renderloop)
}

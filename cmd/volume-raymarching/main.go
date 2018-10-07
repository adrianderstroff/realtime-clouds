package main

import (
	"math"
	"runtime"
	"time"

	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/fbo"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/interaction"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
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

	// setup fbos
	fbo1 := fbo.Make(WIDTH, HEIGHT)
	fbo2 := fbo.Make(WIDTH, HEIGHT)
	fbo3 := fbo.Make(WIDTH, HEIGHT)
	fbo4 := fbo.Make(WIDTH, HEIGHT)

	// generate 3D texture with worley noise
	worleydata := noise.Worley3D(128, 128, 128, 5)
	worleytex, err := texture.Make3DFromData(worleydata, 128, 128, 128, gl.RED, gl.RED)
	if err != nil {
		panic(err)
	}
	worleydata = noise.Worley3D(128, 128, 128, 10)
	worleytex2, err := texture.Make3DFromData(worleydata, 128, 128, 128, gl.RED, gl.RED)
	if err != nil {
		panic(err)
	}

	// generate 3D texture with perlin noise
	perlindata := noise.Perlin3D(128, 128, 128, 5)
	perlintex, err := texture.Make3DFromData(perlindata, 128, 128, 128, gl.RED, gl.RED)
	if err != nil {
		panic(err)
	}
	perlindata = noise.Perlin3D(128, 128, 128, 10)
	perlintex2, err := texture.Make3DFromData(perlindata, 128, 128, 128, gl.RED, gl.RED)
	if err != nil {
		panic(err)
	}

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
		M := mgl32.HomogRotate3DY(delta * math.Pi * 0.5)
		V := camera.GetView()
		P := camera.GetPerspective()

		// do raymarching passes
		raymarchingpass.Render(&fbo1, &worleytex, M, V, P, 10)
		raymarchingpass.Render(&fbo2, &perlintex, M, V, P, 10)
		raymarchingpass.Render(&fbo3, &worleytex2, M, V, P, 10)
		raymarchingpass.Render(&fbo4, &perlintex2, M, V, P, 10)

		// copy textures to screen
		fbo1.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), 0, 0, int32(WIDTH/2), int32(HEIGHT/2))
		fbo2.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), int32(WIDTH/2), 0, int32(WIDTH/2), int32(HEIGHT/2))
		fbo3.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), 0, int32(HEIGHT/2), int32(WIDTH/2), int32(HEIGHT/2))
		fbo4.CopyToScreenRegion(0, 0, 0, int32(WIDTH), int32(HEIGHT), int32(WIDTH/2), int32(HEIGHT/2), int32(WIDTH/2), int32(HEIGHT/2))
	}
	window.RunMainLoop(renderloop)
}

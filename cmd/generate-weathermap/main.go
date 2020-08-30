package main

import (
	"runtime"

	"github.com/adrianderstroff/realtime-clouds/pkg/gui"

	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera/trackball"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/window"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/"
	CUBEMAP_PATH = "./assets/images/cubemap/"
	OUT_PATH     = "./"

	WIDTH  int = 1200
	HEIGHT int = 1000
)

type State struct {
	// perlin
	useperlin    bool
	octaves      int32
	pbrightness  float32
	pcontrast    float32
	prepeat      int32
	pz           int32
	pscale       float32
	ppersistance float32
	// worley
	useworley    bool
	resolution   int32
	radius       float32
	wbrightness  float32
	wcontrast    float32
	wscale       float32
	wpersistance float32
	// post processing
	threshold int32
	// state
	dirty bool
}

func main() {
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "Generate Weather Map"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	defer window.Close()

	// set initial gl settings
	gl.ClearColor(1, 1, 1, 1)

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 8)
	camera.Rotate(90, 0)

	// create state
	state := State{
		// perlin
		useperlin:    true,
		octaves:      1,
		threshold:    120,
		pbrightness:  0,
		pcontrast:    0.5,
		prepeat:      2,
		pz:           0,
		pscale:       1,
		ppersistance: 1,
		// worley
		useworley:    true,
		resolution:   16,
		radius:       70,
		wbrightness:  0,
		wcontrast:    0.5,
		wscale:       1,
		wpersistance: 1,
		// state
		dirty: false,
	}

	// make texture
	weathermaptexture, err := texture.MakeFromData(createAndFillImage(1024*1024*4, 120), 1024, 1024, gl.RGBA32F, gl.RGBA)
	if err != nil {
		panic(err)
	}

	// make render pass
	renderpass := MakeRenderpass(SHADER_PATH)

	// create worley noise generator
	perlin := MakePerlin(SHADER_PATH)
	//worley := MakeWorley(SHADER_PATH)

	// create gui
	gamegui := gui.New(window.Window)

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// generate perlin
		perlin.UpdateState(&state)
		perlin.GenerateTexture(&weathermaptexture)
		err := gl.GetError()
		if err != nil {
			panic(err)
		}

		// generate worley
		//worley.UpdateState(&state)
		//worley.GenerateTexture(&weathermaptexture)

		renderpass.Render(&camera, &weathermaptexture)

		// gui
		gamegui.Begin()
		if gamegui.BeginWindow("Options", 0, 0, 250, float32(HEIGHT)) {
			if gamegui.BeginGroup("Perlin", 380) {
				gamegui.Checkbox("Use Perlin", &state.useperlin)
				gamegui.SliderInt32("Resolution", &state.resolution, 1, 64, 1)
				gamegui.SliderInt32("Repeat", &state.prepeat, 1, 10, 1)
				gamegui.SliderInt32("Z", &state.pz, 0, 1024, 1)
				gamegui.Label("Fbm")
				gamegui.SliderInt32("Octaves", &state.octaves, 1, 5, 1)
				gamegui.SliderFloat32("Scale", &state.pscale, 1, 3, 0.01)
				gamegui.SliderFloat32("Persistance", &state.ppersistance, 0, 1, 0.01)
				gamegui.Label("Post")
				gamegui.SliderFloat32("Brightness", &state.pbrightness, 0, 1, 0.01)
				gamegui.SliderFloat32("Contrast", &state.pcontrast, 0, 1, 0.01)
				gamegui.EndGroup()
			}
			if gamegui.BeginGroup("Worley", 320) {
				gamegui.Checkbox("Use Worley", &state.useworley)
				gamegui.SliderInt32("Resolution", &state.resolution, 1, 64, 1)
				gamegui.SliderFloat32("Radius", &state.radius, 0, 200, 0.1)
				gamegui.Label("Fbm")
				gamegui.SliderFloat32("Scale", &state.wscale, 1, 3, 0.01)
				gamegui.SliderFloat32("Persistance", &state.wpersistance, 0, 1, 0.01)
				gamegui.Label("Post")
				gamegui.SliderFloat32("Brightness", &state.wbrightness, 0, 1, 0.01)
				gamegui.SliderFloat32("Contrast", &state.wcontrast, 0, 1, 0.01)
				gamegui.EndGroup()
			}
			if gamegui.BeginGroup("Post Process", 80) {
				gamegui.SliderInt32("Threshold", &state.threshold, 0, 255, 1)
				gamegui.EndGroup()
			}
		}
		gamegui.EndWindow()
		gamegui.End()
	}
	window.RunMainLoop(renderloop)
}

package main

import (
	"fmt"
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

// State holds the state of the application
type State struct {
	// perlin
	Useperlin    bool
	Poctaves     int32
	Presolution  int32
	Pbrightness  float32
	Pcontrast    float32
	Pz           int32
	Pscale       float32
	Ppersistance float32
	// worley
	Useworley    bool
	Woctaves     int32
	Wresolution  int32
	Wradius      float32
	Wbrightness  float32
	Wcontrast    float32
	Wscale       float32
	Wpersistance float32
	// post processing
	Operation1 int32
	Operation2 int32
	Threshold  float32
}

// initializeState sets the state with initial values
func initializeState() *State {
	state := &State{
		// perlin
		Useperlin:    true,
		Poctaves:     1,
		Presolution:  1,
		Pbrightness:  0,
		Pcontrast:    0.5,
		Pz:           1,
		Pscale:       1,
		Ppersistance: 1,
		// worley
		Useworley:    true,
		Woctaves:     1,
		Wresolution:  16,
		Wradius:      70,
		Wbrightness:  0,
		Wcontrast:    0.5,
		Wscale:       1,
		Wpersistance: 1,
		// general
		Operation1: 0,
		Operation2: 0,
		Threshold:  0.5,
	}
	return state
}

func loadState(state *State) {
	err := Load("./state.json", state)
	if err != nil {
		fmt.Println("Couldn't find state.json, use initial state instead")
	}
}
func saveState(state *State) {
	err := Save("./state.json", state)
	if err != nil {
		panic(err)
	}
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
	state := initializeState()
	loadState(state)
	saveState(state)

	// make textures
	image := createAndFillImage(1024*1024*4, 120)
	weathermaptexture, err := texture.MakeFromData(image, 1024, 1024, gl.RGBA32F, gl.RGBA)
	if err != nil {
		panic(err)
	}
	perlintexture, err := texture.MakeFromData(image, 1024, 1024, gl.RGBA32F, gl.RGBA)
	if err != nil {
		panic(err)
	}
	worleytexture, err := texture.MakeFromData(image, 1024, 1024, gl.RGBA32F, gl.RGBA)
	if err != nil {
		panic(err)
	}

	// make render pass
	renderpass := MakeRenderpass(SHADER_PATH)

	// create worley noise generator
	clear := MakeClear(SHADER_PATH)
	perlin := MakePerlin(SHADER_PATH)
	worley := MakeWorley(SHADER_PATH)
	merge := MakeMerge(SHADER_PATH)
	postprocess := MakePostProcess(SHADER_PATH)

	// create gui
	gamegui := gui.New(window.Window)

	// operations
	operations := make([]string, 3)
	operations[0] = "Lerp"
	operations[1] = "Multiply"
	operations[2] = "Mapping"

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// clear weather texture
		clear.ClearTexture(&weathermaptexture)

		// generate perlin
		if state.Useperlin {
			perlin.UpdateState(state)
			perlin.GenerateTexture(&perlintexture)
			merge.UpdateState(state.Operation1)
			merge.MergeTextures(&weathermaptexture, &perlintexture, &weathermaptexture)
		}

		// generate worley
		if state.Useworley {
			worley.UpdateState(state)
			worley.GenerateTexture(&worleytexture)
			merge.UpdateState(state.Operation2)
			merge.MergeTextures(&weathermaptexture, &worleytexture, &weathermaptexture)
		}

		postprocess.UpdateState(state)
		postprocess.Apply(&weathermaptexture, &weathermaptexture)

		renderpass.Render(&camera, &weathermaptexture)

		// gui
		gamegui.Begin()
		if gamegui.BeginWindow("Options", 0, 0, 250, float32(HEIGHT)) {
			if gamegui.BeginGroup("Perlin", 350) {
				gamegui.Checkbox("Use Perlin", &state.Useperlin)
				gamegui.SliderInt32("Resolution", &state.Presolution, 0, 10, 1)
				gamegui.SliderInt32("Z", &state.Pz, 0, 1024, 1)
				gamegui.Label("Fbm")
				gamegui.SliderInt32("Octaves", &state.Poctaves, 1, 5, 1)
				gamegui.SliderFloat32("Scale", &state.Pscale, 1, 3, 0.01)
				gamegui.SliderFloat32("Persistance", &state.Ppersistance, 0, 2, 0.01)
				gamegui.Label("Post")
				gamegui.SliderFloat32("Brightness", &state.Pbrightness, 0, 1, 0.01)
				gamegui.SliderFloat32("Contrast", &state.Pcontrast, 0, 1, 0.01)
				gamegui.EndGroup()
			}
			if gamegui.BeginGroup("Worley", 360) {
				gamegui.Checkbox("Use Worley", &state.Useworley)
				gamegui.SliderInt32("Resolution", &state.Wresolution, 1, 64, 1)
				gamegui.SliderFloat32("Radius", &state.Wradius, 0, 200, 0.1)
				gamegui.Label("Fbm")
				gamegui.SliderInt32("Octaves", &state.Woctaves, 1, 5, 1)
				gamegui.SliderFloat32("Scale", &state.Wscale, 1, 3, 0.01)
				gamegui.SliderFloat32("Persistance", &state.Wpersistance, 0, 1, 0.01)
				gamegui.Label("Post")
				gamegui.SliderFloat32("Brightness", &state.Wbrightness, 0, 1, 0.01)
				gamegui.SliderFloat32("Contrast", &state.Wcontrast, 0, 1, 0.01)
				gamegui.EndGroup()
			}
			if gamegui.BeginGroup("Post Process", 150) {
				gamegui.Selector("Merge Op1", operations, &state.Operation1)
				gamegui.Selector("Merge Op2", operations, &state.Operation2)
				gamegui.SliderFloat32("Threshold", &state.Threshold, 0, 1, 0.01)
				gamegui.EndGroup()
			}
		}
		gamegui.EndWindow()
		gamegui.End()
	}
	window.RunMainLoop(renderloop)

	saveState(state)
}

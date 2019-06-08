package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/shader"
	"github.com/adrianderstroff/realtime-clouds/pkg/scene/camera"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/mesh/plane"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type RaymarchingPass struct {
	cloudbasefbo   texture.Texture
	clouddetailfbo texture.Texture
	turbulencefbo  texture.Texture
	cloudmapfbo    texture.Texture
	raymarchshader shader.Shader
	// uniform variables
	globaldensity  float32
	globalcoverage float32
}

func MakeRaymarchingPass(width, height int, texpath, shaderpath string) RaymarchingPass {
	// create textures
	cloudbasefbo, err := texture.Make3DFromPath(MakePathsFromDirectory(texpath+"cloud-base/", "base", "png", 0, 127), gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	clouddetailfbo, err := texture.Make3DFromPath(MakePathsFromDirectory(texpath+"cloud-detail/", "detail", "png", 0, 31), gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	turbulencefbo, err := texture.MakeFromPath(texpath+"cloud-turbulence/turbulence.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	cloudmapfbo, err := texture.MakeFromPath(texpath+"cloud-map/cloud-map.png", gl.RGBA, gl.RGBA)
	//cloudmapfbo, err := texture.MakeFromPath(texpath+"debug.jpg", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}

	// change wrap to repeat
	cloudbasefbo.SetWrap3D(gl.REPEAT, gl.REPEAT, gl.REPEAT)
	clouddetailfbo.SetWrap3D(gl.REPEAT, gl.REPEAT, gl.REPEAT)
	turbulencefbo.SetWrap2D(gl.REPEAT, gl.REPEAT)
	cloudmapfbo.SetWrap2D(gl.REPEAT, gl.REPEAT)

	// create shaders
	plane := plane.Make(2, 2, gl.TRIANGLES)
	raymarchshader, err := shader.Make(shaderpath+"/realtimeclouds/clouds.vert", shaderpath+"/realtimeclouds/clouds.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(plane)

	return RaymarchingPass{
		cloudbasefbo:   cloudbasefbo,
		clouddetailfbo: clouddetailfbo,
		turbulencefbo:  turbulencefbo,
		cloudmapfbo:    cloudmapfbo,
		raymarchshader: raymarchshader,
		// uniform variables
		globaldensity:  0.5,
		globalcoverage: 0.5,
	}
}

func (rmp *RaymarchingPass) Render(camera camera.Camera, time float32) {
	rmp.cloudbasefbo.Bind(0)
	rmp.clouddetailfbo.Bind(1)
	rmp.turbulencefbo.Bind(2)
	rmp.cloudmapfbo.Bind(3)

	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateVec3("uCamera.pos", camera.GetPos())
	rmp.raymarchshader.UpdateMat4("uCamera.V", camera.GetView())
	rmp.raymarchshader.UpdateMat4("uCamera.P", camera.GetPerspective())
	rmp.raymarchshader.UpdateFloat32("uCamera.fov", 45.0)
	rmp.raymarchshader.UpdateFloat32("uCamera.aspect", 800.0/600.0)
	rmp.raymarchshader.UpdateMat4("M", mgl32.Ident4())
	rmp.raymarchshader.UpdateFloat32("uTime", time)
	rmp.raymarchshader.Render()
	rmp.raymarchshader.Release()

	rmp.cloudbasefbo.Unbind()
	rmp.clouddetailfbo.Unbind()
	rmp.turbulencefbo.Unbind()
	rmp.cloudmapfbo.Unbind()
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (rmp *RaymarchingPass) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (rmp *RaymarchingPass) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (rmp *RaymarchingPass) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (rmp *RaymarchingPass) OnKeyPress(key, action, mods int) bool {
	// update global density
	if key == int(glfw.KeyQ) {
		rmp.globaldensity -= 0.01
	} else if key == int(glfw.KeyE) {
		rmp.globaldensity += 0.01
	}
	rmp.globaldensity = cgm.Clamp(rmp.globaldensity, 0, 1)

	// update global coverage
	if key == int(glfw.KeyZ) {
		rmp.globalcoverage -= 0.01
	} else if key == int(glfw.KeyC) {
		rmp.globalcoverage += 0.01
	}
	rmp.globalcoverage = cgm.Clamp(rmp.globalcoverage, 0, 1)

	// update uniforms
	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateFloat32("uGlobalDensity", rmp.globaldensity)
	rmp.raymarchshader.UpdateFloat32("uGlobalCoverage", rmp.globalcoverage)
	rmp.raymarchshader.Release()

	return false
}

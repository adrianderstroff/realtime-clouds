// Package shader provides a way to load shader programs, adding renderable
// objects to the shader and updating values of the shader as well as
// executing the shader.
package shader

// Renderable is an object that can be drawn by a renderer.
type Renderable interface {
	Build(shaderProgramHandle uint32)
	Render()
	RenderInstanced(instancecount int32)
}

// Package core provides an abstraction layer on top of OpenGL.
// It contains entities that provide utilities to simplify rendering.
package core

// Renderable is an object that can be drawn by a renderer.
type Renderable interface {
	Build(shaderProgramHandle uint32)
	Render()
	RenderInstanced(instancecount int32)
}

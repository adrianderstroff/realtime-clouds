// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

// Renderable is an object that can be drawn by a renderer.
type Renderable interface {
	Build(shaderProgramHandle uint32)
	Render()
	RenderInstanced(instancecount int32)
}

// Cube is a Renderable that can be specified with different dimensions..
type Cube struct {
	vao VAO
}

// MakeCube creates a Cube with the specified half width, height and depth.
func MakeCube(halfWidth, halfHeight, halfDepth float32) Cube {
	// vertex positions
	v1 := []float32{-halfWidth, halfHeight, halfDepth}
	v2 := []float32{-halfWidth, -halfHeight, halfDepth}
	v3 := []float32{halfWidth, halfHeight, halfDepth}
	v4 := []float32{halfWidth, -halfHeight, halfDepth}
	v5 := []float32{-halfWidth, halfHeight, -halfDepth}
	v6 := []float32{-halfWidth, -halfHeight, -halfDepth}
	v7 := []float32{halfWidth, halfHeight, -halfDepth}
	v8 := []float32{halfWidth, -halfHeight, -halfDepth}
	vertices := combine(
		// front
		v1, v2, v3,
		v3, v2, v4,
		// back
		v7, v8, v5,
		v5, v8, v6,
		// left
		v5, v6, v1,
		v1, v6, v2,
		// right
		v3, v4, v7,
		v7, v4, v8,
		// top
		v5, v1, v7,
		v7, v1, v3,
		// bottom
		v2, v6, v4,
		v4, v6, v8,
	)
	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := repeat(combine(t1, t2, t3, t3, t2, t4), 6)
	// normals
	right := []float32{1.0, 0.0, 0.0}
	left := []float32{-1.0, 0.0, 0.0}
	top := []float32{0.0, 1.0, 0.0}
	bottom := []float32{0.0, -1.0, 0.0}
	front := []float32{0.0, 0.0, -1.0}
	back := []float32{0.0, 0.0, 1.0}
	normals := combine(
		repeat(bottom, 6),
		repeat(top, 6),
		repeat(left, 6),
		repeat(right, 6),
		repeat(front, 6),
		repeat(back, 6),
	)

	vertexBuffer := MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vertexBuffer.AddVertexAttribute("vert", 3, gl.FLOAT)
	uvBuffer := MakeVBO(uvs, 2, gl.STATIC_DRAW)
	uvBuffer.AddVertexAttribute("uv", 2, gl.FLOAT)
	normalBuffer := MakeVBO(normals, 3, gl.STATIC_DRAW)
	normalBuffer.AddVertexAttribute("normal", 3, gl.FLOAT)

	vao := MakeVAO(gl.TRIANGLES)
	vao.AddVertexBuffer(&vertexBuffer)
	vao.AddVertexBuffer(&uvBuffer)
	vao.AddVertexBuffer(&normalBuffer)

	return Cube{vao}
}

// Delete destroys this Renderable.
func (cube *Cube) Delete() {
	cube.vao.Delete()
}

// Build prepares the vertex buffers for rendering.
// It is called by the RenderProgram after adding it using AddRenderable
// thus it is usually not advised to call it on your own.
func (cube Cube) Build(shaderProgramHandle uint32) {
	cube.vao.BuildBuffers(shaderProgramHandle)
}

// Render draws the geometry of this Renderable using the currently bound shader.
func (cube Cube) Render() {
	cube.vao.Render()
}

// RenderInstanced draws the geometry of this Renderable multiple times according to the instancecount.
func (cube Cube) RenderInstanced(instancecount int32) {
	cube.vao.RenderInstanced(instancecount)
}

// Skybox is a cube of size 2 ranging from -1 to 1 in all three dimensions.
// It requires a cubemap as texture for its six sides.
type Skybox struct {
	vao     VAO
	texture Texture
}

// MakeSkybox constructs a Skybox object with the given cubemapTexture.
func MakeSkybox(cubeTexture Texture) Skybox {
	var size float32 = 1.0
	v0 := []float32{-size, -size, -size}
	v1 := []float32{-size, -size, size}
	v2 := []float32{size, -size, size}
	v3 := []float32{size, -size, -size}
	v4 := []float32{-size, size, -size}
	v5 := []float32{-size, size, size}
	v6 := []float32{size, size, size}
	v7 := []float32{size, size, -size}
	vertices := combine(
		// right face
		v2, v7, v3, v2, v6, v7,
		// left face
		v0, v4, v5, v0, v5, v1,
		// top face
		v7, v6, v5, v7, v5, v4,
		// bottom face
		v0, v1, v2, v0, v2, v3,
		// back face
		v0, v7, v4, v0, v3, v7,
		// front face
		v6, v2, v5, v5, v2, v1,
	)
	vbo := MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vbo.AddVertexAttribute("vert", 3, gl.FLOAT)
	vao := MakeVAO(gl.TRIANGLES)
	vao.AddVertexBuffer(&vbo)

	return Skybox{vao, cubeTexture}
}

// Delete destroys this Renderable as well as the cubemapTexture.
func (skybox *Skybox) Delete() {
	skybox.vao.Delete()
	skybox.texture.Delete()
}

// Build prepares the vertex buffers for rendering.
// It is called by the RenderProgram after adding it using AddRenderable
// thus it is usually not advised to call it on your own.
func (skybox Skybox) Build(shaderProgramHandle uint32) {
	skybox.vao.BuildBuffers(shaderProgramHandle)
}

// Render draws the geometry of this Renderable using the currently bound shader.
func (skybox Skybox) Render() {
	gl.DepthMask(false)
	skybox.texture.Bind(0)
	skybox.vao.Render()
	skybox.texture.Unbind()
	gl.DepthMask(true)
}

// RenderInstanced draws the geometry of this Renderable multiple times according to the instancecount.
func (skybox Skybox) RenderInstanced(instancecount int32) {
	gl.DepthMask(false)
	skybox.texture.Bind(0)
	skybox.vao.RenderInstanced(instancecount)
	skybox.texture.Unbind()
	gl.DepthMask(true)
}

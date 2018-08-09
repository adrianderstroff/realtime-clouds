package skybox

import (
	"github.com/go-gl/gl/v4.3-core/gl"

	buf "github.com/adrianderstroff/realtime-clouds/pkg/buffer"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/texture"
)

// Skybox is a cube of size 2 ranging from -1 to 1 in all three dimensions.
// It requires a cubemap as texture for its six sides.
type Skybox struct {
	vao     buf.VAO
	texture tex.Texture
}

// MakeSkybox constructs a Skybox object with the given cubemapTexture.
func MakeSkybox(cubeTexture tex.Texture) Skybox {
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
	vbo := buf.MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vbo.AddVertexAttribute("vert", 3, gl.FLOAT)
	vao := buf.MakeVAO(gl.TRIANGLES)
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

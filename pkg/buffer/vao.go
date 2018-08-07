// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

// VAO is a buffer that holds multiple vertex buffers and zero to one index buffer.
type VAO struct {
	handle        uint32
	mode          uint32
	vertexBuffers []*VBO
	indexBuffer   *IBO
}

// Creates a new VAO.
// 'mode' specified the drawing mode used.
// Some modes would be TRIANGLE, TRIANGLE_STRIP, TRIANGLE_FAN
func MakeVAO(mode uint32) VAO {
	vao := VAO{0, mode, nil, nil}
	gl.GenVertexArrays(1, &vao.handle)
	return vao
}

// Delete destroys this and all vertex and index buffers associated with this VAO.
func (vao *VAO) Delete() {
	// delete buffers
	if vao.vertexBuffers != nil {
		for _, vertBuf := range vao.vertexBuffers {
			vertBuf.Delete()
		}
	}
	vao.indexBuffer.Delete()

	// delete vertex array
	gl.DeleteVertexArrays(1, &vao.handle)
}

// Render draws the geometry.
// It uses indexed rendering if a index buffer is present.
func (vao *VAO) Render() {
	gl.BindVertexArray(vao.handle)
	if vao.indexBuffer != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.indexBuffer.handle)
		gl.DrawElements(vao.mode, vao.indexBuffer.count, gl.UNSIGNED_SHORT, nil)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		gl.DrawArrays(vao.mode, 0, vao.vertexBuffers[0].count)
	}
	gl.BindVertexArray(0)
}

// RenderInstanced draws the geomtry multiple times defined by the instancecount.
// It uses indexed rendering if a index buffer is present.
func (vao *VAO) RenderInstanced(instancecount int32) {
	gl.BindVertexArray(vao.handle)
	if vao.indexBuffer != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.indexBuffer.handle)
		gl.DrawElementsInstanced(vao.mode, vao.indexBuffer.count, gl.UNSIGNED_SHORT, nil, instancecount)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		gl.DrawArraysInstanced(vao.mode, 0, vao.vertexBuffers[0].count, instancecount)
	}
	gl.BindVertexArray(0)
}

// AddVertexBuffer adds a vertex buffer at the end.
func (vao *VAO) AddVertexBuffer(vbo *VBO) {
	vao.vertexBuffers = append(vao.vertexBuffers, vbo)
}

// AddIndexBuffer sets the index buffer.
func (vao *VAO) AddIndexBuffer(ibo *IBO) {
	vao.indexBuffer = ibo
}

// GetVertexBuffer returns the vertex buffer at the specifed index.
func (vao *VAO) GetVertexBuffer(idx int) *VBO {
	return vao.vertexBuffers[idx]
}

// GetIndexBuffer returns the only index buffer.
func (vao *VAO) GetIndexBuffer() *IBO {
	return vao.indexBuffer
}

// BuildBuffers gets called by the Shader to setup all added buffers.
func (vao *VAO) BuildBuffers(shaderProgramHandle uint32) {
	gl.BindVertexArray(vao.handle)
	for _, vbo := range vao.vertexBuffers {
		vbo.BuildVertexAttributes(shaderProgramHandle)
	}
	gl.BindVertexArray(0)
}

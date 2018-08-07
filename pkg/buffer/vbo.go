// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

// VertexAttribute specifies the layout of the vertex buffer.
type VertexAttribute struct {
	name   string
	count  int32
	glType uint32
}

// VBO is a buffer that stores vertex attributes.
type VBO struct {
	handle     uint32
	count      int32
	stride     uint32
	usage      uint32
	attributes []VertexAttribute
}

// MakeVBO construct a VBO width the specified data and the size of one element.
// The usage gives the GPU a hint how to treat the data usage.
// Usage patterns are GL_STREAM_DRAW, GL_STREAM_READ, GL_STREAM_COPY, GL_STATIC_DRAW,
// GL_STATIC_READ, GL_STATIC_COPY, GL_DYNAMIC_DRAW, GL_DYNAMIC_READ, or GL_DYNAMIC_COPY
func MakeVBO(data []float32, elementsPerVertex uint32, usage uint32) VBO {
	vbo := VBO{
		handle:     0,
		count:      int32(len(data)) / int32(elementsPerVertex),
		stride:     elementsPerVertex,
		usage:      usage,
		attributes: nil,
	}

	gl.GenBuffers(1, &vbo.handle)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
	if len(data) != 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), usage)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, 4, gl.Ptr([]float32{0.0}), usage)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vbo
}

// UpdateData replaces the previous data with this data.
// Make sure that the data follows the same layout as specified by the vertex attributes.
func (vbo *VBO) UpdateData(data []float32) {
	// update buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), vbo.usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// update size
	vbo.count = int32(len(data)) / int32(vbo.stride)
}

// Delete detroys this VBO and all vertex attributes.
func (vbo *VBO) Delete() {
	vbo.count = 0
	vbo.stride = 0
	vbo.attributes = nil
	gl.DeleteBuffers(1, &vbo.handle)
}

// AddVertexAttribute adds a new vertex layout for the given name the number of elements and the data type of the elements.
func (vbo *VBO) AddVertexAttribute(name string, count int32, glType uint32) {
	vbo.attributes = append(vbo.attributes, VertexAttribute{name, count, glType})
}

// BuildVertexAttributes is called by the shader and binds the vertex data to the variable specified by name.
func (vbo *VBO) BuildVertexAttributes(shaderProgramHandle uint32) {
	// specify all vertex attributes
	var offset int = 0
	for _, attrib := range vbo.attributes {
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
		location := gl.GetAttribLocation(shaderProgramHandle, gl.Str(attrib.name+"\x00"))
		if location != -1 {
			gl.EnableVertexAttribArray(uint32(location))
			gl.VertexAttribPointer(uint32(location), attrib.count, attrib.glType, false, int32(vbo.stride*4), gl.PtrOffset(offset*4))
		}
		offset += int(attrib.count)
	}

	// unbind vbo to prevent overwrites
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

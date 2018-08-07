// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import "github.com/go-gl/gl/v4.3-core/gl"

// IBO contains indices to vertex attributes.
// It has to be used together with a VAO.
type IBO struct {
	handle uint32
	count  int32
}

// MakeIBO constructs an IBO with the indices specified in data and the usage.
func MakeIBO(data []uint16, usage uint32) IBO {
	ibo := IBO{0, int32(len(data))}
	gl.GenBuffers(1, &ibo.handle)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.handle)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data)*2, gl.Ptr(data), usage)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	return ibo
}

// Delete detroys this IBO.
func (ibo *IBO) Delete() {
	gl.DeleteBuffers(1, &ibo.handle)
	ibo.count = 0
}

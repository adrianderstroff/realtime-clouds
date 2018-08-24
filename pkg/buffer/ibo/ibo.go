// Package ibo contains a buffer with indices of a vertex array.
package ibo

import (
	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
)

// IBO contains indices to vertex attributes.
// It has to be used together with a VAO.
type IBO struct {
	handle uint32
	count  int32
}

// Make constructs an IBO with the indices specified in data and the usage.
func Make(data []uint16, usage uint32) IBO {
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

// Len returns the size of the buffer.
func (ibo *IBO) Len() int32 {
	return ibo.count
}

// GetHandle returns the OpenGL handle of this buffer.
func (ibo *IBO) GetHandle() uint32 {
	return ibo.handle
}

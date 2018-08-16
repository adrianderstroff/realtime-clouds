// Package geometry provides geometric primitives that can be used in meshes.
// It also offers a way to create custom geometric shapes.
package geometry

import (
	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	geometry "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry"
)

// MakeQuad creates a Quad with the specified width, height and depth.
// If the normals should be inside the quad the inside parameter should be true.
func MakeQuad(width, height, depth float32, inside bool) geometry.Geometry {
	// half side lengths
	halfWidth := width / 2.0
	halfHeight := height / 2.0
	halfDepth := depth / 2.0

	// vertex positions
	v1 := []float32{-halfWidth, halfHeight, halfDepth}
	v2 := []float32{-halfWidth, -halfHeight, halfDepth}
	v3 := []float32{halfWidth, halfHeight, halfDepth}
	v4 := []float32{halfWidth, -halfHeight, halfDepth}
	v5 := []float32{-halfWidth, halfHeight, -halfDepth}
	v6 := []float32{-halfWidth, -halfHeight, -halfDepth}
	v7 := []float32{halfWidth, halfHeight, -halfDepth}
	v8 := []float32{halfWidth, -halfHeight, -halfDepth}
	positions := geometry.Combine(
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
	uvs := geometry.Repeat(geometry.Combine(t1, t2, t3, t3, t2, t4), 6)

	// normals
	right := []float32{1.0, 0.0, 0.0}
	left := []float32{-1.0, 0.0, 0.0}
	top := []float32{0.0, 1.0, 0.0}
	bottom := []float32{0.0, -1.0, 0.0}
	front := []float32{0.0, 0.0, -1.0}
	back := []float32{0.0, 0.0, 1.0}
	// swap normals if inside is true
	if inside {
		right, left = left, right
		top, bottom = bottom, top
		front, back = back, front
	}
	normals := geometry.Combine(
		geometry.Repeat(bottom, 6),
		geometry.Repeat(top, 6),
		geometry.Repeat(left, 6),
		geometry.Repeat(right, 6),
		geometry.Repeat(front, 6),
		geometry.Repeat(back, 6),
	)

	// setup data
	data := [][]float32{
		positions,
		uvs,
		normals,
	}

	// setup layout
	layout := []geometry.VertexAttribute{
		geometry.MakeVertexAttribute("pos", gl.FLOAT, 3, gl.STATIC_DRAW),
		geometry.MakeVertexAttribute("uv", gl.FLOAT, 2, gl.STATIC_DRAW),
		geometry.MakeVertexAttribute("normal", gl.FLOAT, 3, gl.STATIC_DRAW),
	}

	return geometry.Make(layout, data)
}

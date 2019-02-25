// Package geometry provides geometric primitives that can be used in meshes.
// It also offers a way to create custom geometric shapes.
package quad

import (
	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	geometry "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry"
)

// Make creates a Quad with the specified width and height on the x-y plane
// with the normal pointing up the y-axis
func Make(width, height float32) geometry.Geometry {
	// half side lengths
	halfWidth := width / 2.0
	halfHeight := height / 2.0

	// vertex positions
	v1 := []float32{-halfWidth, 0, halfHeight}
	v2 := []float32{-halfWidth, 0, -halfHeight}
	v3 := []float32{halfWidth, 0, halfHeight}
	v4 := []float32{halfWidth, 0, -halfHeight}
	positions := geometry.Combine(v1, v2, v3, v3, v2, v4)

	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := geometry.Combine(t1, t2, t3, t3, t2, t4)

	// normals
	up := []float32{0.0, 1.0, 0.0}
	normals := geometry.Repeat(up, 6)

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

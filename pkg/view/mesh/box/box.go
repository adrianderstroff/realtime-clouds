// Package box is used for creating a simple box mesh.
package box

import (
	mesh "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
)

// Make constructs a box with the specified dimensions. If inside is true
// then the triangles are specified in an order in which the normals will
// point inwards.
func Make(width, height, depth float32, inside bool, mode uint32) mesh.Mesh {
	geometry := cube.Make(width, height, depth, inside)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

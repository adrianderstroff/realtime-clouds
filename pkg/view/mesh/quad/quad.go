// Package mesh is used for creating meshes from geometry and textures.
// Meshes are entities that can be assigned to a ShaderProgram in order to render them.
package mesh

import (
	quad "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry/quad"
	mesh "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
)

func Make(width, height, depth float32, inside bool, mode uint32) mesh.Mesh {
	geometry := quad.Make(width, height, depth, inside)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

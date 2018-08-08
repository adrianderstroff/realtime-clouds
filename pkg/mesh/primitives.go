// Package mesh is used for creating meshes from geometry and textures.
// Meshes are entities that can be assigned to a ShaderProgram in order to render them.
package mesh

import (
	geom "github.com/adrianderstroff/realtime-clouds/pkg/geometry"
)

func MakeQuad(width, height, depth float32, inside bool, mode uint32) Mesh {
	geometry := geom.MakeQuad(width, height, depth, inside)
	mesh := MakeMesh(geometry, nil, mode)
	return mesh
}

// Package box is used for creating a simple box mesh.
package plane

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/view/geometry/quad"
	mesh "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
)

// Make constructs a plane with the specified dimensions. The plane is on the
// x-y axis and the normal points up the y-axis.
func Make(width, height float32, mode uint32) mesh.Mesh {
	geometry := quad.Make(width, height)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

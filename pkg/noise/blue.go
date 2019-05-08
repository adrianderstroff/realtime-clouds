package noise

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

// ToroidalDistance2D returns the distance between p1 and p2 in toroidal space
// p1 and p2 have to be between 0 and 1
func ToroidalDistance2D(p1, p2 mgl32.Vec2) float32 {
	dx := cgm.Abs32(p2.X() - p1.X())
	dy := cgm.Abs32(p2.Y() - p1.Y())

	if dx > 0.5 {
		dx = 1.0 - dx
	}

	if dy > 0.5 {
		dy = 1.0 - dy
	}

	return cgm.Sqrt32(dx*dx + dy*dy)
}

// Blue returns a blue noise image of size (width, height) with n points
func Blue(width, height, n int) []uint8 {
	data := []uint8{}

	return data
}

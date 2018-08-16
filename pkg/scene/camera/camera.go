package scene

import "github.com/go-gl/mathgl/mgl32"

// Camera abstracts a camera model with either perspective or orthographic projection.
type Camera interface {
	Update()
	Rotate(theta, phi float32)
	Zoom(distance float32)

	GetView() mgl32.Mat4
	GetPerspective() mgl32.Mat4
	GetViewPerspective() mgl32.Mat4
}

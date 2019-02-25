// Package camera provides implementations of different camera models.
package trackball

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

var MIN_THETA = 0.000001
var MAX_THETA = math.Pi - MIN_THETA

// Trackball moves on a sphere around a target point with a specified radius.
type Trackball struct {
	width  int
	height int
	radius float32
	theta  float32
	phi    float32

	Pos    mgl32.Vec3
	Target mgl32.Vec3
	Up     mgl32.Vec3
	Fov    float32
	Near   float32
	Far    float32

	leftButtonPressed bool
}

// MakeDefault creates a Trackball with the viewport of width and height and a radius from the origin.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func MakeDefault(width, height int, radius float32) Trackball {
	return Make(
		width, height, radius,
		mgl32.Vec3{0.0, 0.0, 0.0}, 45,
		0.1, 100.0,
	)
}

// NewDefault creates a reference to a Trackball with the viewport of width and height and a radius from the origin.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func NewDefault(width, height int, radius float32) *Trackball {
	return New(
		width, height, radius,
		mgl32.Vec3{0.0, 0.0, 0.0}, 45,
		0.1, 100.0,
	)
}

// Make creates a Trackball with the viewport of width and height, the radius from the target,
// the target position the camera is orbiting around, the field of view and the distance of the near and far plane.
func Make(width, height int, radius float32, target mgl32.Vec3, fov, near, far float32) Trackball {
	camera := Trackball{
		width:  width,
		height: height,
		radius: radius,
		theta:  90.0,
		phi:    90.0,
		Target: target,
		Fov:    fov,
		Near:   near,
		Far:    far,
	}
	camera.Update()

	return camera
}

// New creates a reference to a Trackball with the viewport of width and height, the radius from the target,
// the target position the camera is orbiting around, the field of view and the distance of the near and far plane.
func New(width, height int, radius float32, target mgl32.Vec3, fov, near, far float32) *Trackball {
	camera := Make(width, height, radius, target, fov, near, far)
	return &camera
}

// Update recalculates the position of the camera.
// Call it  every time after calling Rotate or Zoom.
func (camera *Trackball) Update() {
	theta := mgl32.DegToRad(camera.theta)
	phi := mgl32.DegToRad(camera.phi)

	// limit angles
	theta = float32(math.Max(float64(theta), MIN_THETA))
	theta = float32(math.Min(float64(theta), MAX_THETA))

	// sphere coordinates
	btheta := float64(theta)
	bphi := float64(phi)
	pos := mgl32.Vec3{
		camera.radius * float32(math.Sin(btheta)*math.Cos(bphi)),
		camera.radius * float32(math.Cos(btheta)),
		camera.radius * float32(math.Sin(btheta)*math.Sin(bphi)),
	}
	camera.Pos = pos.Add(camera.Target)

	look := camera.Pos.Sub(camera.Target).Normalize()
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
	right := worldUp.Cross(look)
	camera.Up = look.Cross(right)
}

// Rotate adds delta angles in degrees to the theta and phi angles.
// Where theta is the vertical angle and phi the horizontal angle.
func (camera *Trackball) Rotate(theta, phi float32) {
	camera.theta += theta
	camera.phi += phi
}

// Zoom changes the radius of the camera to the target point.
func (camera *Trackball) Zoom(distance float32) {
	camera.radius -= distance
	// limit radius
	if camera.radius < 0.1 {
		camera.radius = 0.1
	}
}

// GetPos returns the position of the camera in worldspace
func (camera *Trackball) GetPos() mgl32.Vec3 {
	return camera.Pos
}

// GetView returns the view matrix of the camera.
func (camera *Trackball) GetView() mgl32.Mat4 {
	return mgl32.LookAtV(camera.Pos, camera.Target, camera.Up)
}

// GetPerspective returns the perspective projection of the camera.
func (camera *Trackball) GetPerspective() mgl32.Mat4 {
	fov := mgl32.DegToRad(camera.Fov)
	aspect := float32(camera.width) / float32(camera.height)
	return mgl32.Perspective(fov, aspect, camera.Near, camera.Far)
}

// GetOrtho returns the orthographic projection of the camera.
func (camera *Trackball) GetOrtho() mgl32.Mat4 {
	angle := camera.Fov * math.Pi / 180.0
	dfar := float32(math.Tan(float64(angle/2.0))) * camera.Far
	d := dfar
	return mgl32.Ortho(-d, d, -d, d, camera.Near, camera.Far)
}

// GetViewPerspective returns P*V.
func (camera *Trackball) GetViewPerspective() mgl32.Mat4 {
	return camera.GetPerspective().Mul4(camera.GetView())
}

// SetPos updates the target point of the camera.
// It requires to call Update to take effect.
func (camera *Trackball) SetPos(pos mgl32.Vec3) {
	camera.Target = pos
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (camera *Trackball) OnCursorPosMove(x, y, dx, dy float64) bool {
	if camera.leftButtonPressed {
		dPhi := float32(-dx) / 2.0
		dTheta := float32(-dy) / 2.0
		camera.Rotate(dTheta, -dPhi)
	}
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (camera *Trackball) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	camera.leftButtonPressed = leftPressed
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (camera *Trackball) OnMouseScroll(x, y float64) bool {
	camera.Zoom(float32(y))
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (camera *Trackball) OnKeyPress(key, action, mods int) bool {
	return false
}

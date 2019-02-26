// Package camera provides implementations of different camera models.
package fps

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// FPS moves in the view direction while the viewing direction can be changed.
type FPS struct {
	width  int
	height int
	theta  float32
	phi    float32
	dir    mgl32.Vec3
	speed  float32

	Pos    mgl32.Vec3
	Target mgl32.Vec3
	Up     mgl32.Vec3
	Right  mgl32.Vec3
	Fov    float32
	Near   float32
	Far    float32
}

// MakeDefault creates a FPS camera with the viewport of width and height and a position.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func MakeDefault(width, height int, pos mgl32.Vec3, speed float32) FPS {
	return Make(width, height, pos, speed, 45, 0.1, 100.0)
}

// NewDefault creates a reference to a FPS camera with the viewport of width and height and a position.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func NewDefault(width, height int, pos mgl32.Vec3, speed float32) *FPS {
	return New(width, height, pos, speed, 45, 0.1, 100.0)
}

// Make creates a FPS with the viewport of width and height and a radius from the origin.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func Make(width, height int, pos mgl32.Vec3, speed, fov, near, far float32) FPS {
	dir := mgl32.Vec3{0.0, 0.0, 1.0}
	camera := FPS{
		width:  width,
		height: height,
		theta:  90.0,
		phi:    90.0,
		dir:    dir,
		speed:  speed,

		Pos:    pos,
		Target: pos.Add(dir),
		Up:     mgl32.Vec3{0, 1, 0},
		Right:  mgl32.Vec3{1, 0, 0},
		Fov:    fov,
		Near:   near,
		Far:    far,
	}
	camera.Update()

	return camera
}

// New creates a reference to a FPS with the viewport of width and height and a radius from the origin.
// It assumes a field of view of 45 degrees and a near and far plane at 0.1 and 100.0 respectively.
func New(width, height int, pos mgl32.Vec3, speed, fov, near, far float32) *FPS {
	camera := Make(width, height, pos, speed, fov, near, far)
	return &camera
}

// Update recalculates the position of the camera.
// Call it  every time after calling Rotate or Zoom.
func (camera *FPS) Update() {
	theta := mgl32.DegToRad(camera.theta)
	phi := mgl32.DegToRad(camera.phi)

	// sphere coordinates with inverse y
	btheta := float64(theta)
	bphi := float64(phi)
	camera.dir = mgl32.Vec3{
		float32(math.Sin(btheta) * math.Cos(bphi)),
		-float32(math.Cos(btheta)),
		float32(math.Sin(btheta) * math.Sin(bphi)),
	}
	camera.dir = camera.dir.Normalize()

	// set target
	camera.Target = camera.Pos.Add(camera.dir)

	// calculate up vector
	look := camera.dir.Mul(-1)
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
	camera.Right = worldUp.Cross(look).Normalize()
	camera.Up = look.Cross(camera.Right)
}

// Rotate adds delta angles in degrees to the theta and phi angles.
// Where theta is the vertical angle and phi the horizontal angle.
func (camera *FPS) Rotate(theta, phi float32) {
	camera.theta += theta
	camera.phi += phi

	// limit angles
	camera.theta = float32(math.Max(math.Min(float64(camera.theta), 179.9), 0.01))
	if camera.phi < 0 {
		camera.phi = 360 + camera.phi
	} else if camera.phi >= 360 {
		camera.phi = camera.phi - 360
	}
}

// Zoom changes the radius of the camera to the target point.
func (camera *FPS) Zoom(distance float32) {}

// GetPos returns the position of the camera in worldspace
func (camera *FPS) GetPos() mgl32.Vec3 {
	return camera.Pos
}

// GetView returns the view matrix of the camera.
func (camera *FPS) GetView() mgl32.Mat4 {
	return mgl32.LookAtV(camera.Pos, camera.Target, camera.Up)
}

// GetPerspective returns the perspective projection of the camera.
func (camera *FPS) GetPerspective() mgl32.Mat4 {
	fov := mgl32.DegToRad(camera.Fov)
	aspect := float32(camera.width) / float32(camera.height)
	return mgl32.Perspective(fov, aspect, camera.Near, camera.Far)
}

// GetOrtho returns the orthographic projection of the camera.
func (camera *FPS) GetOrtho() mgl32.Mat4 {
	angle := camera.Fov * math.Pi / 180.0
	dfar := float32(math.Tan(float64(angle/2.0))) * camera.Far
	d := dfar
	return mgl32.Ortho(-d, d, -d, d, camera.Near, camera.Far)
}

// GetViewPerspective returns P*V.
func (camera *FPS) GetViewPerspective() mgl32.Mat4 {
	return camera.GetPerspective().Mul4(camera.GetView())
}

// SetPos updates the target point of the camera.
// It requires to call Update to take effect.
func (camera *FPS) SetPos(pos mgl32.Vec3) {
	camera.Pos = pos
	camera.Target = camera.Pos.Add(camera.dir)
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (camera *FPS) OnCursorPosMove(x, y, dx, dy float64) bool {
	dPhi := float32(-dx) / 2.0
	dTheta := float32(-dy) / 2.0
	camera.Rotate(dTheta, -dPhi)
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (camera *FPS) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (camera *FPS) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (camera *FPS) OnKeyPress(key, action, mods int) bool {
	dir := camera.dir.Mul(camera.speed)
	right := camera.Right.Mul(camera.speed)
	if key == int(glfw.KeyW) {
		camera.Pos = camera.Pos.Add(dir)
	} else if key == int(glfw.KeyS) {
		camera.Pos = camera.Pos.Sub(dir)
	} else if key == int(glfw.KeyA) {
		camera.Pos = camera.Pos.Sub(right)
	} else if key == int(glfw.KeyD) {
		camera.Pos = camera.Pos.Add(right)
	}
	return false
}

// Package interaction handles the user interaction with mouse and keyboard.
package interaction

import (
	window "github.com/adrianderstroff/realtime-clouds/pkg/core/window"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Interaction holds all eventhandlers and propagates them to it's registered handlers.
type Interaction struct {
	ctx *window.Window

	cursorPosHandlers   []CursorPosHandler
	mouseButtonHandlers []MouseButtonHandler
	mouseScrollHandlers []MouseScrollHandler
	keyPressHandlers    []KeyPressHandler

	prevPosX, prevPosY float64
	posInit            bool
	leftPressed        bool
	rightPressed       bool
	loopCursor         bool
}

// Interactable is an entity that listens to different events and reacts to them.
type Interactable interface {
	OnCursorPosMove(x, y, dx, dy float64) bool
	OnMouseButtonPress(leftPressed, rightPressed bool) bool
	OnMouseScroll(x, y float64) bool
	OnKeyPress(key, action, mods int) bool
}

// CursorPosHandler is called every time the cursor position changes.
type CursorPosHandler func(float64, float64, float64, float64) bool

// MouseButtonHandler is called every time the left or right mouse button is pressed or released.
type MouseButtonHandler func(bool, bool) bool

// MouseScrollHandler is called every time the mouse scroll is used.
type MouseScrollHandler func(float64, float64) bool

// KeyPressHandler is called every time a keyboard key is pressed or released.
type KeyPressHandler func(int, int, int) bool

// Make constructs an Interaction and registers all necessary handlers for the window.
func Make(window *window.Window) Interaction {
	// construct Interaction
	interaction := Interaction{
		ctx: window,

		prevPosX:     0.0,
		prevPosY:     0.0,
		posInit:      false,
		leftPressed:  false,
		rightPressed: false,
		loopCursor:   false,
	}

	// add handlers to window
	window.Window.SetCursorPosCallback(interaction.onCursorPos)
	window.Window.SetCursorEnterCallback(interaction.onCursorEnter)
	window.Window.SetMouseButtonCallback(interaction.onMouseButton)
	window.Window.SetScrollCallback(interaction.onMouseScroll)
	window.Window.SetKeyCallback(interaction.onKeyPress)

	return interaction
}

// EnableCursorLoop hides the cursor and loops it inside the window in x and y direction.
func (interaction *Interaction) EnableCursorLoop() {
	interaction.loopCursor = true
	interaction.ctx.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}

// AddInteractable registers all handlers of the interactable in the Windowmanager.
func (interaction *Interaction) AddInteractable(interactable Interactable) {
	interaction.AddCursorPosHandler(interactable.OnCursorPosMove)
	interaction.AddMouseButtonHandler(interactable.OnMouseButtonPress)
	interaction.AddMouseScrollHandler(interactable.OnMouseScroll)
	interaction.AddKeyPressHandler(interactable.OnKeyPress)
}

// AddCursorPosHandler registers a CursorPosHandler in the Window.
func (interaction *Interaction) AddCursorPosHandler(handler CursorPosHandler) {
	interaction.cursorPosHandlers = append(interaction.cursorPosHandlers, handler)
}

// AddMouseButtonHandler registers a MouseButtonHandler in the Window.
func (interaction *Interaction) AddMouseButtonHandler(handler MouseButtonHandler) {
	interaction.mouseButtonHandlers = append(interaction.mouseButtonHandlers, handler)
}

// AddMouseScrollHandler registers a MouseScrollHandler in the Window.
func (interaction *Interaction) AddMouseScrollHandler(handler MouseScrollHandler) {
	interaction.mouseScrollHandlers = append(interaction.mouseScrollHandlers, handler)
}

// AddKeyPressHandler registers a KeyPressHandler in the Window.
func (interaction *Interaction) AddKeyPressHandler(handler KeyPressHandler) {
	interaction.keyPressHandlers = append(interaction.keyPressHandlers, handler)
}

// onCursorPos receives the cursor x and y pos and propagates it to all CusorPosHandlers.
func (interaction *Interaction) onCursorPos(w *glfw.Window, x float64, y float64) {
	if !interaction.posInit {
		interaction.posInit = true
		interaction.prevPosX = x
		interaction.prevPosY = y
	}
	deltaX := x - interaction.prevPosX
	deltaY := y - interaction.prevPosY
	for _, handler := range interaction.cursorPosHandlers {
		if handler(x, y, deltaX, deltaY) {
			break
		}
	}
	interaction.prevPosX = x
	interaction.prevPosY = y
}

// onMouseButton receives the button the button action and if modifier keys had been pressed and propagates it to all MouseButtonHandlers.
func (interaction *Interaction) onMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	// save pressed button
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			interaction.leftPressed = true
		} else if action == glfw.Release {
			interaction.leftPressed = false
		}
	} else if button == glfw.MouseButtonRight {
		if action == glfw.Press {
			interaction.rightPressed = true
		} else if action == glfw.Release {
			interaction.rightPressed = false
		}
	}

	// inform all handlers
	for _, handler := range interaction.mouseButtonHandlers {
		if handler(interaction.leftPressed, interaction.rightPressed) {
			break
		}
	}
}

// onMouseScroll receives the x and y scroll and propagates it to all MouseScrollHandlers.
func (interaction *Interaction) onMouseScroll(w *glfw.Window, x float64, y float64) {
	for _, handler := range interaction.mouseScrollHandlers {
		if handler(x, y) {
			break
		}
	}
}

// onCursorEnter receives the event whether the cursor entered or left the window.
func (interaction *Interaction) onCursorEnter(w *glfw.Window, entered bool) {
	if !entered {
		interaction.posInit = false

		// loop
		if interaction.loopCursor {
			x := interaction.prevPosX
			y := interaction.prevPosY
			w := float64(interaction.ctx.Width)
			h := float64(interaction.ctx.Height)
			var border float64 = 20

			if x < border {
				x = w - 1
			} else if x > w-border {
				x = 1
			}

			if y < border {
				y = h - 1
			} else if y > h-border {
				y = 1
			}

			interaction.ctx.Window.SetCursorPos(x, y)
		}
	}
}

// onKeyPress receives the pressed button the scan code of the key the key action and if modifier keys had been pressed.
func (interaction *Interaction) onKeyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for _, handler := range interaction.keyPressHandlers {
		if handler(int(key), int(action), int(mods)) {
			break
		}
	}
}

package gui

import (
	"fmt"

	"github.com/adrianderstroff/nuklear/nk"
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowFlags      = nk.WindowMovable | nk.WindowMinimizable | nk.WindowTitle
	groupFlags       = nk.WindowTitle | nk.WindowBorder
	maxVertexBuffer  = 1024 * 1024
	maxElementBuffer = 512 * 1024
)

// GUI is used for rendering a graphical user interface using the underlying
// rendering api used by simgl
type GUI struct {
	ctx *nk.Context
}

// Make creates a new GUI. It attaches itself to the specified window.
func Make(window *glfw.Window) GUI {
	gui := GUI{nil}

	// init nuklear
	gui.ctx = nk.NkPlatformInit(window, nk.PlatformDefault)

	// set style
	styleTable := make([]nk.Color, nk.ColorCount)
	styleTable[nk.ColorText] = nk.NkRgba(190, 190, 190, 255)
	styleTable[nk.ColorWindow] = nk.NkRgba(30, 33, 40, 215)
	styleTable[nk.ColorHeader] = nk.NkRgba(181, 45, 69, 220)
	styleTable[nk.ColorBorder] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorButton] = nk.NkRgba(181, 45, 69, 255)
	styleTable[nk.ColorButtonHover] = nk.NkRgba(190, 50, 70, 255)
	styleTable[nk.ColorButtonActive] = nk.NkRgba(195, 55, 75, 255)
	styleTable[nk.ColorToggle] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorToggleHover] = nk.NkRgba(45, 60, 60, 255)
	styleTable[nk.ColorToggleCursor] = nk.NkRgba(181, 45, 69, 255)
	styleTable[nk.ColorSelect] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorSelectActive] = nk.NkRgba(181, 45, 69, 255)
	styleTable[nk.ColorSlider] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorSliderCursor] = nk.NkRgba(181, 45, 69, 255)
	styleTable[nk.ColorSliderCursorHover] = nk.NkRgba(186, 50, 74, 255)
	styleTable[nk.ColorSliderCursorActive] = nk.NkRgba(191, 55, 79, 255)
	styleTable[nk.ColorProperty] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorEdit] = nk.NkRgba(51, 55, 67, 225)
	styleTable[nk.ColorEditCursor] = nk.NkRgba(190, 190, 190, 255)
	styleTable[nk.ColorCombo] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorChart] = nk.NkRgba(51, 55, 67, 255)
	styleTable[nk.ColorChartColor] = nk.NkRgba(170, 40, 60, 255)
	styleTable[nk.ColorChartColorHighlight] = nk.NkRgba(255, 0, 0, 255)
	styleTable[nk.ColorScrollbar] = nk.NkRgba(30, 33, 40, 255)
	styleTable[nk.ColorScrollbarCursor] = nk.NkRgba(64, 84, 95, 255)
	styleTable[nk.ColorScrollbarCursorHover] = nk.NkRgba(70, 90, 100, 255)
	styleTable[nk.ColorScrollbarCursorActive] = nk.NkRgba(75, 95, 105, 255)
	styleTable[nk.ColorTabHeader] = nk.NkRgba(181, 45, 69, 220)
	nk.NkStyleFromTable(gui.ctx, styleTable)

	// add font atlas
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	fontConfig := nk.NkFontConfig(12)
	fontConfig.SetOversample(1, 1)
	nk.NkFontStashEnd()

	// set padding
	pad := nk.NkVec2(0, 0)
	nk.SetPadding(gui.ctx, pad)
	sp := nk.NkVec2(0, 0)
	nk.SetSpacing(gui.ctx, sp)

	return gui
}

// New returns a pointer to the newly created GUI object. The GUI attaches
// itself to the specified window.
func New(window *glfw.Window) *GUI {
	gui := Make(window)
	return &gui
}

// Begin marks the new frame of rendering the GUI. Make sure that the gui has
// at least one menu item else it will crash the program.
func (gui *GUI) Begin() {
	nk.NkPlatformNewFrame()
}

// End renders all specified GUI objects that had been specified.
func (gui *GUI) End() {
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)

	// reset stuff set by nuklear
	gl.ForcedEnable(gl.CULL_FACE)
	gl.ForcedEnable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
}

// BeginWindow starts a new menu window.
func (gui *GUI) BeginWindow(name string, x, y, width, height float32) bool {
	bounds := nk.NkRect(x, y, width, height)
	update := nk.NkBegin(gui.ctx, name, bounds, windowFlags)
	return update > 0
}

// EndWindow ends a new menu window.
func (gui *GUI) EndWindow() {
	nk.NkEnd(gui.ctx)
}

// BeginGroup starts a new group
func (gui *GUI) BeginGroup(name string, height float32) bool {
	nk.NkLayoutRowDynamic(gui.ctx, height, 1)
	update := nk.NkGroupBegin(gui.ctx, name, groupFlags)
	return update > 0
}

// EndGroup ends the current group
func (gui *GUI) EndGroup() {
	nk.NkGroupEnd(gui.ctx)
}

// Label draws a label with the specified text.
func (gui *GUI) Label(name string) {
	nk.NkLayoutRowDynamic(gui.ctx, 30, 1)
	{
		nk.NkLabel(gui.ctx, name, nk.TextCentered)
	}
}

// Button draws a button of the specified name and pressed state.
func (gui *GUI) Button(name string, isPressed *bool) {
	nk.NkLayoutRowDynamic(gui.ctx, 30, 1)
	{
		*isPressed = nk.NkButtonLabel(gui.ctx, name) > 0
	}
}

// Checkbox draw a checkbox with the specified label and checked state.
func (gui *GUI) Checkbox(name string, ischecked *bool) {
	nk.NkLayoutRowDynamic(gui.ctx, 30, 1)
	{
		var isactive int32 = 0
		if *ischecked {
			isactive = 1
		}
		nk.NkCheckboxLabel(gui.ctx, name, &isactive)
		*ischecked = isactive > 0
	}
}

// Selector draws a combo box with the specified items and returns the selected
// item index.
func (gui *GUI) Selector(name string, items []string, selectedIdx *int32) bool {
	width := gui.getWindowWidth()
	hasChanged := false
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 2)
	{
		nk.NkLayoutRowPush(gui.ctx, width*1/3)
		nk.NkLabel(gui.ctx, name, nk.TextLeft)
		nk.NkLayoutRowPush(gui.ctx, width*2/3)
		dim := nk.NkVec2(width*2/3, 200)
		oldIdx := *selectedIdx
		nk.NkCombobox(gui.ctx, items, int32(len(items)), selectedIdx, 20.0, dim)
		hasChanged = (oldIdx != *selectedIdx)
	}
	return hasChanged
}

// SliderInt32 draws a slider with the specified bounds, a value step size for
// the resulution of the slider values and returns the current value.
func (gui *GUI) SliderInt32(name string, value *int32, min, max, step int32) {
	width := gui.getWindowWidth()
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 2)
	{
		nk.NkLayoutRowPush(gui.ctx, width*1/3)
		nk.NkLabel(gui.ctx, name, nk.TextLeft)
		nk.NkLayoutRowPush(gui.ctx, width*2/3)
		nk.NkSliderInt(gui.ctx, min, value, max, step)
	}
}

// SliderFloat32 draws a slider with the specified bounds, a value step size for
// the resulution of the slider values and returns the current value.
func (gui *GUI) SliderFloat32(name string, value *float32, min, max, step float32) {
	width := gui.getWindowWidth()
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 2)
	{
		nk.NkLayoutRowPush(gui.ctx, width*1/3)
		nk.NkLabel(gui.ctx, name, nk.TextLeft)
		nk.NkLayoutRowPush(gui.ctx, width*2/3)
		nk.NkSliderFloat(gui.ctx, min, value, max, step)
	}
}

// Slider3 draws a slider with the specified bounds, a value step size for
// the resulution of the slider values and returns the current value.
func (gui *GUI) Slider3(name string, vec *mgl32.Vec3, min, max, step float32) {
	width := gui.getWindowWidth()
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 1)
	nk.NkLayoutRowPush(gui.ctx, width)
	nk.NkLabel(gui.ctx, name, nk.TextLeft)

	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 2)

	nk.NkLayoutRowPush(gui.ctx, width*1/3)
	nk.NkLabel(gui.ctx, fmt.Sprint(vec[0]), nk.TextRight)
	nk.NkLayoutRowPush(gui.ctx, width*2/3)
	nk.NkSliderFloat(gui.ctx, min, &vec[0], max, step)

	nk.NkLayoutRowPush(gui.ctx, width*1/3)
	nk.NkLabel(gui.ctx, fmt.Sprint(vec[1]), nk.TextRight)
	nk.NkLayoutRowPush(gui.ctx, width*2/3)
	nk.NkSliderFloat(gui.ctx, min, &vec[1], max, step)

	nk.NkLayoutRowPush(gui.ctx, width*1/3)
	nk.NkLabel(gui.ctx, fmt.Sprint(vec[2]), nk.TextRight)
	nk.NkLayoutRowPush(gui.ctx, width*2/3)
	nk.NkSliderFloat(gui.ctx, min, &vec[2], max, step)
}

// Input3 draws three value inputs for representing a vector 3D.
func (gui *GUI) Input3(name string, vec *mgl32.Vec3, min, max, step float32) {
	width := gui.getWindowWidth()
	widthl := width / 3
	widthr := (width-widthl)/3 - 2.5
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 30, 4)
	{
		nk.NkLayoutRowPush(gui.ctx, widthl)
		nk.NkLabel(gui.ctx, name, nk.TextLeft)
		nk.NkLayoutRowPush(gui.ctx, widthr)
		nk.NkPropertyFloat(gui.ctx, "X", min, &vec[0], max, step, 0.005)
		nk.NkLayoutRowPush(gui.ctx, widthr)
		nk.NkPropertyFloat(gui.ctx, "Y", min, &vec[1], max, step, 0.005)
		nk.NkLayoutRowPush(gui.ctx, widthr)
		nk.NkPropertyFloat(gui.ctx, "Z", min, &vec[2], max, step, 0.005)
	}
}

// ColorPicker provides a color picker widget that returns the components of the
// selected RGBa color.
func (gui *GUI) ColorPicker(name string, rgba *mgl32.Vec4) {
	color := nk.NkRgbaF(rgba.X(), rgba.Y(), rgba.Z(), rgba.W())
	colorf := nk.NkColorCf(color)
	width := gui.getWindowWidth()
	widthContent := width * 2 / 3
	nk.NkLayoutRowBegin(gui.ctx, nk.Static, 20, 2)
	nk.NkLayoutRowPush(gui.ctx, width*1/3)
	nk.NkLabel(gui.ctx, name, nk.TextLeft)
	nk.NkLayoutRowPush(gui.ctx, widthContent)
	dimensions := nk.NkVec2(widthContent, 200)
	if status := nk.NkComboBeginColor(gui.ctx, color, dimensions); status > 0 {
		nk.NkLayoutRowDynamic(gui.ctx, 190, 1)
		nk.NkColorPick(gui.ctx, &colorf, nk.ColorFormatRGBA)
		nk.NkComboEnd(gui.ctx)
	}
	rgba[0] = *colorf.GetR()
	rgba[1] = *colorf.GetG()
	rgba[2] = *colorf.GetB()
	rgba[3] = *colorf.GetA()
}

// OnCursorPosMove is a callback function for a cursor move event.
func (gui *GUI) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback function for a mouse button press event.
func (gui *GUI) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	status := nk.NkWindowIsAnyHovered(gui.ctx)
	return status == 1
}

// OnMouseScroll is a callback function for a mouse scroll event.
func (gui *GUI) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback function for a key press event.
func (gui *GUI) OnKeyPress(key, action, mods int) bool {
	return false
}

// OnResize is a callback function for a window resize event.
func (gui *GUI) OnResize(width, height int) bool {
	return false
}

func (gui *GUI) getWindowWidth() float32 {
	size := nk.NkWindowGetContentRegionSize(gui.ctx)
	return size.X() - 8
}

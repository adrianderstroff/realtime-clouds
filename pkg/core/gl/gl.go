// Package core provides an abstraction layer on top of OpenGL.
// It contains entities that provide utilities to simplify rendering.
package core

import (
	"fmt"

	ogl "github.com/go-gl/gl/v4.3-core/gl"
)

var state map[uint32]bool

func Init() error {
	// init opengl context
	if err := ogl.Init(); err != nil {
		return err
	}
	// print opengl version
	version := ogl.GoStr(ogl.GetString(ogl.VERSION))
	fmt.Println("OpenGL version", version)

	// set clear color
	Enable(DEPTH_TEST)
	FrontFace(CCW)
	CullFace(BACK)
	DepthFunc(LESS)
	ClearColor(0.0, 0.0, 0.0, 1.0)

	return nil
}

// Enable changes the state of the specified value.
// If the state of this value is already true then nothing happens.
func Enable(val uint32) {
	// check if the state is already true
	if state, ok := state[val]; ok && state == true {
		return
	}

	// case when the state of this value is false or unknown
	ogl.Enable(val)
	state[val] = true
}

// Disable changes the state of the specified value.
// If the state of this value is already false then nothing happens.
func Disable(val uint32) {
	// check if the state is already false
	if state, ok := state[val]; ok && state == false {
		return
	}

	// case when the state of this value is true or unknown
	ogl.Enable(val)
	state[val] = false
}

// GetError gets the last occured error and returns it or returns nil if no error occured.
// After receiving this error it is cleared from the perspective of OpenGL.
func GetError() error {
	err := ogl.GetError()
	if err != NO_ERROR {
		errorType := string(err)
		switch err {
		case INVALID_ENUM:
			errorType = "INVALID ENUM"
		case INVALID_VALUE:
			errorType = "INVALID VALUE"
		case INVALID_OPERATION:
			errorType = "INVALID OPERATION"
		case STACK_OVERFLOW:
			errorType = "STACK OVERFLOW"
		case STACK_UNDERFLOW:
			errorType = "STACK UNDERFLOW"
		case OUT_OF_MEMORY:
			errorType = "OUT OF MEMORY"
		case INVALID_FRAMEBUFFER_OPERATION:
			errorType = "INVALID FRAMEBUFFER OPERATION"
		case CONTEXT_LOST:
			errorType = "CONTEXT LOST"
		default:
			errorType = "UNKNOWN ERROR"
		}
		return fmt.Errorf("OpenGL error: %v", errorType)
	} else {
		return nil
	}
}

// Functions adapted from go-gl
var (
	Clear      = ogl.Clear
	ClearColor = ogl.ClearColor
	FrontFace  = ogl.FrontFace
	CullFace   = ogl.CullFace
	DepthFunc  = ogl.DepthFunc
)

// Errors
const (
	NO_ERROR                      = ogl.NO_ERROR
	INVALID_ENUM                  = ogl.INVALID_ENUM
	INVALID_VALUE                 = ogl.INVALID_VALUE
	INVALID_OPERATION             = ogl.INVALID_OPERATION
	STACK_OVERFLOW                = ogl.STACK_OVERFLOW
	STACK_UNDERFLOW               = ogl.STACK_UNDERFLOW
	OUT_OF_MEMORY                 = ogl.OUT_OF_MEMORY
	INVALID_FRAMEBUFFER_OPERATION = ogl.INVALID_FRAMEBUFFER_OPERATION
	CONTEXT_LOST                  = ogl.CONTEXT_LOST
)

// Capabilities that can be enabled or disabled
const (
	BLEND                         = ogl.BLEND
	CLIP_DISTANCE0                = ogl.CLIP_DISTANCE0
	CLIP_DISTANCE1                = ogl.CLIP_DISTANCE1
	CLIP_DISTANCE2                = ogl.CLIP_DISTANCE2
	CLIP_DISTANCE3                = ogl.CLIP_DISTANCE3
	CLIP_DISTANCE4                = ogl.CLIP_DISTANCE4
	CLIP_DISTANCE5                = ogl.CLIP_DISTANCE5
	CLIP_DISTANCE6                = ogl.CLIP_DISTANCE6
	CLIP_DISTANCE7                = ogl.CLIP_DISTANCE7
	COLOR_LOGIC_OP                = ogl.COLOR_LOGIC_OP
	CULL_FACE                     = ogl.CULL_FACE
	DEBUG_OUTPUT                  = ogl.DEBUG_OUTPUT
	DEBUG_OUTPUT_SYNCHRONOUS      = ogl.DEBUG_OUTPUT_SYNCHRONOUS
	DEPTH_CLAMP                   = ogl.DEPTH_CLAMP
	DEPTH_TEST                    = ogl.DEPTH_TEST
	DITHER                        = ogl.DITHER
	FRAMEBUFFER_SRGB              = ogl.FRAMEBUFFER_SRGB
	LINE_SMOOTH                   = ogl.LINE_SMOOTH
	MULTISAMPLE                   = ogl.MULTISAMPLE
	POLYGON_OFFSET_FILL           = ogl.POLYGON_OFFSET_FILL
	POLYGON_OFFSET_LINE           = ogl.POLYGON_OFFSET_LINE
	POLYGON_OFFSET_POINT          = ogl.POLYGON_OFFSET_POINT
	POLYGON_SMOOTH                = ogl.POLYGON_SMOOTH
	PRIMITIVE_RESTART             = ogl.PRIMITIVE_RESTART
	PRIMITIVE_RESTART_FIXED_INDEX = ogl.PRIMITIVE_RESTART_FIXED_INDEX
	RASTERIZER_DISCARD            = ogl.RASTERIZER_DISCARD
	SAMPLE_ALPHA_TO_COVERAGE      = ogl.SAMPLE_ALPHA_TO_COVERAGE
	SAMPLE_ALPHA_TO_ONE           = ogl.SAMPLE_ALPHA_TO_ONE
	SAMPLE_COVERAGE               = ogl.SAMPLE_COVERAGE
	SAMPLE_SHADING                = ogl.SAMPLE_SHADING
	SAMPLE_MASK                   = ogl.SAMPLE_MASK
	SCISSOR_TEST                  = ogl.SCISSOR_TEST
	STENCIL_TEST                  = ogl.STENCIL_TEST
	TEXTURE_CUBE_MAP_SEAMLESS     = ogl.TEXTURE_CUBE_MAP_SEAMLESS
	PROGRAM_POINT_SIZE            = ogl.PROGRAM_POINT_SIZE
)

// Values
const (
	CW                 = ogl.CW
	CCW                = ogl.CCW
	FRONT              = ogl.FRONT
	BACK               = ogl.BACK
	FRONT_AND_BACK     = ogl.FRONT_AND_BACK
	NEVER              = ogl.NEVER
	LESS               = ogl.LESS
	EQUAL              = ogl.EQUAL
	LEQUAL             = ogl.LEQUAL
	GREATER            = ogl.GREATER
	NOTEQUAL           = ogl.NOTEQUAL
	GEQUAL             = ogl.GEQUAL
	ALWAYS             = ogl.ALWAYS
	COLOR_BUFFER_BIT   = ogl.COLOR_BUFFER_BIT
	DEPTH_BUFFER_BIT   = ogl.DEPTH_BUFFER_BIT
	STENCIL_BUFFER_BIT = ogl.STENCIL_BUFFER_BIT
)

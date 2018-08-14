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
	Ptr                   = ogl.Ptr
	Clear                 = ogl.Clear
	ClearColor            = ogl.ClearColor
	FrontFace             = ogl.FrontFace
	CullFace              = ogl.CullFace
	DepthFunc             = ogl.DepthFunc
	GenTextures           = ogl.GenTextures
	DeleteTextures        = ogl.DeleteTextures
	BindTexture           = ogl.BindTexture
	BindImageTexture      = ogl.BindImageTexture
	ActiveTexture         = ogl.ActiveTexture
	TexParameteri         = ogl.TexParameteri
	TexParameteriv        = ogl.TexParameteriv
	TexParameterf         = ogl.TexParameterf
	TexParameterfv        = ogl.TexParameterfv
	TexParameterIiv       = ogl.TexParameterIiv
	TexParameterIuiv      = ogl.TexParameterIuiv
	TexImage1D            = ogl.TexImage1D
	TexImage2D            = ogl.TexImage2D
	TexImage3D            = ogl.TexImage3D
	TexImage2DMultisample = ogl.TexImage2DMultisample
	TexImage3DMultisample = ogl.TexImage3DMultisample
	GenerateMipmap        = ogl.GenerateMipmap
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
	CW                               = ogl.CW
	CCW                              = ogl.CCW
	FRONT                            = ogl.FRONT
	BACK                             = ogl.BACK
	FRONT_AND_BACK                   = ogl.FRONT_AND_BACK
	NEVER                            = ogl.NEVER
	LESS                             = ogl.LESS
	EQUAL                            = ogl.EQUAL
	LEQUAL                           = ogl.LEQUAL
	GREATER                          = ogl.GREATER
	NOTEQUAL                         = ogl.NOTEQUAL
	GEQUAL                           = ogl.GEQUAL
	ALWAYS                           = ogl.ALWAYS
	COLOR_BUFFER_BIT                 = ogl.COLOR_BUFFER_BIT
	DEPTH_BUFFER_BIT                 = ogl.DEPTH_BUFFER_BIT
	STENCIL_BUFFER_BIT               = ogl.STENCIL_BUFFER_BIT
	TEXTURE_1D                       = ogl.TEXTURE_1D
	PROXY_TEXTURE_1D                 = ogl.PROXY_TEXTURE_1D
	TEXTURE_2D                       = ogl.TEXTURE_2D
	PROXY_TEXTURE_2D                 = ogl.PROXY_TEXTURE_2D
	TEXTURE_3D                       = ogl.TEXTURE_3D
	TEXTURE_RECTANGLE                = ogl.TEXTURE_RECTANGLE
	PROXY_TEXTURE_RECTANGLE          = ogl.PROXY_TEXTURE_RECTANGLE
	TEXTURE_BUFFER                   = ogl.TEXTURE_BUFFER
	TEXTURE_CUBE_MAP                 = ogl.TEXTURE_CUBE_MAP
	PROXY_TEXTURE_CUBE_MAP           = ogl.PROXY_TEXTURE_CUBE_MAP
	TEXTURE_1D_ARRAY                 = ogl.TEXTURE_1D_ARRAY
	PROXY_TEXTURE_1D_ARRAY           = ogl.PROXY_TEXTURE_1D_ARRAY
	TEXTURE_2D_ARRAY                 = ogl.TEXTURE_2D_ARRAY
	TEXTURE_CUBE_MAP_ARRAY           = ogl.TEXTURE_CUBE_MAP_ARRAY
	TEXTURE_2D_MULTISAMPLE           = ogl.TEXTURE_2D_MULTISAMPLE
	TEXTURE_2D_MULTISAMPLE_ARRAY     = ogl.TEXTURE_2D_MULTISAMPLE_ARRAY
	TEXTURE_CUBE_MAP_POSITIVE_X      = ogl.TEXTURE_CUBE_MAP_POSITIVE_X
	TEXTURE_CUBE_MAP_POSITIVE_Y      = ogl.TEXTURE_CUBE_MAP_POSITIVE_Y
	TEXTURE_CUBE_MAP_POSITIVE_Z      = ogl.TEXTURE_CUBE_MAP_POSITIVE_Z
	TEXTURE_CUBE_MAP_NEGATIVE_X      = ogl.TEXTURE_CUBE_MAP_NEGATIVE_X
	TEXTURE_CUBE_MAP_NEGATIVE_Y      = ogl.TEXTURE_CUBE_MAP_NEGATIVE_Y
	TEXTURE_CUBE_MAP_NEGATIVE_Z      = ogl.TEXTURE_CUBE_MAP_NEGATIVE_Z
	MAX_TEXTURE_SIZE                 = ogl.MAX_TEXTURE_SIZE
	MAX_ARRAY_TEXTURE_LAYERS         = ogl.MAX_ARRAY_TEXTURE_LAYERS
	MAX_3D_TEXTURE_SIZE              = ogl.MAX_3D_TEXTURE_SIZE
	TEXTURE_BASE_LEVEL               = ogl.TEXTURE_BASE_LEVEL
	TEXTURE_MAX_LEVEL                = ogl.TEXTURE_MAX_LEVEL
	MAX_COMBINED_TEXTURE_IMAGE_UNITS = ogl.MAX_COMBINED_TEXTURE_IMAGE_UNITS
	REPEAT                           = ogl.REPEAT
	MIRRORED_REPEAT                  = ogl.MIRRORED_REPEAT
	CLAMP_TO_EDGE                    = ogl.CLAMP_TO_EDGE
	CLAMP_TO_BORDER                  = ogl.CLAMP_TO_BORDER
	LINEAR                           = ogl.LINEAR
	NEAREST                          = ogl.NEAREST
	TEXTURE_MIN_FILTER               = ogl.TEXTURE_MIN_FILTER
	TEXTURE_MAG_FILTER               = ogl.TEXTURE_MAG_FILTER
	TEXTURE_WRAP_R                   = ogl.TEXTURE_WRAP_R
	TEXTURE_WRAP_S                   = ogl.TEXTURE_WRAP_S
	TEXTURE_WRAP_T                   = ogl.TEXTURE_WRAP_T
	NEAREST_MIPMAP_NEAREST           = ogl.NEAREST_MIPMAP_NEAREST
	LINEAR_MIPMAP_NEAREST            = ogl.LINEAR_MIPMAP_NEAREST
	NEAREST_MIPMAP_LINEAR            = ogl.NEAREST_MIPMAP_LINEAR
	LINEAR_MIPMAP_LINEAR             = ogl.LINEAR_MIPMAP_LINEAR
	TEXTURE_SWIZZLE_R                = ogl.TEXTURE_SWIZZLE_R
	TEXTURE_SWIZZLE_G                = ogl.TEXTURE_SWIZZLE_G
	TEXTURE_SWIZZLE_B                = ogl.TEXTURE_SWIZZLE_B
	TEXTURE_SWIZZLE_A                = ogl.TEXTURE_SWIZZLE_A
	TEXTURE_SWIZZLE_RGBA             = ogl.TEXTURE_SWIZZLE_RGBA
	RED                              = ogl.RED
	RED_INTEGER                      = ogl.RED_INTEGER
	GREEN                            = ogl.GREEN
	BLUE                             = ogl.BLUE
	ALPHA                            = ogl.ALPHA
	RG                               = ogl.RG
	RG_INTEGER                       = ogl.RG_INTEGER
	ZERO                             = ogl.ZERO
	ONE                              = ogl.ONE
	RGB                              = ogl.RGB
	RGB_INTEGER                      = ogl.RGB_INTEGER
	BGR                              = ogl.BGR
	BGR_INTEGER                      = ogl.BGR_INTEGER
	RGBA                             = ogl.RGBA
	RGBA_INTEGER                     = ogl.RGBA_INTEGER
	BGRA                             = ogl.BGRA
	BGRA_INTEGER                     = ogl.BGRA_INTEGER
	DEPTH_STENCIL_TEXTURE_MODE       = ogl.DEPTH_STENCIL_TEXTURE_MODE
	DEPTH_COMPONENT                  = ogl.DEPTH_COMPONENT
	DEPTH_STENCIL                    = ogl.DEPTH_STENCIL
	STENCIL_INDEX                    = ogl.STENCIL_INDEX
	UNSIGNED_BYTE                    = ogl.UNSIGNED_BYTE
	BYTE                             = ogl.BYTE
	UNSIGNED_SHORT                   = ogl.UNSIGNED_SHORT
	SHORT                            = ogl.SHORT
	UNSIGNED_INT                     = ogl.UNSIGNED_INT
	INT                              = ogl.INT
	FLOAT                            = ogl.FLOAT
	UNSIGNED_BYTE_3_3_2              = ogl.UNSIGNED_BYTE_3_3_2
	UNSIGNED_BYTE_2_3_3_REV          = ogl.UNSIGNED_BYTE_2_3_3_REV
	UNSIGNED_SHORT_5_6_5             = ogl.UNSIGNED_SHORT_5_6_5
	UNSIGNED_SHORT_5_6_5_REV         = ogl.UNSIGNED_SHORT_5_6_5_REV
	UNSIGNED_SHORT_4_4_4_4           = ogl.UNSIGNED_SHORT_4_4_4_4
	UNSIGNED_SHORT_4_4_4_4_REV       = ogl.UNSIGNED_SHORT_4_4_4_4_REV
	UNSIGNED_SHORT_5_5_5_1           = ogl.UNSIGNED_SHORT_5_5_5_1
	UNSIGNED_SHORT_1_5_5_5_REV       = ogl.UNSIGNED_SHORT_1_5_5_5_REV
	UNSIGNED_INT_8_8_8_8             = ogl.UNSIGNED_INT_8_8_8_8
	UNSIGNED_INT_8_8_8_8_REV         = ogl.UNSIGNED_INT_8_8_8_8_REV
	UNSIGNED_INT_10_10_10_2          = ogl.UNSIGNED_INT_10_10_10_2
	UNSIGNED_INT_2_10_10_10_REV      = ogl.UNSIGNED_INT_2_10_10_10_REV
	TEXTURE0                         = ogl.TEXTURE0
	TEXTURE1                         = ogl.TEXTURE1
	TEXTURE2                         = ogl.TEXTURE2
	TEXTURE3                         = ogl.TEXTURE3
	TEXTURE4                         = ogl.TEXTURE4
	TEXTURE5                         = ogl.TEXTURE5
	TEXTURE6                         = ogl.TEXTURE6
	TEXTURE7                         = ogl.TEXTURE7
	TEXTURE8                         = ogl.TEXTURE8
	TEXTURE9                         = ogl.TEXTURE9
	TEXTURE10                        = ogl.TEXTURE10
	TEXTURE11                        = ogl.TEXTURE11
	TEXTURE12                        = ogl.TEXTURE12
	TEXTURE13                        = ogl.TEXTURE13
	TEXTURE14                        = ogl.TEXTURE14
	TEXTURE15                        = ogl.TEXTURE15
	TEXTURE16                        = ogl.TEXTURE16
	TEXTURE17                        = ogl.TEXTURE17
	TEXTURE18                        = ogl.TEXTURE18
	TEXTURE19                        = ogl.TEXTURE19
	TEXTURE20                        = ogl.TEXTURE20
	TEXTURE21                        = ogl.TEXTURE21
	TEXTURE22                        = ogl.TEXTURE22
	TEXTURE23                        = ogl.TEXTURE23
	TEXTURE24                        = ogl.TEXTURE24
	TEXTURE25                        = ogl.TEXTURE25
	TEXTURE26                        = ogl.TEXTURE26
	TEXTURE27                        = ogl.TEXTURE27
	TEXTURE28                        = ogl.TEXTURE28
	TEXTURE29                        = ogl.TEXTURE29
	TEXTURE30                        = ogl.TEXTURE30
	TEXTURE31                        = ogl.TEXTURE31
)

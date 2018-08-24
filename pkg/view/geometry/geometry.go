// Package geometry provides geometric primitives that can be used in meshes.
// It also offers a way to create custom geometric shapes.
package geometry

// Alignment options for different vertex attribute layouts
const (
	ALIGN_MULTI_BATCH  = 0 // having each attribute in a different slice (1111)(2222)(3333)
	ALIGN_SINGLE_BATCH = 1 // having all attributes in one slice but batched (111122223333)
	ALIGN_INTERLEAVED  = 2 // having all attributes in one slice but interleaved (123123123123)
)

// Geometry is a collection of vertex data and the way it's attributes are layed out.
// A vertex can have multiple attributes like position, normal, uv coordinate, etc.
// Those attributes are specified in the VertexAttribute.
// The data is a slice of slices to account for two data alignments.
// One alignment is having all attributes in different slices (1111)(2222)(3333).
// The other alignment is interleaving all attributes in one slice (123123123123).
// In the former case the length of each data slice must be a multiple of the
// correspondings dataType and count. In the latter case there is only one
// slice that has to be a length that is a multiple of all attributes dataTypes
// and counts combined.
type Geometry struct {
	Layout    []VertexAttribute
	Data      [][]float32
	Alignment int
}

// Make constructs a Geometry with it's layout and the data.
func Make(layout []VertexAttribute, data [][]float32) Geometry {
	// determine alignment
	alignment := ALIGN_MULTI_BATCH
	if len(data) == 1 {
		alignment = ALIGN_INTERLEAVED
	}

	return Geometry{
		Layout:    layout,
		Data:      data,
		Alignment: alignment,
	}
}

// New constructs a reference to Geometry with it's layout and the data.
func New(layout []VertexAttribute, data [][]float32) *Geometry {
	geometry := Make(layout, data)
	return &geometry
}

// VertexAttribute specifies the layout of one vertex attribute.
// The id has to match the name of the vertex attribute used in the shader.
// The glType is the type of one element of the vertex attribute to specify.
// The count specifies of how many elements this vertex attribute consists of.
// The usage is a hint to how the vertex data is going to be used. Allowed
// usage options are gl.STREAM_DRAW, gl.STREAM_READ, gl.STREAM_COPY,
// gl.STATIC_DRAW, gl.STATIC_READ, gl.STATIC_COPY, gl.DYNAMIC_DRAW,
// gl.DYNAMIC_READ, or gl.DYNAMIC_COPY.
type VertexAttribute struct {
	Id     string
	GlType uint32
	Count  int32
	Usage  int32
}

// VertexAttribute constructs a VertexAttribute with the given id, the type of one element, the number
// of elements and the OpenGL usage option.
func MakeVertexAttribute(id string, glType uint32, count int32, usage int32) VertexAttribute {
	return VertexAttribute{
		Id:     id,
		GlType: glType,
		Count:  count,
		Usage:  usage,
	}
}

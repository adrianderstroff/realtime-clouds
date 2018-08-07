package mesh

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/buffer/vao"
	geom "github.com/adrianderstroff/realtime-clouds/pkg/geometry/geometry"
	tex  "github.com/adrianderstroff/realtime-clouds/pkg/texture/texture"
)

// Mesh holds geometry data and textures that should be used to render this object.
// It uses the geometry to construct the vertex array object.
type struct Mesh{
	geometry geom.Geometry
	textures []tex.Texture
	vao VAO
}

// MakeMesh constructs a Mesh from it's geometry and a set of textures.
// By passing no textures only the geometry will be used to render this mesh.
func MakeMesh(geometry &geom.Geometry, textures []tex.Texture) Mesh {
	// make vao
	vao := buffer.MakeVAO(mode)

	// populate vao depending on the alignment of the geometry
	switch geometry.Alignment {
	case geom.ALIGN_MULTI_BATCH:
		// add multiple vbos specified by the geometries layout to the vao
		for i := 0; i < len(geometry.Layout); i++ {
			data := geometry.Data[i]
			attrib := geometry.Layout[i]
			vbo := MakeVBO(data, attrib.Count, attrib.Usage)
			vbo.AddVertexAttribute(attrib.Id, attrib.Count, attrib.GlType)
			vao.AddVertexBuffer(&vbo)
		}
	case geom.ALIGN_SINGLE_BATCH:
		// just for future compatibility
	case geom.ALIGN_INTERLEAVED:
		// count number of all elements of all vertex attributes
		count := 0
		for attrib := range geometry.Layout {
			count += attrib.Count
		}

		// add all vertex attributes to one vbo
		vbo := MakeVBO(geometry.Data[0], count, geometry.Layout[0].Usage)
		for attrib := range geometry.Layout {
			vbo.AddVertexAttribute(attrib.Id, attrib.Count, attrib.GlType)
		}
		vao.AddVertexBuffer(&vbo)
	}

	return Mesh{
		geometry: geometry,
		textures: textures,
		vao: vao
	}
}

// Delete destroy the Mesh and it's buffers.
func (mesh *Mesh) Delete() {
	mesh.vao.Delete()
}

// Build is called by the Shader.
// It sets up it's buffers.
func (mesh Mesh) Build(shaderProgramHandle uint32) {
	mesh.vao.BuildBuffers(shaderProgramHandle)
}

// Render draws the Mesh using the currently bound Shader.
func (mesh Mesh) Render() {
	mesh.vao.Render()
}

// RenderInstanced draws the Mesh multiple times specified by instancecount using the currently bound Shader.
func (mesh Mesh) RenderInstanced(instancecount int32) {
	mesh.vao.RenderInstanced(instancecount)
}

// GetVAO returns a pointer to the VAO.
func (mesh *Mesh) GetVAO() *VAO {
	return &mesh.vao
}

// SetVAO updates the VAO.
func (mesh *Mesh) SetVAO(vao VAO) {
	mesh.vao = vao
}
// Package mesh is used for creating meshes from geometry and textures.
// Meshes are entities that can be assigned to a ShaderProgram in order to render them.
package mesh

import (
	vao "github.com/adrianderstroff/realtime-clouds/pkg/buffer/vao"
	vbo "github.com/adrianderstroff/realtime-clouds/pkg/buffer/vbo"
	geom "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// Mesh holds geometry data and textures that should be used to render this object.
// It uses the geometry to construct the vertex array object.
type Mesh struct {
	geometry     geom.Geometry
	textures     []tex.Texture
	vao          vao.VAO
	onPreRender  func()
	onPostRender func()
}

// MakeMesh constructs a Mesh from it's geometry and a set of textures.
// By passing no textures only the geometry will be used to render this mesh.
func Make(geometry geom.Geometry, textures []tex.Texture, mode uint32) Mesh {
	// make vao
	vao := vao.Make(mode)

	// populate vao depending on the alignment of the geometry
	switch geometry.Alignment {
	case geom.ALIGN_MULTI_BATCH:
		// add multiple vbos specified by the geometries layout to the vao
		for i := 0; i < len(geometry.Layout); i++ {
			data := geometry.Data[i]
			attrib := geometry.Layout[i]
			vbo := vbo.Make(data, uint32(attrib.Count), uint32(attrib.Usage))
			vbo.AddVertexAttribute(attrib.Id, attrib.Count, attrib.GlType)
			vao.AddVertexBuffer(&vbo)
		}
	case geom.ALIGN_SINGLE_BATCH:
		// just for future compatibility
	case geom.ALIGN_INTERLEAVED:
		// count number of all elements of all vertex attributes
		var count int32 = 0
		for _, attrib := range geometry.Layout {
			count += attrib.Count
		}

		// add all vertex attributes to one vbo
		vbo := vbo.Make(geometry.Data[0], uint32(count), uint32(geometry.Layout[0].Usage))
		for _, attrib := range geometry.Layout {
			vbo.AddVertexAttribute(attrib.Id, attrib.Count, attrib.GlType)
		}
		vao.AddVertexBuffer(&vbo)
	}

	return Mesh{
		geometry: geometry,
		textures: textures,
		vao:      vao,
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
	// bind all textures in order
	for i, texture := range mesh.textures {
		texture.Bind(uint32(i))
	}
	// pre render event
	if mesh.onPreRender != nil {
		mesh.onPreRender()
	}
	// render geometry
	mesh.vao.Render()
	// post render event
	if mesh.onPostRender != nil {
		mesh.onPostRender()
	}
	// unbind all textures
	for _, texture := range mesh.textures {
		texture.Unbind()
	}
}

// RenderInstanced draws the Mesh multiple times specified by instancecount using the currently bound Shader.
func (mesh Mesh) RenderInstanced(instancecount int32) {
	// bind all textures in order
	for i, texture := range mesh.textures {
		texture.Bind(uint32(i))
	}
	// pre render event
	if mesh.onPreRender != nil {
		mesh.onPreRender()
	}
	// render geometry instanced
	mesh.vao.RenderInstanced(instancecount)
	// post render event
	if mesh.onPostRender != nil {
		mesh.onPostRender()
	}
	// unbind all textures
	for _, texture := range mesh.textures {
		texture.Unbind()
	}
}

// AddTexture adds a texture to the list of textures.
func (mesh *Mesh) AddTexture(texture tex.Texture) {
	mesh.textures = append(mesh.textures, texture)
}

// GetVAO returns a pointer to the VAO.
func (mesh *Mesh) GetVAO() *vao.VAO {
	return &mesh.vao
}

// SetVAO updates the VAO.
func (mesh *Mesh) SetVAO(vao vao.VAO) {
	mesh.vao = vao
}

// SetPreRenderAction sets an action that is executed each time
// before the mesh is being rendered.
func (mesh *Mesh) SetPreRenderAction(action func()) {
	mesh.onPreRender = action
}

// SetPostRenderAction sets an action that is executed each time
// after the mesh has been rendered.
func (mesh *Mesh) SetPostRenderAction(action func()) {
	mesh.onPostRender = action
}

package skybox

import (
	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	quad "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry/quad"
	mesh "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// Make constructs a skybox made from a quad with the cube map textures
// specified by the provided paths as well as the rendering mode.
func Make(sidelength float32, right, left, top, bottom, front, back string, mode uint32) (mesh.Mesh, error) {
	// make geometry
	geometry := quad.Make(sidelength, sidelength, sidelength, true)
	// make texture
	cubemap, err := tex.MakeCubeMap(right, left, top, bottom, front, back, true)
	if err != nil {
		return mesh.Mesh{}, err
	}
	textures := []tex.Texture{cubemap}
	// make mesh
	mesh := mesh.Make(geometry, textures, mode)
	// add actions
	prerender := func() {
		gl.DepthMask(false)
	}
	postrender := func() {
		gl.DepthMask(true)
	}
	mesh.SetPreRenderAction(prerender)
	mesh.SetPostRenderAction(postrender)
	return mesh, nil
}

// MakeFromDirectory constructs a skybox made from a quad with the specified side length
// in all 3 dimensions as well as the  the cube map textures specified by the provided
// directory and fileending as well as the rendering mode.
// The specified directory has to have all images in the same file format and the names
// of the files have to be named right, left, top, bottom, front and back respectively.
func MakeFromDirectory(sidelength float32, dir, fileending string, mode uint32) (mesh.Mesh, error) {
	right := dir + "right." + fileending
	left := dir + "left." + fileending
	top := dir + "top." + fileending
	bottom := dir + "bottom." + fileending
	front := dir + "front." + fileending
	back := dir + "back." + fileending
	return Make(sidelength, right, left, top, bottom, front, back, mode)
}

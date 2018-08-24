package skybox

import (
	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	quad "github.com/adrianderstroff/realtime-clouds/pkg/view/geometry/quad"
	mesh "github.com/adrianderstroff/realtime-clouds/pkg/view/mesh"
	tex "github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

// Make constructs a skybox made from a quad with the cube map textures
// specified by the provided paths as well as the rendering mode.
func Make(right, left, top, bottom, front, back string, mode uint32) (mesh.Mesh, error) {
	// make geometry
	geometry := quad.Make(50, 50, 50, false)
	// make texture
	cubemap, err := tex.MakeCubeMap(right, left, top, bottom, front, back)
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

// MakeFromDirectory constructs a skybox made from a quad with the cube map textures
// specified by the provided directory and fileending as well as the rendering mode.
// The specified directory has to have all images in the same file format and the names
// of the files have to be right, left, top, bottom, front, back respectively.
func MakeFromDirectory(dir, fileending string, mode uint32) (mesh.Mesh, error) {
	right := dir + "right." + fileending
	left := dir + "left." + fileending
	top := dir + "top." + fileending
	bottom := dir + "bottom." + fileending
	front := dir + "front." + fileending
	back := dir + "back." + fileending
	return Make(right, left, top, bottom, front, back, mode)
}

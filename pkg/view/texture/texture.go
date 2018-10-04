// Package texture provides classes for creating and storing images and textures.
package texture

import (
	"unsafe"

	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image2d"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image3d"
)

// Texture holds no to several images.
type Texture struct {
	handle uint32
	target uint32
	texPos uint32 // e.g. gl.TEXTURE0
}

// GetHandle returns the OpenGL of this texture.
func (tex *Texture) GetHandle() uint32 {
	return tex.handle
}

// MakeEmptyTexture creates a Texture with no image data.
func MakeEmpty() Texture {
	return Texture{0, gl.TEXTURE_2D, 0}
}

// Make creates a texture the given width and height.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func Make(width, height int, internalformat int32, format, pixelType uint32, data unsafe.Pointer, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t)

	// specify a texture image
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalformat, int32(width), int32(height), 0, format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColorTexture creates a color texture of the specified size.
func MakeColor(width, height int) Texture {
	return Make(width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepthTexture creates a depth texture of the specfied size.
func MakeDepth(width, height int) Texture {
	tex := Make(width, height, gl.DEPTH_COMPONENT, gl.DEPTH_COMPONENT, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
	return tex
}

// MakeCupeMapTexture creates a cube map with the images specfied from the path.
// For usage with skyboxes where textures are on the inside of the cube, set the
// inside parameter to true to flip all textures horizontally, otherwise set this
// parameter to false.
func MakeCubeMap(right, left, top, bottom, front, back string, inside bool) (Texture, error) {
	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// load images
	imagePaths := []string{right, left, top, bottom, front, back}
	for i, path := range imagePaths {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)
		image, err := image2d.MakeFromPath(path)
		if err != nil {
			return Texture{}, err
		}
		// if inside (e.g. for skyboxes) flip images horizontally
		if inside {
			image.FlipX()
		}
		gl.TexImage2D(target, 0, gl.RGBA, int32(image.GetWidth()), int32(image.GetHeight()),
			0, gl.RGBA, image.GetPixelType(), image.GetDataPointer())
	}

	// format texture
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// unset active texture
	tex.Unbind()

	return tex, nil
}

// MakeFromPath creates a texture with the image data specifed in path.
func MakeFromPath(path string, internalformat int32, format uint32) (Texture, error) {
	image, err := image2d.MakeFromPath(path)
	if err != nil {
		return Texture{}, err
	}

	image.FlipY()

	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// MakeFromImage grabs the dimensions and information from the image
func MakeFromImage(image *image2d.Image2D, internalformat int32, format uint32) Texture {
	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
}

// MakeFromData creates a texture
func MakeFromData(width, height int, internalformat int32, format uint32, data []uint8) (Texture, error) {
	image, err := image2d.MakeFromData(width, height, data)
	if err != nil {
		return Texture{}, err
	}

	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// MakeMultisampleTexture creates a multisample texture of the given width and height and the number of samples that should be used.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func MakeMultisample(width, height, samples int, format uint32, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D_MULTISAMPLE, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	/* gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_T, t) */

	// specify a texture image
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, int32(samples), format, int32(width), int32(height), false)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColorMultisampleTexture creates a multisample color texture of the given width and height and the number of samples that should be used.
func MakeColorMultisample(width, height, samples int) Texture {
	return MakeMultisample(width, height, samples, gl.RGBA,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepthMultisampleTexture creates a multisample depth texture of the given width and height and the number of samples that should be used.
func MakeDepthMultisample(width, height, samples int) Texture {
	return MakeMultisample(width, height, samples, gl.DEPTH_COMPONENT,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// Make3D constructs a 3D texture of the width and height of each image per slice and depth describing the number of slices.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded. The data layout is slices first then rows and lastly columns.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image. r specified the behaviour between the slices.
func Make3D(width, height, depth, internalformat int32, format, pixelType uint32, data unsafe.Pointer, min, mag, s, t, r int32) Texture {
	texture := Texture{0, gl.TEXTURE_3D, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_T, t)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_R, r)

	// specify a texture image
	gl.TexImage3D(gl.TEXTURE_3D, 0, internalformat, width, height, depth, 0, format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}

// Make3DFromPaths creates a 3D texture with the data of the images specifed by the provided paths.
func Make3DFromPath(paths []string, internalformat int32, format uint32) (Texture, error) {
	// load images from the specified paths and accumulate the loaded data
	images := []image2d.Image2D{}
	data := []uint8{}
	for _, path := range paths {
		image, err := image2d.MakeFromPath(path)
		if err != nil {
			return Texture{}, err
		}

		image.FlipY()

		data = append(data, image.GetData()...)
		images = append(images, image)
	}

	image := images[0]
	layers := int32(len(paths))
	return Make3D(int32(image.GetWidth()), int32(image.GetHeight()), layers, internalformat, format,
		image.GetPixelType(), gl.Ptr(data), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Make3DFromImage creates a 3D texture with the data of the 3D image.
func Make3DFromImage(image3d *image3d.Image3D, internalformat int32, format uint32) (Texture, error) {
	// load images from the specified paths and accumulate the loaded data
	data := image3d.GetData()

	return Make3D(int32(image3d.GetWidth()), int32(image3d.GetHeight()), int32(image3d.GetSlices()), internalformat, format,
		image3d.GetPixelType(), gl.Ptr(data), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Make3DFromImage creates a 3D texture with the data of the 3D image.
func Make3DFromData(data []uint8, width, height, slices int, internalformat int32, format uint32) (Texture, error) {
	return Make3D(int32(width), int32(height), int32(slices), internalformat, format,
		gl.UNSIGNED_BYTE, gl.Ptr(data), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Delete destroys the Texture.
func (tex *Texture) Delete() {
	gl.DeleteTextures(1, &tex.handle)
}

// GenMipmap generates mipmap levels.
// Chooses the two mipmaps that most closely match the size of the pixel being textured and uses the GL_LINEAR criterion to produce a texture value.
func (tex *Texture) GenMipmap() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

// GenMipmap generates mipmap levels.
// Chooses the mipmap that most closely matches the size of the pixel being textured and uses the GL_LINEAR criterion to produce a texture value.
func (tex *Texture) GenMipmapNearest() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

// Bind makes the texure available at the specified position.
func (tex *Texture) Bind(index uint32) {
	tex.texPos = gl.TEXTURE0 + index
	gl.ActiveTexture(tex.texPos)
	gl.BindTexture(tex.target, tex.handle)
}

// Unbind makes the texture unavailable for reading.
func (tex *Texture) Unbind() {
	tex.texPos = 0
	gl.BindTexture(tex.target, 0)
}

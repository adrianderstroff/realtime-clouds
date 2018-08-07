// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import (
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// Texture holds no to several images.
type Texture struct {
	handle uint32
	target uint32
	texPos uint32 // e.g. gl.TEXTURE0
}

// MakeEmptyTexture creates a Texture with no image data.
func MakeEmptyTexture() Texture {
	return Texture{0, gl.TEXTURE_2D, 0}
}

// MakeTexture creates a texture the given width and height.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func MakeTexture(width, height, internalformat int32, format, pixelType uint32, data unsafe.Pointer, min, mag, s, t int32) Texture {
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
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalformat, width, height, 0, format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColorTexture creates a color texture of the specified size.
func MakeColorTexture(width, height int32) Texture {
	return MakeTexture(width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepthTexture creates a depth texture of the specfied size.
func MakeDepthTexture(width, height int32) Texture {
	tex := MakeTexture(width, height, gl.DEPTH_COMPONENT, gl.DEPTH_COMPONENT, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
	return tex
}

// MakeCupeMapTexture creates a cube map with the images specfied from the path.
func MakeCubeMapTexture(right, left, top, bottom, front, back string) (Texture, error) {
	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// load images
	imagePaths := []string{right, left, top, bottom, front, back}
	for i, path := range imagePaths {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)
		image, err := MakeImage(path)
		if err != nil {
			return Texture{}, err
		}
		gl.TexImage2D(target, 0, image.internalFormat, image.width, image.height,
			0, image.format, image.pixelType, image.data)
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

// MakeTextureFromPath creates a texture with the image data specifed in path.
func MakeTextureFromPath(path string) (Texture, error) {
	image, err := MakeImage(path)
	if err != nil {
		return Texture{}, err
	}

	return MakeTexture(image.width, image.height, image.internalFormat, image.format,
		image.pixelType, image.data, gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// MakeMultisampleTexture creates a multisample texture of the given width and height and the number of samples that should be used.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func MakeMultisampleTexture(width, height, samples int32, format uint32, min, mag, s, t int32) Texture {
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
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, samples, format, width, height, false)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColorMultisampleTexture creates a multisample color texture of the given width and height and the number of samples that should be used.
func MakeColorMultisampleTexture(width, height, samples int32) Texture {
	return MakeMultisampleTexture(width, height, samples, gl.RGBA,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepthMultisampleTexture creates a multisample depth texture of the given width and height and the number of samples that should be used.
func MakeDepthMultisampleTexture(width, height, samples int32) Texture {
	return MakeMultisampleTexture(width, height, samples, gl.DEPTH_COMPONENT,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
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

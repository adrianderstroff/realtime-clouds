// Package texture provides classes for creating and storing images and textures.
package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// Image stores the dimensions, data format and it's pixel data.
// It can be used to manipulate single pixels and is used to
// upload it's data to a texture.
type Image struct {
	format         uint32
	internalFormat int32
	width          int32
	height         int32
	pixelType      uint32
	data           []uint8
}

// MakeImage constructs an image of the specified width and height and with all pixels set to the specified rgba value.
func MakeImage(width, height int32, r, g, b, a uint8) (Image, error) {
	// create image data
	var data []uint8
	length := width * height
	var i int32
	for i = 0; i < length; i++ {
		data = append(data, r)
		data = append(data, g)
		data = append(data, b)
		data = append(data, a)
	}

	return Image{
		format:         uint32(gl.RGBA),
		internalFormat: int32(gl.RGBA),
		width:          width,
		height:         height,
		pixelType:      uint32(gl.UNSIGNED_BYTE),
		data:           data,
	}, nil
}

// MakeImageFromPath constructs the image data from the specified path.
// If there is no image at the specified path an error is returned instead.
func MakeImageFromPath(path string) (Image, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return Image{}, err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return Image{}, err
	}

	// exctract rgba values
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Image{}, fmt.Errorf("Image not power of 2")
	}

	// get image dimensions
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)

	// rearrange data by flipping the rows of the image
	// this has to be done as the texture assumes the first row
	// to be the bottom row and the last row at the top.
	data := rgba.Pix
	var tempdata []uint8
	var row, col int32
	for row = height - 1; row >= 0; row-- {
		for col = 0; col < width; col++ {
			idx := (row*width + col) * 4
			tempdata = append(tempdata, data[idx])
			tempdata = append(tempdata, data[idx+1])
			tempdata = append(tempdata, data[idx+2])
			tempdata = append(tempdata, data[idx+3])
		}
	}
	data = tempdata

	return Image{
		format:         uint32(gl.RGBA),
		internalFormat: int32(gl.RGBA),
		width:          width,
		height:         height,
		pixelType:      uint32(gl.UNSIGNED_BYTE),
		data:           data,
	}, nil
}

// GetDataPointer returns an pointer to the beginning of the image data.
func (image *Image) GetDataPointer() unsafe.Pointer {
	return gl.Ptr(image.data)
}

// GetWidth returns the width of the image.
func (image *Image) GetWidth() int32 {
	return image.width
}

// GetHeight returns the height of the image.
func (image *Image) GetHeight() int32 {
	return image.height
}

// GetR returns the red value of the pixel at (x,y).
func (image *Image) GetR(x, y int32) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx]
}

// GetG returns the green value of the pixel at (x,y).
func (image *Image) GetG(x, y int32) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+1]
}

// GetB returns the blue value of the pixel at (x,y).
func (image *Image) GetB(x, y int32) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+2]
}

// GetA returns the alpha value of the pixel at (x,y).
func (image *Image) GetA(x, y int32) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+3]
}

// GetRGB returns the RGB values of the pixel at (x,y).
func (image *Image) GetRGB(x, y int32) (uint8, uint8, uint8) {
	idx := image.getIdx(x, y)
	return image.data[idx], image.data[idx+1], image.data[idx+2]
}

// GetRGBA returns the RGBA value of the pixel at (x,y).
func (image *Image) GetRGBA(x, y int32) (uint8, uint8, uint8, uint8) {
	idx := image.getIdx(x, y)
	return image.data[idx], image.data[idx+1], image.data[idx+2], image.data[idx+3]
}

// SetR sets the red value of the pixel at (x,y).
func (image *Image) SetR(x, y int32, r uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
}

// SetG sets the green value of the pixel at (x,y).
func (image *Image) SetG(x, y int32, g uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+1] = g
}

// SetB sets the blue value of the pixel at (x,y).
func (image *Image) SetB(x, y int32, b uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+2] = b
}

// SetA sets the alpha value of the pixel at (x,y).
func (image *Image) SetA(x, y int32, a uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+3] = a
}

// SetRGB sets the RGB values of the pixel at (x,y).
func (image *Image) SetRGB(x, y int32, r, g, b uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
	image.data[idx+1] = g
	image.data[idx+2] = b
}

// SetRGBA sets the RGBA values of the pixel at (x,y).
func (image *Image) SetRGBA(x, y int32, r, g, b, a uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
	image.data[idx+1] = g
	image.data[idx+2] = b
	image.data[idx+3] = a
}

// getIdx turns the x and y indices into a 1D index.
// The y coordinate is inverted as the first row of the
// image is at the bottom and the last one at the top.
// This reflects the way opengl interprets the pixel data
// when uploading it into a texture.
func (image *Image) getIdx(x, y int32) int32 {
	return ((image.height-1-y)*image.width + x) * 4
}

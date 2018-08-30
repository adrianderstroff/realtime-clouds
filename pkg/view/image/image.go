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

	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
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
func Make(width, height int32, r, g, b, a uint8) (Image, error) {
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

// MakeFromData constructs an image of the specified width and height and the specified data.
// The format specifies the format of the data e.g. RED, RG, RGB, RGBA, BGR, BGRA, RED_INTEGER,
// RG_INTEGER,  RGB_INTEGER, BGR_INTEGER, RGBA_INTEGER, BGRA_INTEGER, STENCIL_INDEX,
// DEPTH_COMPONENT or DEPTH_STENCIL.
func MakeFromData(width, height int32, format int, data []uint8) (Image, error) {
	return Image{
		format:         uint32(format),
		internalFormat: int32(gl.RGBA),
		width:          width,
		height:         height,
		pixelType:      uint32(gl.UNSIGNED_BYTE),
		data:           data,
	}, nil
}

// MakeImageFromPath constructs the image data from the specified path.
// If there is no image at the specified path an error is returned instead.
func MakeFromPath(path string) (Image, error) {
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

	return Image{
		format:         uint32(gl.RGBA),
		internalFormat: int32(gl.RGBA),
		width:          width,
		height:         height,
		pixelType:      uint32(gl.UNSIGNED_BYTE),
		data:           rgba.Pix,
	}, nil
}

// FlipX changes the order of the columns by swapping the first column of a row with the
// last column of the same row, the second column of this row with the second last column of this row etc.
func (image *Image) FlipX() {
	var tempdata []uint8
	var row, col int32
	for row = 0; row < image.height; row++ {
		for col = image.width - 1; col >= 0; col-- {
			idx := (row*image.width + col) * 4
			tempdata = append(tempdata, image.data[idx])
			tempdata = append(tempdata, image.data[idx+1])
			tempdata = append(tempdata, image.data[idx+2])
			tempdata = append(tempdata, image.data[idx+3])
		}
	}
	image.data = tempdata
}

// FlipY changes the order of the rows by swapping the first row with the
// last row, the second row with the second last row etc.
func (image *Image) FlipY() {
	var tempdata []uint8
	var row, col int32
	for row = image.height - 1; row >= 0; row-- {
		for col = 0; col < image.width; col++ {
			idx := (row*image.width + col) * 4
			tempdata = append(tempdata, image.data[idx])
			tempdata = append(tempdata, image.data[idx+1])
			tempdata = append(tempdata, image.data[idx+2])
			tempdata = append(tempdata, image.data[idx+3])
		}
	}
	image.data = tempdata
}

// GetFormat gets the format of the pixel data.
func (image *Image) GetFormat() uint32 {
	return image.format
}

// GetInternalFormat gets the number of color components in the texture
func (image *Image) GetInternalFormat() int32 {
	return image.internalFormat
}

// GetWidth returns the width of the image.
func (image *Image) GetWidth() int32 {
	return image.width
}

// GetHeight returns the height of the image.
func (image *Image) GetHeight() int32 {
	return image.height
}

// GetPixelType gets the data type of the pixel data.
func (image *Image) GetPixelType() uint32 {
	return image.pixelType
}

// GetDataPointer returns an pointer to the beginning of the image data.
func (image *Image) GetDataPointer() unsafe.Pointer {
	return gl.Ptr(image.data)
}

// GetData returns a copy of the images data
func (image *Image) GetData() []uint8 {
	return image.data
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
	return (y*image.width + x) * 4
}

// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

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

// Image holds image information of the data, dimensions and formats.
type Image struct {
	format         uint32
	internalFormat int32
	width          int32
	height         int32
	pixelType      uint32
	data           unsafe.Pointer
}

// MakeImage loads an Image from the specified path.
func MakeImage(path string) (Image, error) {
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

	return Image{
		uint32(gl.RGBA),
		int32(gl.RGBA),
		//int32(gl.SRGB_ALPHA),
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		uint32(gl.UNSIGNED_BYTE),
		gl.Ptr(rgba.Pix),
	}, nil
}

// RawImageData stores the image data on the CPU.
type RawImageData struct {
	data   []uint8
	width  int32
	height int32
}

// MakeRawImageData loads the raw image data from the specified path.
func MakeRawImageData(path string) (RawImageData, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return RawImageData{}, err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return RawImageData{}, err
	}

	// exctract rgba values
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return RawImageData{}, fmt.Errorf("Image not power of 2")
	}

	return RawImageData{
		rgba.Pix,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
	}, nil
}

// GetWidth returns the width of the image.
func (data *RawImageData) GetWidth() int32 {
	return data.width
}

// GetHeight returns the height of the image.
func (data *RawImageData) GetHeight() int32 {
	return data.height
}

// GetR returns the red value of the pixel at (x,y).
func (data *RawImageData) GetR(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx]
}

// GetG returns the green value of the pixel at (x,y).
func (data *RawImageData) GetG(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+1]
}

// GetB returns the blue value of the pixel at (x,y).
func (data *RawImageData) GetB(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+2]
}

// GetA returns the alpha value of the pixel at (x,y).
func (data *RawImageData) GetA(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+3]
}

// GetRGB returns the RGB values of the pixel at (x,y).
func (data *RawImageData) GetRGB(x, y int32) (uint8, uint8, uint8) {
	idx := data.getIdx(x, y)
	return data.data[idx], data.data[idx+1], data.data[idx+2]
}

// GetRGBA returns the RGBA value of the pixel at (x,y).
func (data *RawImageData) GetRGBA(x, y int32) (uint8, uint8, uint8, uint8) {
	idx := data.getIdx(x, y)
	return data.data[idx], data.data[idx+1], data.data[idx+2], data.data[idx+3]
}

// SetR sets the red value of the pixel at (x,y).
func (data *RawImageData) SetR(x, y int32, r uint8) {
	idx := data.getIdx(x, y)
	data.data[idx] = r
}

// SetG sets the green value of the pixel at (x,y).
func (data *RawImageData) SetG(x, y int32, g uint8) {
	idx := data.getIdx(x, y)
	data.data[idx+1] = g
}

// SetB sets the blue value of the pixel at (x,y).
func (data *RawImageData) SetB(x, y int32, b uint8) {
	idx := data.getIdx(x, y)
	data.data[idx+2] = b
}

// SetA sets the alpha value of the pixel at (x,y).
func (data *RawImageData) SetA(x, y int32, a uint8) {
	idx := data.getIdx(x, y)
	data.data[idx+3] = a
}

// SetRGB sets the RGB values of the pixel at (x,y).
func (data *RawImageData) SetRGB(x, y int32, r, g, b uint8) {
	idx := data.getIdx(x, y)
	data.data[idx] = r
	data.data[idx+1] = g
	data.data[idx+2] = b
}

// SetRGBA sets the RGBA values of the pixel at (x,y).
func (data *RawImageData) SetRGBA(x, y int32, r, g, b, a uint8) {
	idx := data.getIdx(x, y)
	data.data[idx] = r
	data.data[idx+1] = g
	data.data[idx+2] = b
	data.data[idx+3] = a
}

// getIdx turns the x and y indices into a 1D index.
func (data *RawImageData) getIdx(x, y int32) int32 {
	return (y*data.width + x) * 4
}

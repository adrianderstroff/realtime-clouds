// Package texture provides classes for creating and storing images.
package image2d

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"unsafe"

	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
)

// Image2D stores the dimensions, data format and it's pixel data.
// It can be used to manipulate single pixels and is used to
// upload it's data to a texture.
type Image2D struct {
	pixelType uint32
	width     int
	height    int
	channels  int
	data      []uint8
}

// Make constructs a white image of the specified width and height and number of channels.
func Make(width, height, channels int) (Image2D, error) {
	// early return if invalid dimensions had been specified
	err := checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// create image data
	var data []uint8
	length := width * height
	for i := 0; i < length; i++ {
		for c := 0; c < channels; c++ {
			data = append(data, 255)
		}
	}

	return Image2D{
		pixelType: uint32(gl.UNSIGNED_BYTE),
		width:     width,
		height:    height,
		channels:  channels,
		data:      data,
	}, nil
}

// MakeFromData constructs an image of the specified width and height and the specified data.
func MakeFromData(width, height int, data []uint8) (Image2D, error) {
	// data is stored as rgba value even if data is one channel only
	channels := len(data) / (width * height)

	// early return if invalid dimensions had been specified
	err := checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	return Image2D{
		pixelType: uint32(gl.UNSIGNED_BYTE),
		width:     width,
		height:    height,
		channels:  channels,
		data:      data,
	}, nil
}

// MakeFromPath constructs the image data from the specified path.
// If there is no image at the specified path an error is returned instead.
func MakeFromPath(path string) (Image2D, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return Image2D{}, err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return Image2D{}, err
	}

	// get image dimensions
	rect := img.Bounds()
	size := rect.Size()
	width := size.X
	height := size.Y

	// determine number of channels
	colormodel := img.ColorModel()
	channels := 4
	if colormodel == color.AlphaModel ||
		colormodel == color.Alpha16Model ||
		colormodel == color.GrayModel ||
		colormodel == color.Gray16Model {
		channels = 1
	}

	// early return if invalid dimensions had been specified
	err = checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// exctract data values
	var data []uint8
	switch channels {
	case 1:
		gray := image.NewGray(rect)
		draw.Draw(gray, rect, img, image.Pt(0, 0), draw.Src)
		data = gray.Pix
	case 4:
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rect, img, image.Pt(0, 0), draw.Src)
		data = rgba.Pix
	}

	return Image2D{
		pixelType: uint32(gl.UNSIGNED_BYTE),
		width:     width,
		height:    height,
		channels:  channels,
		data:      data,
	}, nil
}

// SaveToPath saves the image at the specified path in the png format.
// The specified image path has to have the fileextension .png.
// An error is thrown if the path is not valid or any of the specified
// directories don't exist.
func (img *Image2D) SaveToPath(path string) error {
	// create a file at the specified path
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// write data back into the golang image format
	rect := image.Rect(0, 0, img.width, img.height)
	switch img.channels {
	case 1:
		// fill image data
		out := image.NewGray(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 2:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = 0
				out.Pix[idxdst+3] = 255
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 3:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = img.data[idxsrc+2]
				out.Pix[idxdst+3] = 255
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 4:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
				out.Pix[idx+1] = img.data[idx+1]
				out.Pix[idx+2] = img.data[idx+2]
				out.Pix[idx+3] = img.data[idx+3]
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	}

	return nil
}

// FlipX changes the order of the columns by swapping the first column of a row with the
// last column of the same row, the second column of this row with the second last column of this row etc.
func (image *Image2D) FlipX() {
	var tempdata []uint8
	for row := 0; row < image.height; row++ {
		for col := image.width - 1; col >= 0; col-- {
			idx := image.getIdx(col, row)
			for c := 0; c < image.channels; c++ {
				tempdata = append(tempdata, image.data[idx+c])
			}
		}
	}
	image.data = tempdata
}

// FlipY changes the order of the rows by swapping the first row with the
// last row, the second row with the second last row etc.
func (image *Image2D) FlipY() {
	var tempdata []uint8
	for row := image.height - 1; row >= 0; row-- {
		for col := 0; col < image.width; col++ {
			idx := image.getIdx(col, row)
			for c := 0; c < image.channels; c++ {
				tempdata = append(tempdata, image.data[idx+c])
			}
		}
	}
	image.data = tempdata
}

// GetWidth returns the width of the image.
func (image *Image2D) GetWidth() int {
	return image.width
}

// GetHeight returns the height of the image.
func (image *Image2D) GetHeight() int {
	return image.height
}

// GetChannels return the number of the channels of the image.
func (image *Image2D) GetChannels() int {
	return image.channels
}

// GetPixelType gets the data type of the pixel data.
func (image *Image2D) GetPixelType() uint32 {
	return image.pixelType
}

// GetDataPointer returns an pointer to the beginning of the image data.
func (image *Image2D) GetDataPointer() unsafe.Pointer {
	return gl.Ptr(image.data)
}

// GetData returns a copy of the image's data
func (image *Image2D) GetData() []uint8 {
	cpy := make([]uint8, len(image.data))
	copy(cpy, image.data)
	return cpy
}

// GetR returns the red value of the pixel at (x,y).
func (image *Image2D) GetR(x, y int) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx]
}

// GetG returns the green value of the pixel at (x,y).
func (image *Image2D) GetG(x, y int) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+1]
}

// GetB returns the blue value of the pixel at (x,y).
func (image *Image2D) GetB(x, y int) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+2]
}

// GetA returns the alpha value of the pixel at (x,y).
func (image *Image2D) GetA(x, y int) uint8 {
	idx := image.getIdx(x, y)
	return image.data[idx+3]
}

// GetRGB returns the RGB values of the pixel at (x,y).
func (image *Image2D) GetRGB(x, y int) (uint8, uint8, uint8) {
	idx := image.getIdx(x, y)
	return image.data[idx],
		image.data[idx+1],
		image.data[idx+2]
}

// GetRGBA returns the RGBA value of the pixel at (x,y).
func (image *Image2D) GetRGBA(x, y int) (uint8, uint8, uint8, uint8) {
	idx := image.getIdx(x, y)
	return image.data[idx],
		image.data[idx+1],
		image.data[idx+2],
		image.data[idx+3]
}

// SetR sets the red value of the pixel at (x,y).
func (image *Image2D) SetR(x, y int, r uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
}

// SetG sets the green value of the pixel at (x,y).
func (image *Image2D) SetG(x, y int, g uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+1] = g
}

// SetB sets the blue value of the pixel at (x,y).
func (image *Image2D) SetB(x, y int, b uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+2] = b
}

// SetA sets the alpha value of the pixel at (x,y).
func (image *Image2D) SetA(x, y int, a uint8) {
	idx := image.getIdx(x, y)
	image.data[idx+3] = a
}

// SetRGB sets the RGB values of the pixel at (x,y).
func (image *Image2D) SetRGB(x, y int, r, g, b uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
	image.data[idx+1] = g
	image.data[idx+2] = b
}

// SetRGBA sets the RGBA values of the pixel at (x,y).
func (image *Image2D) SetRGBA(x, y int, r, g, b, a uint8) {
	idx := image.getIdx(x, y)
	image.data[idx] = r
	image.data[idx+1] = g
	image.data[idx+2] = b
	image.data[idx+3] = a
}

func (image Image2D) String() string {
	return fmt.Sprintf("Image2D (%v,%v) %v", image.width, image.height, image.channels)
}

// getIdx turns the x and y indices into a 1D index.
func (image *Image2D) getIdx(x, y int) int {
	return (x + y*image.width) * image.channels
}

func checkDimensions(width, height, channels int) error {
	if width < 1 || height < 1 {
		return errors.New("Width and height must be bigger than 0.")
	}

	if channels < 1 || channels > 4 {
		return errors.New("Number of channels must be between 1 and 4.")
	}

	return nil
}
